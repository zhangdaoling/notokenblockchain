package p2p

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/willf/bloom"
)

// errors
var (
	ErrMessageChannelFull = errors.New("message channel is full")
	ErrDuplicateMessage   = errors.New("reduplicate message")
)

const (
	bloomMaxItemCount = 100000
	bloomErrRate      = 0.001

	msgChanSize          = 1024
	maxDataLength        = 10000000 // 10MB
	routingQueryTimeout  = 10
	maxContinuousTimeout = 10
)

type Peer struct {
	id          peer.ID
	addr        multiaddr.Multiaddr
	direction   connDirection
	conn        network.Conn
	stream      network.Stream
	peerManager *PeerManager

	continuousTimeout    int
	lastRoutingQueryTime int64

	urgentMsgCh chan *p2pMessage
	normalMsgCh chan *p2pMessage

	quitWriteCh chan struct{}
	once        sync.Once

	bloomMutex     sync.Mutex
	bloomItemCount int
	recentMsg      *bloom.BloomFilter
}

// NewPeer returns a new instance of Peer struct.
func NewPeer(stream network.Stream, pm *PeerManager, direction connDirection) *Peer {
	peer := &Peer{
		id:          stream.Conn().RemotePeer(),
		addr:        stream.Conn().RemoteMultiaddr(),
		direction:   direction,
		conn:        stream.Conn(),
		stream:      stream,
		peerManager: pm,
		urgentMsgCh: make(chan *p2pMessage, msgChanSize),
		normalMsgCh: make(chan *p2pMessage, msgChanSize),
		quitWriteCh: make(chan struct{}),
		recentMsg:   bloom.NewWithEstimates(bloomMaxItemCount, bloomErrRate),
	}
	atomic.StoreInt64(&peer.lastRoutingQueryTime, (time.Now().Unix()))
	return peer
}

// ID return the net id.
func (p *Peer) ID() string {
	return p.id.Pretty()
}

// Addr return the address.
func (p *Peer) Addr() string {
	return p.addr.String()
}

// Start starts peer's loop.
func (p *Peer) Start() {
	go p.readLoop()
	go p.writeLoop()
}

// Stop stops peer's loop and cuts off the TCP connection.
func (p *Peer) Stop() {
	p.once.Do(func() {
		close(p.quitWriteCh)
	})
	p.conn.Close()
}

// SendMessage puts message into the corresponding channel.
func (p *Peer) SendMessage(msg *p2pMessage, mp MessagePriority, deduplicate bool) error {
	if deduplicate && msg.needDedup() {
		if p.hasMessage(msg) {
			// ilog.Debug("ignore reduplicate message")
			return ErrDuplicateMessage
		}
	}

	ch := p.urgentMsgCh
	if mp == NormalMessage {
		ch = p.normalMsgCh
	}
	select {
	case ch <- msg:
	default:
		fmt.Errorf("sending message failed. channel is full. messagePriority=%d", mp)
		return ErrMessageChannelFull
	}
	if msg.needDedup() {
		p.recordMessage(msg)
	}
	if msg.messageType() == RoutingTableQuery {
		p.routingQueryNow()
	}
	return nil
}

func (p *Peer) writeLoop() {
	for {
		select {
		case <-p.quitWriteCh:
			fmt.Printf("peer is stopped. pid=%v, addr=%v", p.ID(), p.addr)
			return
		case um := <-p.urgentMsgCh:
			p.write(um)
		case nm := <-p.normalMsgCh:
			for done := false; !done; {
				select {
				case <-p.quitWriteCh:
					fmt.Printf("peer is stopped. pid=%v, addr=%v", p.ID(), p.addr)
					return
				case um := <-p.urgentMsgCh:
					p.write(um)
				default:
					done = true
				}
			}
			p.write(nm)
		}
	}
}

func (p *Peer) readLoop() {
	header := make([]byte, dataBegin)
	for {
		_, err := io.ReadFull(p.stream, header)
		if err != nil {
			fmt.Printf("read header failed. err=%v", err)
			break
		}
		chainID := binary.BigEndian.Uint32(header[chainIDBegin:chainIDEnd])
		if chainID != p.peerManager.config.ChainID {
			fmt.Printf("Mismatched chainID, put peer to blacklist. remotePeer=%v, chainID=%d", p.ID(), chainID)
			break
		}
		length := binary.BigEndian.Uint32(header[dataLengthBegin:dataLengthEnd])
		if length > maxDataLength {
			fmt.Printf("data length too large: %d", length)
			break
		}
		data := make([]byte, dataBegin+length)
		_, err = io.ReadFull(p.stream, data[dataBegin:])
		if err != nil {
			fmt.Printf("read message failed. err=%v", err)
			break
		}
		copy(data[0:dataBegin], header)
		msg, err := parseP2PMessage(data)
		if err != nil {
			fmt.Printf("parse p2pmessage failed. err=%v", err)
			break
		}
		p.handleMessage(msg)
	}

	p.peerManager.RemoveNeighbor(p.id)
}

func (p *Peer) write(m *p2pMessage) error {
	// 5 kB/s
	deadline := time.Now().Add(time.Duration(len(m.content())/1024/5+3) * time.Second)
	if err := p.stream.SetWriteDeadline(deadline); err != nil {
		fmt.Printf("setting write deadline failed. err=%v, pid=%v", err, p.ID())
		p.peerManager.RemoveNeighbor(p.id)
		return err
	}
	_, err := p.stream.Write(m.content())
	if err != nil {
		fmt.Printf("writing message failed. err=%v, pid=%v", err, p.ID())
		if strings.Contains(err.Error(), "i/o timeout") {
			p.continuousTimeout++
			if p.continuousTimeout >= maxContinuousTimeout {
				fmt.Printf("max continuous timeout times, remove peer %v", p.ID())
				p.peerManager.RemoveNeighbor(p.id)
			}
		} else {
			p.peerManager.RemoveNeighbor(p.id)
		}
		return err
	}
	p.continuousTimeout = 0
	return nil
}

func (p *Peer) handleMessage(msg *p2pMessage) error {
	if msg.needDedup() {
		p.recordMessage(msg)
	}
	if msg.messageType() == RoutingTableResponse {
		if p.isRoutingQueryTimeout() {
			fmt.Printf("receive timeout routing response. pid=%v", p.ID())
			return nil
		}
		p.resetRoutingQueryTime()
	}
	p.peerManager.HandleMessage(msg, p.id)
	return nil
}

func (p *Peer) recordMessage(msg *p2pMessage) {
	p.bloomMutex.Lock()
	defer p.bloomMutex.Unlock()

	if p.bloomItemCount >= bloomMaxItemCount {
		p.recentMsg = bloom.NewWithEstimates(bloomMaxItemCount, bloomErrRate)
		p.bloomItemCount = 0
	}

	p.recentMsg.Add(msg.content())
	p.bloomItemCount++
}

func (p *Peer) hasMessage(msg *p2pMessage) bool {
	p.bloomMutex.Lock()
	defer p.bloomMutex.Unlock()

	return p.recentMsg.Test(msg.content())
}

// resetRoutingQueryTime resets last routing query time.
func (p *Peer) resetRoutingQueryTime() {
	atomic.StoreInt64(&p.lastRoutingQueryTime, -1)
}

// isRoutingQueryTimeout returns whether the last routing query time is too old.
func (p *Peer) isRoutingQueryTimeout() bool {
	return time.Now().Unix()-atomic.LoadInt64(&p.lastRoutingQueryTime) > routingQueryTimeout
}

// routingQueryNow sets the routing query time to the current timestamp.
func (p *Peer) routingQueryNow() {
	atomic.StoreInt64(&p.lastRoutingQueryTime, time.Now().Unix())
}

package p2p

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"github.com/zhangdaoling/simplechain/pb"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-dns"
)

var (
	dumpRoutingTableInterval = 5 * time.Minute
	syncRoutingTableInterval = 30 * time.Second
	metricsStatInterval      = 3 * time.Second
	findBPInterval           = 2 * time.Second

	dialTimeout        = 10 * time.Second
	deadPeerRetryTimes = 5
)

type connDirection int

const (
	inbound connDirection = iota
	outbound
)

const (
	defaultOutboundConn = 10
	defaultInboundConn  = 20

	bucketSize        = 1024
	peerResponseCount = 20
	maxPeerQuery      = 30
	maxAddrCount      = 10

	incomingMsgChanSize = 4096

	routingTableFile = "routing.table"
)

var (
	errInvalidMultiaddr = errors.New("invalid multiaddr string")
)

type PeerManager struct {
	quitCh chan struct{}
	config *P2PConfig

	host    host.Host
	started uint32
	wg      *sync.WaitGroup
	subs    *sync.Map //  map[MessageType]map[string]chan IncomingMessage

	neighbors     map[peer.ID]*Peer
	neighborCount map[connDirection]int
	neighborCap   map[connDirection]int
	neighborMutex sync.RWMutex

	routingTable   *kbucket.RoutingTable
	peerStore      peerstore.Peerstore
	retryTimes     map[string]int
	lastUpdateTime int64
	rtMutex        sync.RWMutex

	/*
		blackPIDs  map[string]bool
		blackIPs   map[string]bool
		blackMutex sync.RWMutex
	*/

	/*
		bpIDs   []peer.ID
		bpMutex sync.RWMutex
	*/
}

// NewPeerManager returns a new instance of PeerManager struct.
func NewPeerManager(host host.Host, config *P2PConfig) *PeerManager {
	routingTable := kbucket.NewRoutingTable(bucketSize, kbucket.ConvertPeerID(host.ID()), time.Second, host.Peerstore())
	pm := &PeerManager{
		quitCh: make(chan struct{}),
		config: config,

		host: host,
		wg:   new(sync.WaitGroup),
		subs: new(sync.Map),

		neighbors:     make(map[peer.ID]*Peer),
		neighborCount: make(map[connDirection]int),
		neighborCap:   make(map[connDirection]int),
		routingTable:  routingTable,
		peerStore:     host.Peerstore(),
		retryTimes:    make(map[string]int),
	}
	if config.InboundConn <= 0 {
		pm.neighborCap[inbound] = defaultOutboundConn
	} else {
		pm.neighborCap[inbound] = config.InboundConn
	}

	if config.OutboundConn <= 0 {
		pm.neighborCap[outbound] = defaultInboundConn
	} else {
		pm.neighborCap[outbound] = config.OutboundConn
	}

	return pm
}

func (pm *PeerManager) Start() error {
	//start once
	if !atomic.CompareAndSwapUint32(&pm.started, 0, 1) {
		return nil
	}

	pm.parseSeeds()

	pm.wg.Add(1)
	go pm.syncRoutingTableLoop()
	//pm.wg.Add(1)
	//go pm.metricsStatLoop()
	return nil
}

func (pm *PeerManager) Stop() {
	//stop once
	if !atomic.CompareAndSwapUint32(&pm.started, 1, 0) {
		return
	}

	close(pm.quitCh)
	pm.wg.Wait()
	pm.CloseAllNeighbors()
	return
}

func (pm *PeerManager) isStopped() bool {
	return atomic.LoadUint32(&pm.started) == 0
}

// Broadcast sends message to all the neighbors.
func (pm *PeerManager) Broadcast(data []byte, typ MessageType, mp MessagePriority) {
	msg := newP2PMessage(pm.config.ChainID, typ, pm.config.Version, defaultReservedFlag, data)

	wg := new(sync.WaitGroup)
	for _, p := range pm.GetAllNeighbors() {
		wg.Add(1)
		go func(p *Peer) {
			p.SendMessage(msg, mp, true)
			wg.Done()
		}(p)
	}
	wg.Wait()
}

// SendToPeer sends message to the specified peer.
func (pm *PeerManager) SendToPeer(peerID peer.ID, data []byte, typ MessageType, mp MessagePriority) {
	msg := newP2PMessage(pm.config.ChainID, typ, pm.config.Version, defaultReservedFlag, data)

	peer := pm.GetNeighbor(peerID)
	if peer != nil {
		peer.SendMessage(msg, mp, false)
	}
}

// Register registers a message channel of the given types.
func (pm *PeerManager) Register(id string, mTyps ...MessageType) chan IncomingMessage {
	if len(mTyps) == 0 {
		return nil
	}
	c := make(chan IncomingMessage, incomingMsgChanSize)
	for _, typ := range mTyps {
		m, _ := pm.subs.LoadOrStore(typ, new(sync.Map))
		m.(*sync.Map).Store(id, c)
	}
	return c
}

// Deregister deregisters a message channel of the given types.
func (pm *PeerManager) Deregister(id string, mTyps ...MessageType) {
	for _, typ := range mTyps {
		if m, exist := pm.subs.Load(typ); exist {
			m.(*sync.Map).Delete(id)
		}
	}
}

// HandleStream handles the incoming stream.
//
// It checks whether the remote peer already exists.
// If the peer is new and the neighbor count doesn't reach the threshold, it adds the peer into the neighbor list.
// If peer already exits, just add the stream to the peer.
// In other cases, reset the stream.
func (pm *PeerManager) HandleStream(s network.Stream, direction connDirection) {
	remotePID := s.Conn().RemotePeer()
	pm.freshPeer(remotePID)

	fmt.Printf("handle new stream. pid=%s, addr=%v, direction=%v", remotePID.Pretty(), s.Conn().RemoteMultiaddr(), direction)

	peer := pm.GetNeighbor(remotePID)
	if peer != nil {
		s.Reset()
		return
	}

	if pm.NeighborCount(direction) < pm.neighborCap[direction] {
		pm.AddNeighbor(NewPeer(s, pm, direction))
		return
	}
	pm.kickNormalNeighbors(direction)
	pm.AddNeighbor(NewPeer(s, pm, direction))
	return
}

// HandleMessage handles messages according to its type.
func (pm *PeerManager) HandleMessage(msg *p2pMessage, peerID peer.ID) {
	data, err := msg.data()
	if err != nil {
		fmt.Errorf("get message data failed. err=%v", err)
		return
	}
	switch msg.messageType() {
	case RoutingTableQuery:
		go pm.handleRoutingTableQuery(msg, peerID)
	case RoutingTableResponse:
		go pm.handleRoutingTableResponse(msg, peerID)
	default:
		inMsg := NewIncomingMessage(peerID, data, msg.messageType())
		if m, exist := pm.subs.Load(msg.messageType()); exist {
			m.(*sync.Map).Range(func(k, v interface{}) bool {
				select {
				case v.(chan IncomingMessage) <- *inMsg:
				default:
					fmt.Printf("sending incoming message failed. type=%s", msg.messageType())
				}
				return true
			})
		}
	}
}

// handleRoutingTableQuery picks the nearest peers of the given peerIDs and sends the result to inquirer.
func (pm *PeerManager) handleRoutingTableQuery(msg *p2pMessage, from peer.ID) {
	data, _ := msg.data()

	query := &pb.RoutingQuery{}
	err := proto.Unmarshal(data, query)
	if err != nil {
		fmt.Errorf("pb decode failed. err=%v, bytes=%v", err, data)
		return
	}

	queryIDs := query.GetIds()
	fmt.Printf("handling routing table query. %v", queryIDs)
	bytes, _ := pm.getRoutingResponse(queryIDs)
	if len(bytes) > 0 {
		pm.SendToPeer(from, bytes, RoutingTableResponse, UrgentMessage)
	}
}

// handleRoutingTableResponse stores the peer information received.
func (pm *PeerManager) handleRoutingTableResponse(msg *p2pMessage, from peer.ID) { // nolint

	data, _ := msg.data()
	resp := &pb.RoutingResponse{}
	err := proto.Unmarshal(data, resp)
	if err != nil {
		fmt.Errorf("Decoding pb failed. err=%v, bytes=%v", err, data)
		return
	}
	fmt.Printf("Receiving peer infos: %v, from=%v", resp, from.Pretty())
	for _, peerInfo := range resp.Peers {
		if len(peerInfo.Addrs) > 0 {
			pid, err := peer.IDB58Decode(peerInfo.Id)
			if err != nil {
				fmt.Printf("Decoding peerID failed. err=%v, id=%v", err, peerInfo.Id)
				continue
			}

			if pm.isDead(pid) { // ignore bad node's addr
				continue
			}

			if pid == pm.host.ID() { // ignore self's addr
				continue
			}

			// TODO: Take more reasonable measures to avoid routing pollution
			if pm.GetNeighbor(pid) != nil && from != pid { // ignore neighbor's addr
				continue
			}

			addrs := make([]string, 0, len(peerInfo.Addrs))
			if from != pid {
				for _, addr := range peerInfo.Addrs {
					if isPublicMaddr(addr) {
						addrs = append(addrs, addr)
					}
				}
			} else {
				// choose public multiaddr if exists, else take out port and combine it with remoteAddr
				var port string
				var hasPublicMaddr bool
				for _, addr := range peerInfo.Addrs {
					if isPublicMaddr(addr) {
						addrs = append(addrs, addr)
						hasPublicMaddr = true
					} else {
						port = addr[strings.LastIndex(addr, "/")+1:]
					}
				}
				if !hasPublicMaddr {
					neighbor := pm.GetNeighbor(pid)
					if neighbor != nil {
						remoteAddr := neighbor.addr.String()
						remoteListenAddr := remoteAddr[:strings.LastIndex(remoteAddr, "/")+1] + port
						addrs = append(addrs, remoteListenAddr)
					}
				}
			}
			maddrs := make([]multiaddr.Multiaddr, 0, len(addrs))
			for _, addr := range addrs {
				ma, err := multiaddr.NewMultiaddr(addr)
				if err != nil {
					fmt.Printf("Parsing multiaddr failed. err=%v, addr=%v", err, addr)
					continue
				}
				maddrs = append(maddrs, ma)
			}
			if len(maddrs) > maxAddrCount {
				maddrs = maddrs[:maxAddrCount]
			}
			if len(maddrs) > 0 {
				pm.storePeerInfo(pid, maddrs)
			}
		}
	}
}

func (pm *PeerManager) parseSeeds() {
	for _, seed := range pm.config.SeedNodes {
		peerID, addr, err := parseMultiaddr(seed)
		if err != nil {
			fmt.Errorf("parse seed nodes error. seed=%s, err=%v", seed, err)
			continue
		}

		if madns.Matches(addr) {
			resAddrs, err := madns.Resolve(context.Background(), addr)
			if err != nil {
				fmt.Errorf("resolve multiaddr failed. err=%v, addr=%v", err, addr)
				return
			}
			pm.storePeerInfo(peerID, resAddrs)
		} else {
			pm.storePeerInfo(peerID, []multiaddr.Multiaddr{addr})
		}
	}
}

func (pm *PeerManager) syncRoutingTableLoop() {
	pm.routingQuery([]string{pm.host.ID().Pretty()})

	defer pm.wg.Done()
	for {
		select {
		case <-pm.quitCh:
			return
		case <-time.After(syncRoutingTableInterval):
			pid, _ := randomPID()
			pm.routingQuery([]string{pid.Pretty()})
		}
	}
}

func (pm *PeerManager) routingQuery(ids []string) {
	if len(ids) == 0 {
		return
	}
	query := &pb.RoutingQuery{
		Ids: ids,
	}
	bytes, err := proto.Marshal(query)
	if err != nil {
		fmt.Errorf("pb encode failed. err=%v, obj=%+v", err, query)
		return
	}

	pm.Broadcast(bytes, RoutingTableQuery, UrgentMessage)
	outboundNeighborCount := pm.NeighborCount(outbound)
	if outboundNeighborCount >= pm.neighborCap[outbound] {
		return
	}
	allPeerIDs := pm.routingTable.ListPeers()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(allPeerIDs))

	for i, t := 0, 0; i < len(perm) && t < pm.neighborCap[outbound]-outboundNeighborCount; i++ {
		if pm.isStopped() {
			return
		}

		peerID := allPeerIDs[perm[i]]
		if peerID == pm.host.ID() {
			continue
		}
		if pm.GetNeighbor(peerID) != nil {
			continue
		}
		fmt.Printf("dial peer: pid=%v", peerID.Pretty())
		stream, err := pm.newStream(peerID)
		if err != nil {
			fmt.Printf("create stream failed. pid=%s, err=%v", peerID.Pretty(), err)

			if strings.Contains(err.Error(), "connected to wrong peer") {
				pm.deletePeerInfo(peerID)
				continue
			}

			pm.recordDialFail(peerID)
			if pm.isDead(peerID) {
				pm.deletePeerInfo(peerID)
			}
			continue
		}
		pm.HandleStream(stream, outbound)
		pm.SendToPeer(peerID, bytes, RoutingTableQuery, UrgentMessage)
		t++
	}
}

func (pm *PeerManager) getRoutingResponse(peerIDs []string) ([]byte, error) {
	queryIDs := peerIDs
	if len(queryIDs) > maxPeerQuery {
		queryIDs = queryIDs[:maxPeerQuery]
	}

	pidSet := make(map[peer.ID]struct{})
	for _, queryID := range queryIDs {
		pid, err := peer.IDB58Decode(queryID)
		if err != nil {
			fmt.Printf("decode peerID failed. err=%v, id=%v", err, queryID)
			continue
		}
		peerIDs := pm.routingTable.NearestPeers(kbucket.ConvertPeerID(pid), peerResponseCount)
		for _, id := range peerIDs {
			if !pm.isDead(id) {
				pidSet[id] = struct{}{}
			}
		}
	}

	resp := &pb.RoutingResponse{}
	for pid := range pidSet {
		info := pm.peerStore.PeerInfo(pid)
		if len(info.Addrs) > 0 {
			peerInfo := &pb.PeerInfo{
				Id: info.ID.Pretty(),
			}
			for _, addr := range info.Addrs {
				if isPublicMaddr(addr.String()) {
					peerInfo.Addrs = append(peerInfo.Addrs, addr.String())
				}
			}
			if len(peerInfo.Addrs) > maxAddrCount {
				peerInfo.Addrs = peerInfo.Addrs[:maxAddrCount]
			}
			if len(peerInfo.Addrs) > 0 {
				resp.Peers = append(resp.Peers, peerInfo)
			}
		}
	}
	selfInfo := &pb.PeerInfo{Id: pm.host.ID().Pretty()}
	for _, addr := range pm.host.Addrs() {
		selfInfo.Addrs = append(selfInfo.Addrs, addr.String())
	}
	resp.Peers = append(resp.Peers, selfInfo)

	bytes, err := proto.Marshal(resp)
	if err != nil {
		fmt.Errorf("pb encode failed. err=%v, obj=%+v", err, resp)
		return nil, err
	}
	return bytes, nil
}

func (pm *PeerManager) storePeerInfo(peerID peer.ID, addrs []multiaddr.Multiaddr) {
	// SetAddrs won't overwrite old value. clear and add.
	pm.peerStore.ClearAddrs(peerID)
	pm.peerStore.AddAddrs(peerID, addrs, peerstore.PermanentAddrTTL)
	pm.routingTable.Update(peerID)
	atomic.StoreInt64(&pm.lastUpdateTime, (time.Now().Unix()))
}

// deletePeerInfo deletes peer information in peerStore and routingTable. It doesn't need lock since the
// peerStore.ClearAddrs and routingTable.Update function are thread safe.
func (pm *PeerManager) deletePeerInfo(peerID peer.ID) {
	pm.peerStore.ClearAddrs(peerID)
	pm.routingTable.Remove(peerID)
	atomic.StoreInt64(&pm.lastUpdateTime, (time.Now().Unix()))
}

// AddNeighbor starts a peer and adds it to the neighbor list.
func (pm *PeerManager) AddNeighbor(p *Peer) {

	pm.neighborMutex.Lock()
	defer pm.neighborMutex.Unlock()

	if pm.neighbors[p.id] == nil {
		p.Start()
		// pm.storePeerInfo(p.id, []multiaddr.Multiaddr{p.addr})
		pm.neighbors[p.id] = p
		pm.neighborCount[p.direction]++
	}
}

// RemoveNeighbor stops a peer and removes it from the neighbor list.
func (pm *PeerManager) RemoveNeighbor(peerID peer.ID) {

	pm.neighborMutex.Lock()
	defer pm.neighborMutex.Unlock()

	p := pm.neighbors[peerID]
	if p != nil {
		p.Stop()
		delete(pm.neighbors, peerID)
		pm.neighborCount[p.direction]--
	}
}

// GetNeighbor returns the peer of the given peerID from the neighbor list.
func (pm *PeerManager) GetNeighbor(peerID peer.ID) *Peer {
	pm.neighborMutex.RLock()
	defer pm.neighborMutex.RUnlock()

	return pm.neighbors[peerID]
}

// GetAllNeighbors returns the peer of the given peerID from the neighbor list.
func (pm *PeerManager) GetAllNeighbors() []*Peer {
	pm.neighborMutex.RLock()
	defer pm.neighborMutex.RUnlock()

	peers := make([]*Peer, 0, len(pm.neighbors))
	for _, p := range pm.neighbors {
		peers = append(peers, p)
	}
	return peers
}

// CloseAllNeighbors close all connections.
func (pm *PeerManager) CloseAllNeighbors() {
	for _, p := range pm.GetAllNeighbors() {
		p.Stop()
	}
}

// AllNeighborCount returns the total neighbor amount.
func (pm *PeerManager) AllNeighborCount() int {
	pm.neighborMutex.RLock()
	defer pm.neighborMutex.RUnlock()

	return len(pm.neighbors)
}

// NeighborCount returns the neighbor amount of the given direction.
func (pm *PeerManager) NeighborCount(direction connDirection) int {
	pm.neighborMutex.RLock()
	defer pm.neighborMutex.RUnlock()

	return pm.neighborCount[direction]
}

// kickNormalNeighbors removes neighbors that are not block producers.
func (pm *PeerManager) kickNormalNeighbors(direction connDirection) {
	pm.neighborMutex.Lock()
	defer pm.neighborMutex.Unlock()

	for _, p := range pm.neighbors {
		if pm.neighborCount[direction] < pm.neighborCap[direction] {
			return
		}
		if direction == p.direction {
			p.Stop()
			delete(pm.neighbors, p.id)
			pm.neighborCount[direction]--
		}
	}
}

func (pm *PeerManager) newStream(pid peer.ID) (network.Stream, error) {
	ctx, _ := context.WithTimeout(context.Background(), dialTimeout) // nolint
	return pm.host.NewStream(ctx, pid, protocolID)
}

func (pm *PeerManager) recordDialFail(pid peer.ID) {
	pm.rtMutex.Lock()
	defer pm.rtMutex.Unlock()

	pm.retryTimes[pid.Pretty()]++
}

func (pm *PeerManager) freshPeer(pid peer.ID) {
	pm.rtMutex.Lock()
	defer pm.rtMutex.Unlock()

	delete(pm.retryTimes, pid.Pretty())
}

func (pm *PeerManager) isDead(pid peer.ID) bool {
	pm.rtMutex.RLock()
	defer pm.rtMutex.RUnlock()

	return pm.retryTimes[pid.Pretty()] > deadPeerRetryTimes
}

func parseMultiaddr(s string) (peer.ID, multiaddr.Multiaddr, error) {
	strs := strings.Split(s, "/ipfs/")
	if len(strs) != 2 {
		return "", nil, errInvalidMultiaddr
	}
	addr, err := multiaddr.NewMultiaddr(strs[0])
	if err != nil {
		return "", nil, err
	}
	peerID, err := peer.IDB58Decode(strs[1])
	if err != nil {
		return "", nil, err
	}
	return peerID, addr, nil
}

func randomPID() (peer.ID, error) {
	_, pubkey, err := crypto.GenerateEd25519Key(cryptorand.Reader)
	if err != nil {
		return "", err
	}
	return peer.IDFromPublicKey(pubkey)
}

// private IP:
// 10.0.0.0    - 10.255.255.255
// 192.168.0.0 - 192.168.255.255
// 172.16.0.0  - 172.31.255.255
func privateIP(ip string) (bool, error) {
	IP := net.ParseIP(ip)
	if IP == nil {
		return false, errors.New("invalid IP")
	}
	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	return private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP), nil
}

func isPublicMaddr(s string) bool {
	ip := getIPFromMaddr(s)
	if ip == "127.0.0.1" {
		return false
	}
	private, err := privateIP(ip)
	if err != nil {
		return false
	}
	return !private
}

var ipReg = regexp.MustCompile(`/((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))/`)

func getIPFromMaddr(s string) string {
	str := ipReg.FindString(s)
	if len(str) > 2 {
		return str[1 : len(str)-1]
	}
	return ""
}

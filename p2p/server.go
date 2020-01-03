package p2p

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

const (
	protocolID  = "notokenblockchain/0.0.1"
	privKeyFile = "priv.key"
)

// p2p service interface
type Server interface {
	//server interface
	Start() error
	Stop()

	//func interface
	Broadcast([]byte, MessageType, MessagePriority)
	SendToPeer(peer.ID, []byte, MessageType, MessagePriority)
	Register(string, ...MessageType) chan IncomingMessage
	Deregister(string, ...MessageType)

	//admin interface
	//GetAllNeighbors() []*Peer
	//ID() string
	//ConnectBPs([]string)
	//PutPeerToBlack(string)
}

var _ Server = &P2PServer{}

type P2PServer struct {
	host host.Host
	*PeerManager
	config *P2PConfig
}

func NewP2PServer(c *P2PConfig) (s *P2PServer, err error) {
	s = &P2PServer{
		config: c,
	}
	if err := os.MkdirAll(c.DataPath, 0755); err != nil {
		fmt.Errorf("failed to create p2p datapath, err=%v, path=%v", err, c.DataPath)
		return nil, err
	}

	privKey, err := getOrCreateKey(filepath.Join(c.DataPath, privKeyFile))
	if err != nil {
		fmt.Errorf("failed to get private key. err=%v, path=%v", err, c.DataPath)
		return nil, err
	}

	host, err := s.startHost(privKey, c.ListenAddr)
	if err != nil {
		fmt.Errorf("failed to start a host. err=%v, listenAddr=%v", err, c.ListenAddr)
		return nil, err
	}
	s.host = host

	s.PeerManager = NewPeerManager(host, c)
	return s, nil
}

func (s *P2PServer) startHost(pk crypto.PrivKey, listenAddr string) (host.Host, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(context.Background(),
		libp2p.Identity(pk),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", tcpAddr.IP, tcpAddr.Port)))
	if err != nil {
		return nil, err
	}
	h.SetStreamHandler(protocolID, s.streamHandler)
	return h, nil
}

func (s *P2PServer) streamHandler(stream network.Stream) {
	s.PeerManager.HandleStream(stream, inbound)
}

func (s *P2PServer) Start() error {
	go s.PeerManager.Start()
	return nil
}

func (s *P2PServer) Stop() {
	s.PeerManager.Stop()
	return
}

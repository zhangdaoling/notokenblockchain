package p2p

// P2PConfig is the config for p2p network.
type P2PConfig struct {
	ChainID  uint32
	Version  uint16
	DataPath string

	ListenAddr   string
	SeedNodes    []string
	InboundConn  int
	OutboundConn int
	BlackPID     []string
	BlackIP      []string
}

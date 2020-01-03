package consensus

var Consensus = mainChainConsensus

var mainChainConsensus = ConsensusInfo{
	ChainID:                   1,
	MessageSizeLimit:          65536,
	MessageExpirationTime:     1000,
	MessageExpirationHeight:   10,
	MessageDifficultyInterval: 1000,
	BlockDifficultyInterval:   100,
	MaxMessageInBlock:         100,
}

type ConsensusInfo struct {
	ChainID                   int64
	MessageSizeLimit          int64
	MessageExpirationTime     int64
	MessageExpirationHeight   int64
	MessageDifficultyInterval int64
	BlockDifficultyInterval   int64
	MaxMessageInBlock         int64
}

package geneis

var Geneis = mainChainGeneis

var mainChainGeneis = GeneisInfo{
	BlockDifficuty:   100,
	MessageDifficuty: 10,
}

type GeneisInfo struct {
	BlockDifficuty   int64
	MessageDifficuty int64
}

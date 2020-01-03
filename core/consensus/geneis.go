package consensus

var Geneis = mainChainGeneis

var mainChainGeneis = GeneisInfo{
	BlockDifficulty:   100,
	MessageDifficuty: 10,
}

type GeneisInfo struct {
	BlockDifficulty   int64
	MessageDifficuty int64
}

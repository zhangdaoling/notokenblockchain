package consensus

import (
	"github.com/zhangdaoling/notokenblockchain/core/types"
	"github.com/zhangdaoling/notokenblockchain/crypto"
	"github.com/zhangdaoling/notokenblockchain/pb"
)

var Geneis = mainChainGeneis

var mainChainGeneis = GeneisInfo{
	BlockDifficulty:  100,
	MessageDifficuty: 10,
}

type GeneisInfo struct {
	BlockDifficulty  int64
	MessageDifficuty int64
}

var GeneisChainState = types.ChainState{
	Difficulty:   0,
	Length:       1,
	TotalWeight:  0,
	TotalMessage: 0,
}

var GeneisBlock = types.Block{
	PBBlock: &pb.Block{
		Head: &pb.BlockHead{
			Version:           Consensus.Version,
			ChainId:           Consensus.ChainID,
			Height:            0,
			Time:              0,
			ParentHash:        nil,
			MessageMerkleHash: nil,
			PublicKey:         nil,
		},
		Signature: &pb.Signature{
			Algorithm: int32(crypto.Ed25519),
			Sig:       nil,
		},
		MessageHashes: nil,
		Messages:      nil,
	},
}

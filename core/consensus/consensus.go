package consensus

import (
	"errors"
)

var (
	ErrDifficultHeightErr = errors.New("new difficulty in wrong interval ")
)

var Consensus = mainChainConsensus

var mainChainConsensus = ConsensusInfo{
	Version:                 1,
	ChainID:                 1,
	MessageSizeLimit:        65536,
	MessageExpirationTime:   1000,
	MessageExpirationHeight: 10,
	DifficultyInterval:      100,
	MaxMessageInBlock:       100,
}

type ConsensusInfo struct {
	Version                 int64
	ChainID                 int64
	MessageSizeLimit        int64
	MessageExpirationTime   int64
	MessageExpirationHeight int64
	DifficultyInterval      int64
	MaxMessageInBlock       int64
	ExpectBlockInterval     int64
}

func MessageWeight(diffiuclty int64) int64 {
	return 4 + diffiuclty*10
}

func BlockWeight(diffiuclty int64) int64 {
	return 80 + diffiuclty*200
}

func CalDifficultyHeigth(height int64) int64 {
	var h int64
	if height == 0 {
		h = 0
	} else if height < Consensus.DifficultyInterval {
		h = 0
	} else {
		h = (height - 1) / Consensus.DifficultyInterval * Consensus.DifficultyInterval
	}
	return h
}

func NewDifficut(lastHeight, lastTime, lastWeight, lastDiffiuclty, currentHeight, currentTime, currentWeight int64) (int64, bool, error) {
	if currentHeight-lastHeight != Consensus.DifficultyInterval || currentHeight%Consensus.DifficultyInterval != 0 || lastHeight%Consensus.DifficultyInterval != 0 {
		return 0, false, ErrDifficultHeightErr
	}
	if currentTime < lastTime {
		return 0, false, ErrDifficultHeightErr
	}
	fix := (float64)(currentTime-lastTime) / (float64)(currentHeight-lastHeight) / (float64)(Consensus.ExpectBlockInterval)
	if fix >= 2 {
		fix = 2
	}
	if fix <= 0.5 {
		fix = 0.5
	}
	x := fix * float64(lastDiffiuclty)
	return int64(x), true, nil
}

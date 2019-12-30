package common

import "time"

const (
	MessageSizeLimit = 65536
)

var (
	MessageMaxExpiration = int64(90 * time.Second)
	ChainID              = int32(0)
)

var (
	BlockDifficultyInterval = int64(100)
	MessageDifficultyInterval = int64(1000)
)

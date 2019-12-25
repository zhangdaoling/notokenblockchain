package common

import "time"

const (
	MessageSizeLimit = 65536
)

var (
	MessageMaxExpiration = int64(90 * time.Second)
	ChainID              = int32(0)
)

package types

type ChainState struct {
	lastHeight        int64
	lastTime          int64
	lastDifficulty    int64
	blockDifficulty   int64
	messageDifficulty int64
}

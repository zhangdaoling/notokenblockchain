package chain

import "github.com/zhangdaoling/notokenblockchain/core/types"

type BlockTreeNode struct {
	Block        *types.Block
	Parent       *BlockTreeNode
	Children     map[*BlockTreeNode]bool
	isLinked     bool
	hash         []byte
	difficult    int64
	sumDifficult int64
	ChainState   *types.ChainState
}

func (n *BlockTreeNode) Height() int64 {
	return n.Block.Height()
}

func (n *BlockTreeNode) Hash() []byte {
	return n.hash
}

func (n *BlockTreeNode) ParentHash() []byte {
	return n.Block.ParentHash()
}

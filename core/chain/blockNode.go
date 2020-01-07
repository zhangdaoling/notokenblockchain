package chain

import "github.com/zhangdaoling/notokenblockchain/core/types"

type BlockTreeNode struct {
	hash            []byte
	Block           *types.Block
	Parent          *BlockTreeNode
	Children        map[*BlockTreeNode]bool
	isLinked        bool
	isMain          bool
	ChainState      types.ChainState
	DifficultyState types.DiffiucltyState
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

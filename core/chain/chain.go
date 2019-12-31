package chain

import (
	"github.com/zhangdaoling/simplechain/core/types"
	"github.com/zhangdaoling/simplechain/core/chainstorage"
)

type ChainConfig struct {
	MinCacheBlockNumber int64
}

type Chain struct {
	db *chainstorage.ChainStorage
	//blocks which unlink to the tree, key is the parent hash
	unlinkNode map[string]*BlockTreeNode
	//the header of the main tree. main tree include the main chain
	root *BlockTreeNode
	//no children block
	leaf map[*BlockTreeNode]int64
	//all block in memory
	blockMaps map[string]*BlockTreeNode
	//search cache
	hashToblock   map[string]*BlockTreeNode
	heightToBlock map[int64]*BlockTreeNode
	hashToMessage map[string]*message.Message
}

type BlockTreeNode struct {
	Block      *block.Block
	Parent     *BlockTreeNode
	Children   map[*BlockTreeNode]bool
	ChainState *chainstorage.ChainState
	isLinked   bool
	isMain     bool
}

func(b *BlockTreeNode) buildChildren(blk *block.Block) *BlockTreeNode{
	h1 := b.
	return &BlockTreeNode{
		Block: blk,
		Parent: b,
		Children:make(map[*BlockTreeNode]bool),
		ChainState:b.ChainState,
		isLinked:b.isLinked,
		isMain:b.isMain,
	}
}

func NewChainManager(db *chainstorage.ChainStorage, config *ChainConfig) (*Chain, error) {
	c := Chain{
		db:         db,
		unlinkNode: make(map[string]*BlockTreeNode, 0),
	}

	var startHeight, i int64
	length := db.Length()
	if length > config.MinCacheBlockNumber {
		startHeight = length - config.MinCacheBlockNumber
	}
	blk, err:= db.GetBlockByNumber(startHeight)
	if err != nil {
		return nil, err
	}
	state := db.GetDifficulty(startHeight)
	parent := &BlockTreeNode{
		Block:blk,
		Parent:nil,
		Children:make(map[*BlockTreeNode]bool),
		ChainState: state,
		isLinked: true,
		isMain:true,
	}
	root := parent
	for i = length - 1; i >= startHeight; i-- {
		blk, err := db.GetBlockByNumber(i)
		if err != nil {
			return nil, nil
		}
		blkNode := &BlockTreeNode{
			Block: blk,
			Parent: parent,
			Children:make(map[*BlockTreeNode]bool),
			ChainState:parent.ChainState,
			isLinked:parent.isLinked,
			isMain:parent.isMain,
		}
		parent = blkNode
	}

	return nil, nil
}



func (c *Chain) AddBlock(blk *block.Block) error {
	return nil
}

func (c *Chain) flush() error {
	return nil
}

func (c *Chain) HasBlock(hash []byte) (bool, error) {
	return false, nil
}

func (c *Chain) GetBlockByHash(hash []byte) (*block.Block, error) {
	return nil, nil
}

func (c *Chain) GetBlockByNumber(number int64) (*block.Block, error) {
	return nil, nil
}

func (c *Chain) GetHashByHeight(number int64) ([]byte, error) {
	return nil, nil
}

func (c *Chain) GetHeightByHash(hash []byte) (int64, error) {
	return 0, nil
}

func (c *Chain) HasMessage(hash []byte) (bool, error) {
	return false, nil
}

func (c *ChainManager) GetMessage(hash []byte) (*message.Message, error) {
	return nil, nil
}

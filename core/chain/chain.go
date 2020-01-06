package chain

import (
	"bytes"
	"errors"
	"sync"

	"github.com/zhangdaoling/notokenblockchain/core/chainstorage"
	"github.com/zhangdaoling/notokenblockchain/core/types"
)

var (
	ErrBLockNotExist   = errors.New("not found block")
	ErrMessageNotExist = errors.New("not found message")
)

type ChainConfig struct {
	MinCacheBlockNumber int64
	MaxCacheBlockNumber int64
}

type Chain struct {
	config *ChainConfig
	rw     sync.RWMutex
	db     *chainstorage.ChainStorage
	//blocks which unlink to the tree, key is the parent hash
	parentToUnlinkNode map[string]*BlockTreeNode
	//the root of the main tree. main tree include the main chain
	root *BlockTreeNode
	//the main chain leaf node
	mainChainLeaf *BlockTreeNode
	//no children block
	//leafDifficulty map[*BlockTreeNode]int64
	//all block in memory, include unlinked
	blockMaps map[string]*BlockTreeNode
	//cache main chain for search
	hashToblock   map[string]*BlockTreeNode
	heightToBlock map[int64]*BlockTreeNode
	hashToMessage map[string]*types.Message
}

func NewChainManager(db *chainstorage.ChainStorage, config *ChainConfig) (*Chain, error) {
	c := &Chain{
		config:             config,
		db:                 db,
		parentToUnlinkNode: make(map[string]*BlockTreeNode),
		root:               nil,
		mainChainLeaf:      nil,
		//leafDifficulty: make(map[*BlockTreeNode]int64),
		blockMaps:     make(map[string]*BlockTreeNode),
		hashToblock:   make(map[string]*BlockTreeNode),
		heightToBlock: make(map[int64]*BlockTreeNode),
		hashToMessage: make(map[string]*types.Message),
	}

	blk, err := db.GetBlockByHeight(db.Get)
	if blk == nil {
		return nil, chainstorage.ErrDBBLockNotExist
	}
	if err != nil {
		return nil, err
	}
	state := db.GetState(db.Length())
	err = blk.Verify()
	if err != nil {
		return nil, err
	}
	hash, err := blk.Hash()
	if err != nil {
		return nil, err
	}
	parent := &BlockTreeNode{
		Block:      blk,
		hash:       hash,
		Parent:     nil,
		Children:   make(map[*BlockTreeNode]bool),
		ChainState: state,
		isLinked:   true,
	}
	c.setRoot(parent)
	c.setMainLeaf(parent)

	for i := int64(0); i < config.MinCacheBlockNumber-1; i++ {
		if c.rootHeight() == 0 {
			break
		}
		err := c.cacheRootParent()
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Chain) setRoot(node *BlockTreeNode) {
	c.blockMaps[string(node.Hash())] = node
	return
}

func (c *Chain) setMainLeaf(node *BlockTreeNode) {
	c.mainChainLeaf = node
}

func (c *Chain) AddBlock(blk *types.Block) error {
	hash, err := blk.Hash()
	if err != nil {
		return err
	}
	//the same blk, return
	n, err := c.GetBlockByHash(hash)
	if err != nil {
		return err
	}
	if n != nil {
		return nil
	}
	_, ok := c.blockMaps[string(hash)]
	if ok {
		return nil
	}

	err = blk.Verify()
	if err != nil {
		return err
	}
	parentHash := blk.ParentHash()

	//cache more blk from db if parent is in db
	exist, err := c.dbHasBlock(parentHash)
	if err != nil {
		return err
	}
	if exist {
		parentBlk, err := c.db.GetBlockByHash(parentHash)
		if err != nil {
			return err
		}
		if parentBlk != nil {
			if c.rootHeight()-parentBlk.Height() > c.config.MaxCacheBlockNumber {
				//nothing to do
				return nil
			}
			for {
				if c.rootHeight() == 0 || bytes.Equal(c.root.Hash(), parentHash) {
					break
				}
				err := c.cacheRootParent()
				if err != nil {
					return err
				}
			}
		}
	}

	//add to parentUnlink
	node := &BlockTreeNode{
		Block:      blk,
		hash:       hash,
		Parent:     nil,
		Children:   make(map[*BlockTreeNode]bool),
		ChainState: nil,
		isLinked:   false,
	}
	unlinkNode, ok := c.parentToUnlinkNode[string(hash)]
	if ok {
		unlinkNode = &BlockTreeNode{
			Block:      blk,
			hash:       hash,
			Parent:     nil,
			Children:   make(map[*BlockTreeNode]bool),
			ChainState: nil,
			isLinked:   false,
		}
	}
	unlinkNode.Children[node] = true
	c.blockMaps[string(hash)] = node

	parentNode, ok := c.blockMaps[string(node.ParentHash())]
	if !ok {
		return nil
	}
	if parentNode.isLinked {

	}
	return nil
}

func(c *Chain) link(node *BlockTreeNode){
	if len(node.Children) == 0 {
		return
	}
	for child, _ := range node.Children{
		//todo verify chainstate
		child.ChainState = node.ChainState
		if {

		}
		child.isLinked = true
		c.link(child)
	}

}

func (c *Chain) memoryGetBlock(hash []byte) *BlockTreeNode {
	return c.blockMaps[string(hash)]
}

func (c *Chain) memorySetBlock(node *BlockTreeNode) {
}

func (c *Chain) HasBlock(hash []byte) (bool, error) {
	_, ok := c.blockMaps[string(hash)]
	if ok {
		return true, nil
	}

	return c.dbHasBlock(hash)
}

func (c *Chain) GetBlockByHash(hash []byte) (*types.Block, error) {
	return nil, nil
}

func (c *Chain) GetBlockByNumber(number int64) (*types.Block, error) {
	return nil, nil
}

func (c *Chain) HasMessage(hash []byte) (bool, error) {
	return false, nil
}

func (c *Chain) GetMessage(hash []byte) (*types.Message, error) {
	return nil, nil
}

func (c *Chain) flush() error {
	return nil
}

func (c *Chain) rootHeight() int64 {
	return c.root.Height()
}

func (c *Chain) rootHash() []byte {
	return c.root.Hash()
}

func (c *Chain) rootParentHash() []byte {
	return c.root.ParentHash()
}

func (c *Chain) dbHasBlock(hash []byte) (bool, error) {
	height := c.rootHeight()
	blk, err := c.db.GetBlockByHash(hash)
	if err != nil {
		return false, err
	}

	if blk == nil {
		return false, nil
	}
	if height <= blk.Height() {
		return false, nil
	}
	return true, nil
}

func (c *Chain) dbGetBlock(hash []byte) (*types.Block, error) {
	blk, err := c.db.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	if blk == nil {
		return nil, nil
	}
	height := c.rootHeight()
	if height <= blk.Height() {
		return nil, nil
	}
	return blk, nil
}

func (c *Chain) cacheRootParent() error {
	if c.rootHeight() == 0 {
		return nil
	}
	parentHash := c.rootParentHash()
	n := c.memoryGetBlock(parentHash)
	if n != nil {
		return nil
	}

	blk, err := c.db.GetBlockByHash(parentHash)
	if err != nil {
		return nil
	}
	err = blk.Verify()
	if err != nil {
		return err
	}
	hash, err := blk.Hash()
	if err != nil {
		return err
	}
	if bytes.Equal(hash, parentHash) {
		return
	}
	if c.rootHeight() != blk.Height() {
		return
	}
	node := &BlockTreeNode{
		Block:      blk,
		hash:       hash,
		Parent:     nil,
		Children:   make(map[*BlockTreeNode]bool),
		ChainState: nil,
		isLinked:   true,
		isMain:     true,
	}
	if c.root != nil {
		node.Children[c.root] = true
	}
	c.setRoot(node)
	return nil
}

package chainstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zhangdaoling/notokenblockchain/common"
	"github.com/zhangdaoling/notokenblockchain/core/consensus"
	"github.com/zhangdaoling/notokenblockchain/core/types"
	"github.com/zhangdaoling/notokenblockchain/db/kv"
	"github.com/zhangdaoling/notokenblockchain/pb"
)

var (
	ErrDBBLockNotExist      = errors.New("db not found block")
	ErrDBMessageNotExist    = errors.New("db not found message")
	ErrDBHeightDiffNotExist = errors.New("db not found difficulty height")
)

type ChainStorage struct {
	db         *kv.Storage
	rw         sync.RWMutex
	chainState *types.ChainState
}

var (
	difficultyStatePrefix = []byte("ds")
	chainStatePrefix      = []byte("cs")
)

var (
	blockHeightPrefix  = []byte("bh")
	blockHashPrefix    = []byte("b")
	messageHashPrefix  = []byte("mh")
	messageIndexPrefix = []byte("mi")
)

func NewBlockChain(path string) (*ChainStorage, error) {
	levelDB, err := kv.NewStorage(path, kv.LevelDBStorage)
	if err != nil {
		return nil, fmt.Errorf("fail to init blockchaindb, %v", err)
	}
	state := consensus.GeneisChainState
	ptr := &state
	ok, err := levelDB.Has(chainStatePrefix)
	if err != nil {
		return nil, err
	}
	if ok {
		stateByte, err := levelDB.Get(chainStatePrefix)
		if err != nil {
			return nil, err
		}
		err = ptr.DBDecode(stateByte)
		if err != nil {
			return nil, err
		}
	} else {
		difficultyState := &types.DiffiucltyState{
			Height:     consensus.GeneisBlock.Height(),
			Time:       consensus.GeneisBlock.Time(),
			Weight:     0,
			Difficulty: 1,
		}
		difficultyByte := difficultyState.DBBytes()
		err = levelDB.Put(append(difficultyStatePrefix, common.Int64ToBytes(0)...), difficultyByte)
		if err != nil {
			return nil, err
		}

		stateByte := ptr.DBBytes()
		err := levelDB.Put(chainStatePrefix, stateByte)
		if err != nil {
			return nil, err
		}
	}
	s := &ChainStorage{
		db:         levelDB,
		chainState: ptr,
	}
	err = s.CheckLength()
	if err != nil {
		return nil, err
	}
	return s, err
}

func (c *ChainStorage) Close() {
	c.db.Close()
}

func (c *ChainStorage) PopBlock() (*types.Block, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	err := c.db.Begin()
	if err != nil {
		return nil, errors.New("fail to begin batch")
	}

	blk, state, err := c.popBlock()
	if err != nil {
		c.db.RollBack()
		return nil, err
	}
	err = c.db.Commit()
	if err != nil {
		c.db.RollBack()
		return nil, err
	}
	c.chainState = state
	return blk, nil
}

func (c *ChainStorage) AddBlock(blk *types.Block) error {
	c.rw.RLock()
	defer c.rw.RUnlock()

	err := c.db.Begin()
	if err != nil {
		return errors.New("fail to begin batch")
	}

	state, err := c.addBlock(blk, c.db)
	if err != nil {
		c.db.RollBack()
		return err
	}
	err = c.db.Commit()
	if err != nil {
		c.db.RollBack()
		return err
	}
	c.chainState = state
	return nil
}

func (c *ChainStorage) State() types.ChainState {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return *c.chainState
}

func (c *ChainStorage) HasBlock(hash []byte) (bool, error) {
	return c.db.Has(append(blockHashPrefix, hash...))
}

func (c *ChainStorage) GetBlockByHash(hash []byte) (*types.Block, error) {
	blockByte, err := c.getBlockByteByHash(hash)
	if err != nil {
		return nil, err
	}
	if len(blockByte) == 0 {
		return nil, nil
	}
	blk := &types.Block{}
	err = blk.DBDecode(blockByte)
	if err != nil {
		return nil, err
	}
	if blk.MessageHashes() != nil {
		msgs := make([]*pb.Message, len(blk.MessageHashes()))
		for _, hash := range blk.MessageHashes() {
			msg, err := c.GetMessage(hash)
			if err != nil {
				return nil, fmt.Errorf("miss the msg, msg hash: %s", hash)
			}
			msgs = append(msgs, msg.PBMessage)
		}
		blk.SetMessages(msgs)
	}
	return blk, nil
}

func (c *ChainStorage) GetBlockByHeight(number int64) (*types.Block, error) {
	hash, err := c.db.Get(append(blockHashPrefix, common.Int64ToBytes(number)...))
	if err != nil {
		return nil, errors.New("fail to get hash by number")
	}
	if len(hash) == 0 {
		return nil, nil
	}
	return c.GetBlockByHash(hash)
}

func (c *ChainStorage) HasMessage(hash []byte) (bool, error) {
	return c.db.Has(append(messageIndexPrefix, hash...))
}

func (c *ChainStorage) GetMessage(hash []byte) (*types.Message, error) {
	msg := &types.Message{}
	msgIndex, err := c.db.Get(append(messageIndexPrefix, hash...))
	if err != nil {
		return nil, fmt.Errorf("failed to Get the msg: %v", err)
	}
	if len(msgIndex) == 0 {
		return nil, nil
	}

	msgData, err := c.db.Get(append(messageHashPrefix, msgIndex...))
	if err != nil {
		return nil, fmt.Errorf("failed to Get the msg: %v", err)
	}
	if len(msgData) == 0 {
		return nil, nil
	}

	err = msg.DBDecode(msgData)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode the msg: %v", err)
	}
	return msg, nil
}

//todo
func (c *ChainStorage) GetDifficulty(height int64) (*types.DiffiucltyState, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.getDifficulty(height)
}

func (c *ChainStorage) getDifficulty(height int64) (*types.DiffiucltyState, error) {
	h := consensus.CalDifficultyHeigth(height)
	state := &types.DiffiucltyState{}
	stateByte, err := c.db.Get(append(chainStatePrefix, common.Int64ToBytes(h)...))
	if err != nil {
		return nil, err
	}
	if len(stateByte) == 0 {
		return nil, nil
	}
	err = state.DBDecode(stateByte)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (c *ChainStorage) popBlock() (*types.Block, *types.ChainState, error) {
	height := c.chainState.Length - 1
	blk, err := c.GetBlockByHeight(height)
	if err != nil {
		return nil, nil, err
	}
	if blk == nil {
		return nil, nil, nil
	}

	if blk.Height()%consensus.Consensus.DifficultyInterval == 0 {
		err = c.db.Delete(append(difficultyStatePrefix, common.Int64ToBytes(blk.Height())...))
		if err != nil {
			return nil, nil, err
		}
	}
	//delete block
	hash, err := blk.Hash()
	if err != nil {
		return nil, nil, err
	}
	c.db.Delete(append(blockHeightPrefix, common.Int64ToBytes(blk.Height())...))
	c.db.Delete(append(blockHashPrefix, hash...))
	for _, m := range blk.Messages() {
		msg := types.ToMessage(m)
		msgHash, err := msg.Hash()
		if err != nil {
			return nil, nil, err
		}
		c.db.Delete(append(messageHashPrefix, msgHash...))
		c.db.Delete(append(messageIndexPrefix, append(hash, msgHash...)...))
	}

	state := &types.ChainState{
		Length:       c.chainState.Length - 1,
		TotalWeight:  c.chainState.TotalMessage - blk.Weight(),
		TotalMessage: c.chainState.TotalMessage - int64(len(blk.Messages())),
	}
	stateByte := state.DBBytes()
	err = c.db.Put(chainStatePrefix, stateByte)
	if err != nil {
		return nil, nil, err
	}
	return blk, state, nil
}

func (c *ChainStorage) addBlock(blk *types.Block, db *kv.Storage) (state *types.ChainState, err error) {
	hash, err := blk.Hash()
	if err != nil {
		return nil, err
	}
	blockByte, err := blk.DBBytes()
	if err != nil {
		return nil, err
	}
	number := blk.Height()
	c.db.Put(append(blockHeightPrefix, common.Int64ToBytes(number)...), hash)
	c.db.Put(append(blockHashPrefix, hash...), blockByte)
	for _, m := range blk.Messages() {
		msg := types.ToMessage(m)
		msgHash, err := msg.Hash()
		if err != nil {
			return nil, err
		}
		msgBytes, err := msg.DBBytes()
		if err != nil {
			return nil, err
		}
		c.db.Put(append(messageHashPrefix, msgHash...), append(hash, msgHash...))
		c.db.Put(append(messageIndexPrefix, append(hash, msgHash...)...), msgBytes)
	}
	//diffucult state
	if blk.Height()%consensus.Consensus.DifficultyInterval == 0 {
		lastState, err := c.getDifficulty(blk.Height())
		if err != nil {
			return nil, err
		}
		if lastState == nil {
			return nil, ErrDBHeightDiffNotExist
		}
		difficulty, bool, err := consensus.NewDifficut(lastState.Height, lastState.Time, lastState.Weight, lastState.Difficulty, blk.Height(), blk.Time(), blk.Weight())
		if err != nil {
			return nil, err
		}
		if bool {
			difficultyState := types.DiffiucltyState{
				Height:     blk.Height(),
				Time:       blk.Time(),
				Weight:     c.chainState.TotalWeight + blk.Weight(),
				Difficulty: difficulty,
			}
			difficultyByte := difficultyState.DBBytes()
			err = c.db.Put(append(difficultyStatePrefix, common.Int64ToBytes(blk.Height())...), difficultyByte)
			if err != nil {
				return nil, err
			}
		}
	}
	//chain state
	state = &types.ChainState{
		Length:       c.chainState.Length + 1,
		TotalWeight:  c.chainState.TotalMessage + blk.Weight(),
		TotalMessage: c.chainState.TotalMessage + int64(len(blk.Messages())),
	}
	stateByte := state.DBBytes()
	err = c.db.Put(chainStatePrefix, stateByte)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (c *ChainStorage) getBlockByteByHash(hash []byte) ([]byte, error) {
	blockByte, err := c.db.Get(append(blockHashPrefix, hash...))
	if err != nil {
		return nil, errors.New("fail to get block byte by hash")
	}
	return blockByte, nil
}

// CheckLength is check length of block in database
func (c *ChainStorage) CheckLength() error {
	c.rw.RLock()
	defer c.rw.RUnlock()
	for i := c.chainState.Length; i > 0; i-- {
		n, err := c.GetBlockByHeight(i - 1)
		if err != nil {
			return err
		}
		if n == nil {
			return fmt.Errorf("not found block: %d", i-1)
		}
	}
	return nil
}

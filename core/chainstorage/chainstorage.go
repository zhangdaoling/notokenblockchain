package chainstorage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zhangdaoling/notokenblockchain/core/consensus"
	"strconv"
	"sync"

	"github.com/zhangdaoling/notokenblockchain/common"
	"github.com/zhangdaoling/notokenblockchain/core/types"
	"github.com/zhangdaoling/notokenblockchain/db/kv"
	"github.com/zhangdaoling/notokenblockchain/pb"
)

var (
	ErrDBBLockNotExist   = errors.New("db not found block")
	ErrDBMessageNotExist = errors.New("db not found message")
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
	if blk.PBBlock.MessageHashes != nil {
		blk.PBBlock.Messages = make([]*pb.Message, len(blk.PBBlock.MessageHashes))
		for _, hash := range blk.PBBlock.MessageHashes {
			msg, err := c.GetMessage(hash)
			if err != nil {
				return nil, fmt.Errorf("miss the msg, msg hash: %s", hash)
			}
			blk.PBBlock.Messages = append(blk.PBBlock.Messages, msg.PBMessage)
		}
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
	err = state.DBDecode(stateByte)
	if err != nil {
		return nil, err
	}
	return state, nil
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
	number := blk.PBBlock.Head.Height
	c.db.Put(append(blockHeightPrefix, common.Int64ToBytes(number)...), hash)
	c.db.Put(append(blockHashPrefix, hash...), blockByte)
	for _, m := range blk.PBBlock.Messages {
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
	if blk.PBBlock.Head.Height%consensus.Consensus.DifficultyInterval == 0 {
		lastDifficultyState, err := c.getDifficulty(blk.PBBlock.Head.Height)
		if err != nil{
			return nil, err
		}
		currentDifficultyState


	}
	state = &types.ChainState{
		Length:       c.chainState.Length + 1,
		TotalWeight:  c.chainState.TotalMessage + blk.PBBlock.Head.Difficulty,
		TotalMessage: c.chainState.TotalMessage + int64(len(blk.PBBlock.Messages)),
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

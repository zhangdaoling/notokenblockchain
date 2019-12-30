package chainstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zhangdaoling/simplechain/common"
	"github.com/zhangdaoling/simplechain/core/block"
	"github.com/zhangdaoling/simplechain/core/message"
	"github.com/zhangdaoling/simplechain/db/kv"
	"github.com/zhangdaoling/simplechain/pb"
)

var (
	ErrDBNotFoundMessage = errors.New("db not found message")
)

type ChainStorage struct {
	db           *kv.Storage
	rw           sync.RWMutex
	length       int64
	weight       int64
	messageTotal int64
}

type ChainState struct {
	lastHeight        int64
	lastTime          int64
	lastDifficulty    int64
	blockDifficulty   int64
	messageDifficulty int64
}

//todo, more state
var (
	blockLengthPrefixes     = []byte("bl")
	messageTotalPrefix      = []byte("mt")
	blockDifficultyPrefix   = []byte("bd")
	messageDifficultyPrefix = []byte("md")
)

//info
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
	var length, weight, messageTotal int64
	ok, err := levelDB.Has(blockLengthPrefixes)
	if err != nil {
		return nil, fmt.Errorf("fail to check has(lengthPrefixes), %v", err)
	}
	if ok {
		lengthByte, err := levelDB.Get(blockLengthPrefixes)
		if err != nil || len(lengthByte) == 0 {
			return nil, errors.New("fail to get lengthPrefixes")
		}
		length = common.BytesToInt64(lengthByte)
		weightByte, err := levelDB.Get(blockLengthPrefixes)
		if err != nil || len(weightByte) == 0 {
			return nil, errors.New("fail to get difficultyPrefix")
		}
		weight = common.BytesToInt64(weightByte)
		messageTotalByte, err := levelDB.Get(messageTotalPrefix)
		if err != nil || len(messageTotalByte) == 0 {
			return nil, errors.New("fail to get msg total")
		}
		messageTotal = common.BytesToInt64(messageTotalByte)
	} else {
		lengthByte := common.Int64ToBytes(0)
		if err := levelDB.Put(blockLengthPrefixes, lengthByte); err != nil {
			return nil, errors.New("fail to put lengthPrefixes")
		}
		weightByte := common.Int64ToBytes(weight)
		if err := levelDB.Put(blockLengthPrefixes, weightByte); err != nil {
			return nil, errors.New("fail to put difficultyPrefix")
		}
		messageTotalByte := common.Int64ToBytes(0)
		if err := levelDB.Put(messageTotalPrefix, messageTotalByte); err != nil {
			return nil, errors.New("fail to put msg total")
		}
	}
	s := &ChainStorage{
		db:           levelDB,
		length:       length,
		messageTotal: messageTotal,
	}
	s.CheckLength()
	return s, err
}

func (c *ChainStorage) Close() {
	c.db.Close()
}

func (c *ChainStorage) AddBlock(blk *block.Block) error {
	c.rw.RLock()
	defer c.rw.RUnlock()

	err := c.db.Begin()
	if err != nil {
		return errors.New("fail to begin batch")
	}

	err = c.addBlock(blk, c.db)
	if err != nil {
		c.db.RollBack()
		return err
	}
	err = c.db.Commit()
	if err != nil {
		c.db.RollBack()
		return err
	}
	return nil
}

func (c *ChainStorage) Length() int64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.length
}

func (c *ChainStorage) Weight() int64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.weight
}

func (c *ChainStorage) MessageTotal() int64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.messageTotal
}

func (c *ChainStorage) HasBlock(hash []byte) (bool, error) {
	return c.db.Has(append(blockHashPrefix, hash...))
}

func (c *ChainStorage) GetBlockByHash(hash []byte) (*block.Block, error) {
	blockByte, err := c.getBlockByteByHash(hash)
	if err != nil {
		return nil, err
	}
	blk := &block.Block{}
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

func (c *ChainStorage) GetBlockByNumber(number int64) (*block.Block, error) {
	hash, err := c.GetHashByHeight(number)
	if err != nil {
		return nil, err
	}
	return c.GetBlockByHash(hash)
}

func (c *ChainStorage) GetHashByHeight(number int64) ([]byte, error) {
	hash, err := c.db.Get(append(blockHashPrefix, common.Int64ToBytes(number)...))
	if err != nil || len(hash) == 0 {
		return nil, errors.New("fail to get hash by number")
	}
	return hash, nil
}

func (c *ChainStorage) GetHeightByHash(hash []byte) (int64, error) {
	blockByte, err := c.getBlockByteByHash(hash)
	if err != nil {
		return 0, err
	}
	blk := &block.Block{}
	err = blk.DBDecode(blockByte)
	if err != nil {
		return 0, err
	}
	return blk.PBBlock.Head.Height, nil
}

func (c *ChainStorage) HasMessage(hash []byte) (bool, error) {
	return c.db.Has(append(messageIndexPrefix, hash...))
}

func (c *ChainStorage) GetMessage(hash []byte) (*message.Message, error) {
	msg := &message.Message{}
	msgIndex, err := c.db.Get(append(messageIndexPrefix, hash...))
	if err != nil {
		return nil, fmt.Errorf("failed to Get the msg: %v", err)
	}
	msgData, err := c.db.Get(append(messageHashPrefix, msgIndex...))
	if err != nil {
		return nil, fmt.Errorf("failed to Get the msg: %v", err)
	}
	if len(msgData) == 0 {
		return nil, fmt.Errorf("failed to Get the msg: not found")
	}
	err = msg.DBDecode(msgData)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode the msg: %v", err)
	}
	return msg, nil
}

func (c *ChainStorage) GetDifficulty(height int64) {
	if height > c.length-1 {
		height = c.length - 1
	}
	h := height / common.BlockDifficultyInterval
	if h == 0 {

	}
}

func (c *ChainStorage) addBlock(blk *block.Block, db *kv.Storage) error {
	hash, err := blk.Hash()
	if err != nil {
		return err
	}
	blockByte, err := blk.DBBytes()
	if err != nil {
		return err
	}
	number := blk.PBBlock.Head.Height
	messageTotal := c.MessageTotal()
	c.db.Put(append(blockHeightPrefix, common.Int64ToBytes(number)...), hash)
	c.db.Put(append(blockHashPrefix, hash...), blockByte)
	c.db.Put(blockLengthPrefixes, common.Int64ToBytes(number+1))
	c.db.Put(messageTotalPrefix, common.Int64ToBytes(messageTotal+int64(len(blk.PBBlock.Messages))))
	for _, m := range blk.PBBlock.Messages {
		msg := message.ToMessage(m)
		msgHash, err := msg.Hash()
		if err != nil {
			return err
		}
		msgBytes, err := msg.DBBytes()
		if err != nil {
			return err
		}
		c.db.Put(append(messageHashPrefix, msgHash...), append(hash, msgHash...))
		c.db.Put(append(messageIndexPrefix, append(hash, msgHash...)...), msgBytes)
	}
	c.setLength(number + 1)
	c.setMessageTotal(messageTotal + int64(len(blk.PBBlock.Messages)))
	return nil
}

func (c *ChainStorage) getBlockByteByHash(hash []byte) ([]byte, error) {
	blockByte, err := c.db.Get(append(blockHashPrefix, hash...))
	if err != nil || len(blockByte) == 0 {
		return nil, errors.New("fail to get block byte by hash")
	}
	return blockByte, nil
}

func (c *ChainStorage) setLength(i int64) {
	c.length = i
}

func (c *ChainStorage) setWeight(i int64) {
	c.weight = i
}

func (c *ChainStorage) setMessageTotal(i int64) {
	c.messageTotal = i
}

// CheckLength is check length of block in database
func (c *ChainStorage) CheckLength() {
	for i := c.Length(); i > 0; i-- {
		_, err := c.GetBlockByNumber(i - 1)
		if err != nil {
			fmt.Println("fail to get the block")
		} else {
			c.db.Put(blockLengthPrefixes, common.Int64ToBytes(i))
			c.setLength(i)
			break
		}
	}
}

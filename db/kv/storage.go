package kv

import (
	"github.com/zhangdaoling/notokenblockchain/db/kv/leveldb"
)

type StorageType uint8

const (
	_ StorageType = iota
	LevelDBStorage
)

type StorageBackend interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Begin() error
	Commit() error
	RollBack() error
	Close() error
	//Size() (int64, error)
	//NewIteratorByPrefix(prefix []byte) interface{}
}

var _ StorageBackend = &leveldb.DB{}

type Storage struct {
	StorageBackend
}

func NewStorage(path string, t StorageType) (*Storage, error) {
	switch t {
	case LevelDBStorage:
		sb, err := leveldb.NewDB(path)
		if err != nil {
			return nil, err
		}
		return &Storage{StorageBackend: sb}, nil
	default:
		sb, err := leveldb.NewDB(path)
		if err != nil {
			return nil, err
		}
		return &Storage{StorageBackend: sb}, nil
	}
}

/*
func (s *Storage) NewIteratorByPrefix(prefix []byte) *Iterator {
	ib := s.StorageBackend.NewIteratorByPrefix(prefix).(IteratorBackend)
	return &Iterator{
		IteratorBackend: ib,
	}
}

type IteratorBackend interface {
	Next() bool
	Key() []byte
	Value() []byte
	Error() error
	Release()
}

// Iterator is the storage iterator
type Iterator struct {
	IteratorBackend
}
*/

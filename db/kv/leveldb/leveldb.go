package leveldb

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

func NewDB(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		db:    db,
		batch: nil,
	}, nil
}

func (d *DB) Get(key []byte) ([]byte, error) {
	value, err := d.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return []byte{}, nil
	}

	return value, err
}

func (d *DB) Has(key []byte) (bool, error) {
	return d.db.Has(key, nil)
}

func (d *DB) Put(key []byte, value []byte) error {
	if d.batch == nil {
		return d.db.Put(key, value, nil)
	}
	d.batch.Put(key, value)
	return nil
}

func (d *DB) Delete(key []byte) error {
	if d.batch == nil {
		return d.db.Delete(key, nil)
	}
	d.batch.Delete(key)
	return nil
}

func (d *DB) Begin() error {
	if d.batch != nil {
		return fmt.Errorf("not support nested batch write")
	}
	d.batch = new(leveldb.Batch)
	return nil
}

func (d *DB) Commit() error {
	if d.batch == nil {
		return fmt.Errorf("no batch write to commit")
	}
	err := d.db.Write(d.batch, nil)
	if err != nil {
		return err
	}
	d.batch = nil
	return nil
}

func (d *DB) RollBack() error {
	d.batch = nil
	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

/*
func (d *DB) Size() (int64, error) {
	stats := &leveldb.DBStats{}
	if err := d.db.Stats(stats); err != nil {
		return 0, err
	}
	total := int64(0)
	for _, size := range stats.LevelSizes {
		total += size
	}
	return total, nil
}

func (d *DB) NewIteratorByPrefix(prefix []byte) interface{} {
	iter := d.db.NewIterator(util.BytesPrefix(prefix), nil)
	return &Iter{
		iter: iter,
	}
}

type Iter struct {
	iter iterator.Iterator
}

func (i *Iter) Next() bool {
	return i.iter.Next()
}

func (i *Iter) Key() []byte {
	return i.iter.Key()
}

func (i *Iter) Value() []byte {
	return i.iter.Value()
}

func (i *Iter) Error() error {
	return i.iter.Error()
}

func (i *Iter) Release() {
	i.iter.Release()
}
*/

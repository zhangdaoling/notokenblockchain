package types

import (
	"encoding/binary"
	"fmt"
)

type ChainState struct {
	Difficulty   int64
	Length       int64
	TotalWeight  int64
	TotalMessage int64
}

func (c *ChainState) DBBytes() []byte {
	b := make([]byte, 32)
	binary.BigEndian.PutUint64(b[0:8], uint64(c.Difficulty))
	binary.BigEndian.PutUint64(b[8:16], uint64(c.Length))
	binary.BigEndian.PutUint64(b[16:24], uint64(c.TotalWeight))
	binary.BigEndian.PutUint64(b[24:], uint64(c.TotalMessage))
	return b
}

func (c *ChainState) DBDecode(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("sdfsdfsdf")
	}
	c.Difficulty = int64(binary.BigEndian.Uint16(b[0:8]))
	c.Length = int64(binary.BigEndian.Uint16(b[8:16]))
	c.TotalWeight = int64(binary.BigEndian.Uint16(b[16:24]))
	c.TotalMessage = int64(binary.BigEndian.Uint16(b[24:]))
	return nil
}

type DiffiucltyState struct {
	Height     int64
	Time       int64
	Weight     int64
	Difficulty int64
}

func (d *DiffiucltyState) DBBytes() []byte {
	b := make([]byte, 32)
	binary.BigEndian.PutUint64(b[0:8], uint64(d.Height))
	binary.BigEndian.PutUint64(b[8:16], uint64(d.Time))
	binary.BigEndian.PutUint64(b[16:24], uint64(d.Weight))
	binary.BigEndian.PutUint64(b[24:0], uint64(d.Difficulty))
	return b
}

func (d *DiffiucltyState) DBDecode(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("sdfsd")
	}
	d.Height = int64(binary.BigEndian.Uint16(b[0:8]))
	d.Time = int64(binary.BigEndian.Uint16(b[8:16]))
	d.Weight = int64(binary.BigEndian.Uint16(b[16:24]))
	d.Difficulty = int64(binary.BigEndian.Uint16(b[24:]))
	return nil
}

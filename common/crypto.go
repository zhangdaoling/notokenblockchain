package common

import (
	"encoding/hex"
	"fmt"
	"encoding/binary"
	"hash/crc32"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Sha3 ...
func Sha3(raw []byte) (b []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("sha3 panic. err=%v", err)
			b = nil
		}
	}()
	b = sha3.Sum256(raw)[:]
	return
}

func NewSha3(raw []byte) ([]byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("sha3 panic. err=%v", err)
		}
	}()
	return sha3.Sum256(raw)[:]
}
// Ripemd160 ...
func Ripemd160(raw []byte) []byte {
	h := ripemd160.New()
	h.Write(raw)
	return h.Sum(nil)
}

// Base58Encode ...
func Base58Encode(raw []byte) string {
	return base58.Encode(raw)
}

// Base58Decode ...
func Base58Decode(s string) []byte {
	return base58.Decode(s)
}

// Parity ...
func Parity(bit []byte) []byte {
	crc32q := crc32.MakeTable(crc32.Koopman)
	crc := crc32.Checksum(bit, crc32q)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, crc)
	return bs
}

// ToHex ...
func ToHex(data []byte) string {
	return hex.EncodeToString(data)
}

// ParseHex ...
func ParseHex(s string) []byte {
	d, err := hex.DecodeString(s)
	if err != nil {
		println(err)
		return nil
	}
	return d
}
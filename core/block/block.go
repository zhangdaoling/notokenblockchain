package block

import (
	"github.com/zhangdaoling/simplechain/common"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/simplechain/pb"
)

type Block struct {
	BaseHash []byte
	FullHash []byte
	PbBLock  *pb.Block
}

func (b *Block) String() string {
	return ""
}

func (b *Block) Encode() ([]byte, error) {
	return b.FullByte()
}

func (b *Block) Decode(buf []byte) error {
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(buf, pbBlock)
	if err != nil {
		return err
	}
	b.PbBLock = pbBlock
	return nil
}

func (b *Block) Hash() ([]byte, error) {
	return b.CalculateBaseHash()
}

func (b *Block) BaseByte() ([]byte, error) {
	buf, err := proto.Marshal(b.PbBLock.Head)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *Block) FullByte() ([]byte, error) {
	buf, err := proto.Marshal(b.PbBLock)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//header hash
func (b *Block) CalculateBaseHash() (buf []byte, err error) {
	if b.BaseHash == nil {
		tmp, err := b.BaseByte()
		if err != nil {
			return nil, nil
		}
		buf, err = common.Sha3(tmp)
		if err != nil{
			return nil, err
		}
	}
	return
}

func (b *Block) CalculateFullHash() (buf []byte, err error) {
	if b.FullHash == nil {
		tmp, err := b.FullByte()
		if err != nil {
			return nil, err
		}
		buf, err = common.Sha3(tmp)
		if err != nil{
			return nil, err
		}
	}
	return
}

func (b *Block) CalculateTxMerkleHash() []byte {
	m := merkletree.MerkleTree{}
	hashes := make([][]byte, 0, len(b.Txs))
	for _, tx := range b.Txs {
		hashes = append(hashes, tx.Hash())
	}
	m.Build(hashes)
	return m.RootHash()
}

func (b *Block) Sign() {
	return
}

func (b *Block) Verify() error {
	return nil
}

package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/simplechain/common"
	"github.com/zhangdaoling/simplechain/core/merkletree"
	"github.com/zhangdaoling/simplechain/pb"
)

type Block struct {
	PBBlock  *pb.Block
}

func ToBLock(b *pb.Block) *Block {
	return &Block{
		PBBlock: b,
	}
}

func ToPbBlock(b *Block) *pb.Block {
	return b.PBBlock
}

func (b *Block) String() string {
	return ""
}

func (b *Block) P2PBytes () ([]byte, error){
	return proto.Marshal(b.PBBlock)
}

func (b *Block) P2PDecode (buf []byte) error{
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(buf, pbBlock)
	if err != nil {
		return err
	}
	b.PBBlock = pbBlock
	return nil
}

func (b *Block) DBBytes () ([]byte, error){
	tmp := b.PBBlock.Messages
	b.PBBlock.Messages = nil
	buf, err := proto.Marshal(b.PBBlock)
	b.PBBlock.Messages = tmp
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *Block) DBDecode (buf []byte) error{
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(buf, pbBlock)
	if err != nil {
		return err
	}
	b.PBBlock = pbBlock
	return nil
}

func (b *Block) HashBytes () ([]byte, error) {
	return proto.Marshal(b.PBBlock.Head)
}

//block hash; just block head
func (b *Block) Hash() ([]byte, error) {
	baseBytes, err := b.HashBytes()
	if err != nil {
		return nil, err
	}
	buf, err := common.Sha3(baseBytes)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *Block) CalculateTxMerkleHash() ([]byte, error) {
	tree := pb.MerkleTree{}
	hashes := make([][]byte, 0, len(b.PBBlock.Messages))
	for _, msg := range b.PBBlock.Messages {
		m := Message{
			PBMessage: msg,
		}
		hash, err := m.Hash()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}
	merkletree.Build(&tree, hashes)
	return merkletree.RootHash(&tree), nil
}

//to do
func (b *Block) Sign() {
	return
}

//to do
//sort msg by hash
func (b *Block) Verify() error {
	return nil
}

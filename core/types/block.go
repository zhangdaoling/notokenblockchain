package types

import (
	"bytes"
	"errors"
	"github.com/zhangdaoling/notokenblockchain/crypto"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/notokenblockchain/common"
	"github.com/zhangdaoling/notokenblockchain/core/consensus"
	"github.com/zhangdaoling/notokenblockchain/core/merkletree"
	"github.com/zhangdaoling/notokenblockchain/pb"
)

var (
	ErrBlockMessagHashesNotFit = errors.New("block messageHashes no fit message hash")
	ErrBlockMessagHashesWrong  = errors.New("block messageHashes no equal message hash")
	ErrBlockMerlkHashWrong     = errors.New("block merlkHash not right")
	ErrBlockChainIDWrong       = errors.New("block chainid not right")
	ErrBlockSignWrong          = errors.New("block sign verify failed")
	ErrBlockPublicKeyWrong     = errors.New("block public key not right")
)

type Block struct {
	hash    []byte
	pbBlock *pb.Block
}

func (b *Block) Height() int64 {
	return b.pbBlock.Head.Height
}

func (b *Block) Time() int64 {
	return b.pbBlock.Head.Time
}

func (b *Block) Weight() int64 {
	return b.pbBlock.Head.Weight
}

func (b *Block) ChainID() int64 {
	return b.pbBlock.Head.ChainId
}

func (b *Block) ParentHash() []byte {
	return b.pbBlock.Head.ParentHash
}

func (b *Block) PublicKey() []byte {
	return b.pbBlock.Head.PublicKey
}

func (b *Block) Messages() []*pb.Message {
	return b.pbBlock.Messages
}

func (b *Block) MessageHashes() [][]byte {
	return b.pbBlock.MessageHashes
}

func (b *Block) SetMessages(m []*pb.Message) {
	b.pbBlock.Messages = m
}

func ToBLock(b *pb.Block) *Block {
	return &Block{
		pbBlock: b,
	}
}

func ToPbBlock(b *Block) *pb.Block {
	return b.pbBlock
}

func (b *Block) String() string {
	return ""
}

func (b *Block) P2PBytes() ([]byte, error) {
	return proto.Marshal(b.pbBlock)
}

func (b *Block) P2PDecode(buf []byte) error {
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(buf, pbBlock)
	if err != nil {
		return err
	}
	b.pbBlock = pbBlock
	return nil
}

func (b *Block) DBBytes() ([]byte, error) {
	tmp := b.pbBlock.Messages
	b.pbBlock.Messages = nil
	buf, err := proto.Marshal(b.pbBlock)
	b.pbBlock.Messages = tmp
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *Block) DBDecode(buf []byte) error {
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(buf, pbBlock)
	if err != nil {
		return err
	}
	b.pbBlock = pbBlock
	return nil
}

func (b *Block) HashBytes() ([]byte, error) {
	return proto.Marshal(b.pbBlock.Head)
}

//block hash; just block head
func (b *Block) Hash() ([]byte, error) {
	if b.hash != nil {
		return b.hash, nil
	}
	baseBytes, err := b.HashBytes()
	if err != nil {
		return nil, err
	}
	return common.Sha3(baseBytes), nil
}

func (b *Block) CalculateTxMerkleHash() ([]byte, error) {
	tree := pb.MerkleTree{}
	hashes := make([][]byte, 0, len(b.pbBlock.Messages))
	if len(b.pbBlock.MessageHashes) != len(b.pbBlock.Messages) {
		return nil, ErrBlockMessagHashesNotFit
	}

	for i, msg := range b.pbBlock.Messages {
		m := Message{
			PBMessage: msg,
		}
		hash, err := m.Hash()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(b.pbBlock.MessageHashes[i], hash) {
			return nil, ErrBlockMessagHashesWrong
		}
		hashes = append(hashes, hash)
	}
	merkletree.Build(&tree, hashes)
	return merkletree.RootHash(&tree), nil
}

func (b *Block) Sign(privateKey []byte) error {
	alg := crypto.Ed25519
	publicKey, err := alg.GeneratePublicKey(privateKey)
	if err != nil {
		return err
	}
	if !bytes.Equal(b.pbBlock.Head.PublicKey, publicKey) {
		return ErrBlockPublicKeyWrong

	}
	err = b.verifyData()
	if err != nil {
		return err
	}
	hash, err := b.Hash()
	if err != nil {
		return err
	}
	sig, err := alg.Sign(hash, privateKey)
	if err != nil {
		return err
	}
	b.pbBlock.Signature.Algorithm = int32(alg)
	b.pbBlock.Signature.Sig = sig
	b.pbBlock.Head.PublicKey = publicKey
	return nil
}

func (b *Block) Verify() error {
	err := b.verifyData()
	if err != nil {
		return err
	}
	return b.VerifySign()
}

func (b *Block) verifyData() error {
	//check message
	for _, msg := range b.pbBlock.Messages {
		m := ToMessage(msg)
		err := m.Verify()
		if err != nil {
		}
		return nil
	}
	//check message hash
	if len(b.pbBlock.MessageHashes) != len(b.pbBlock.Messages) {
		return ErrBlockMessagHashesNotFit
	}

	//todo calculate and check difficulty

	//check merkle
	merkleTree, err := b.CalculateTxMerkleHash()
	if err != nil {
		return err
	}
	if !bytes.Equal(b.pbBlock.Head.MessageMerkleHash, merkleTree) {
		return ErrBlockMerlkHashWrong
	}
	if b.pbBlock.Head.ChainId != consensus.Consensus.ChainID {
		return ErrBlockChainIDWrong
	}
	//check block hash and cache
	baseBytes, err := b.HashBytes()
	if err != nil {
		return err
	}
	b.hash = common.Sha3(baseBytes)

	return nil
}

func (b *Block) VerifySign() error {
	alg := crypto.Algorithm(b.pbBlock.Signature.Algorithm)
	if !alg.Verify(b.hash, b.pbBlock.Head.PublicKey, b.pbBlock.Signature.Sig) {
		return ErrBlockSignWrong
	}
	return nil
}

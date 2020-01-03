package types

import (
	"bytes"
	"errors"
	"github.com/zhangdaoling/simplechain/crypto"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/simplechain/common"
	"github.com/zhangdaoling/simplechain/core/consensus"
	"github.com/zhangdaoling/simplechain/core/merkletree"
	"github.com/zhangdaoling/simplechain/pb"
)


var (
	ErrBlockMessagHashesNotFit = errors.New("block messageHashes no fit message hash")
	ErrBlockMessagHashesWrong = errors.New("block messageHashes no equal message hash")
	ErrBlockMerlkHashWrong = errors.New("block merlkHash not right")
	ErrBlockChainIDWrong = errors.New("block chainid not right")
	ErrBlockSignWrong = errors.New("block sign verify failed")
	ErrBlockPublicKeyWrong = errors.New("block public key not right")
)

type Block struct {
	hash     []byte
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
	if b.hash != nil{
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
	hashes := make([][]byte, 0, len(b.PBBlock.Messages))
	if len(b.PBBlock.MessageHashes) != len(b.PBBlock.Messages) {
		return nil, ErrBlockMessagHashesNotFit
	}

	for i, msg := range b.PBBlock.Messages {
		m := Message{
			PBMessage: msg,
		}
		hash, err := m.Hash()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(b.PBBlock.MessageHashes[i], hash){
			return nil, ErrBlockMessagHashesWrong
		}
		hashes = append(hashes, hash)
	}
	merkletree.Build(&tree, hashes)
	return merkletree.RootHash(&tree), nil
}

func (b *Block) Sign(privateKey []byte) error{
	alg := crypto.Ed25519
	publicKey, err := alg.GeneratePublicKey(privateKey)
	if err !=nil{
		return err
	}
	if !bytes.Equal(b.PBBlock.Head.PublicKey, publicKey){
		return ErrBlockPublicKeyWrong

	}
	err = b.verifyData()
	if err !=nil {
		return err
	}
	hash, err := b.Hash()
	if err != nil {
		return err
	}
	sig, err := alg.Sign(hash, privateKey)
	if err != nil{
		return err
	}
	b.PBBlock.Signature.Algorithm = int32(alg)
	b.PBBlock.Signature.Sig = sig
	b.PBBlock.Head.PublicKey = publicKey
	return nil
}

func (b *Block) Verify() error {
	err := b.verifyData()
	if err != nil{
		return err
	}
	return b.VerifySign()
}

func (b *Block) verifyData() error {
	//check message
	for _, msg := range b.PBBlock.Messages {
		m := ToMessage(msg)
		err := m.Verify()
		if err !=nil {}
		return nil
	}
	//check message hash
	if len(b.PBBlock.MessageHashes) != len(b.PBBlock.Messages) {
		return ErrBlockMessagHashesNotFit
	}
	//check merkle
	merkleTree, err:= b.CalculateTxMerkleHash()
	if err != nil {
		return err
	}
	if !bytes.Equal(b.PBBlock.Head.MessageMerkleHash, merkleTree){
		return ErrBlockMerlkHashWrong
	}
	if b.PBBlock.Head.ChainId != consensus.Consensus.ChainID{
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
	alg :=crypto.Algorithm(b.PBBlock.Signature.Algorithm)
	if !alg.Verify(b.hash, b.PBBlock.Head.PublicKey, b.PBBlock.Signature.Sig){
		return ErrBlockSignWrong
	}
	return nil
}

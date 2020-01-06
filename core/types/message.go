package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/zhangdaoling/notokenblockchain/common"
	"github.com/zhangdaoling/notokenblockchain/core/consensus"
	"github.com/zhangdaoling/notokenblockchain/crypto"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/notokenblockchain/pb"
)

var (
	ErrMessagePublicKeyWrong = errors.New("block public key not right")
)

type Message struct {
	hash      []byte
	PBMessage *pb.Message
}

func ToMessage(m *pb.Message) *Message {
	return &Message{
		PBMessage: m,
	}
}

func ToPbMessage(m *Message) *pb.Message {
	return m.PBMessage
}

func (m *Message) String() string {
	return ""
}

func (m *Message) P2PBytes() ([]byte, error) {
	return proto.Marshal(m.PBMessage)
}

func (m *Message) P2PDecode(buf []byte) error {
	pbMessage := &pb.Message{}
	err := proto.Unmarshal(buf, pbMessage)
	if err != nil {
		return err
	}
	m.PBMessage = pbMessage
	return nil
}

func (m *Message) DBBytes() ([]byte, error) {
	return proto.Marshal(m.PBMessage)
}

func (m *Message) DBDecode(buf []byte) error {
	pbMessage := &pb.Message{}
	err := proto.Unmarshal(buf, pbMessage)
	if err != nil {
		return err
	}
	m.PBMessage = pbMessage
	return nil
}

func (m *Message) HashBytes() ([]byte, error) {
	return proto.Marshal(m.PBMessage.UnsignMessage)
}

func (m *Message) Hash() ([]byte, error) {
	baseBytes, err := m.HashBytes()
	if err != nil {
		return nil, err
	}
	return common.Sha3(baseBytes), nil
}

func (m *Message) Sign(privateKey []byte) error {
	alg := crypto.Ed25519
	publicKey, err := alg.GeneratePublicKey(privateKey)
	if err != nil {
		return err
	}
	if !bytes.Equal(m.PBMessage.UnsignMessage.PublicKey, publicKey) {
		return ErrMessagePublicKeyWrong

	}

	hash, err := m.Hash()
	if err != nil {
		return err
	}
	sig, err := alg.Sign(hash, privateKey)
	if err != nil {
		return err
	}
	m.PBMessage.Signature.Algorithm = int32(alg)
	m.PBMessage.Signature.Sig = sig
	m.PBMessage.UnsignMessage.PublicKey = publicKey
	return nil
}

func (m *Message) Verify() error {
	baseBytes, err := m.HashBytes()
	if err != nil {
		return err
	}
	m.hash = common.Sha3(baseBytes)
	alg := crypto.Algorithm(m.PBMessage.Signature.Algorithm)
	if !alg.Verify(m.hash, m.PBMessage.UnsignMessage.PublicKey, m.PBMessage.Signature.Sig) {
		return ErrBlockSignWrong
	}
	return nil
}

func (m *Message) IsExpired(ct int64) bool {
	if m.PBMessage.UnsignMessage.Time <= ct {
		return true
	}
	if ct-m.PBMessage.UnsignMessage.Time > consensus.Consensus.MessageExpirationTime {
		return true
	}
	return false
}

func (m *Message) CheckSize() error {
	b, err := m.P2PBytes()
	if err != nil {
		return err
	}

	l := int64(len(b))
	if l > consensus.Consensus.MessageSizeLimit {
		return fmt.Errorf("msg size illegal, should <= %v, got %v", consensus.Consensus.MessageSizeLimit, l)
	}
	return nil
}

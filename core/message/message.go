package message

import (
	"fmt"

	"github.com/zhangdaoling/simplechain/common"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/simplechain/pb"
)

type Message struct {
	PBMessage *pb.Message
}

func ToMessage(m *pb.Message) *Message{
	return &Message{
		PBMessage:m,
	}
}

func ToPb(m *Message) *pb.Message{
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
	if err != nil{
		return nil, err
	}
	buf, err := common.Sha3(baseBytes)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//to do
func (m *Message) Sign() {
	return
}

//to do
func (m *Message) Verify() error {
	return nil
}

func (m *Message) IsExpired(ct int64) bool {
	if m.PBMessage.UnsignMessage.Time <= ct {
		return true
	}
	if ct-m.PBMessage.UnsignMessage.Time > common.MessageMaxExpiration {
		return true
	}
	return false
}

func (m *Message) CheckSize() error {
	b, err := m.P2PBytes()
	if err != nil {
		return err
	}

	l := len(b)
	if l > common.MessageSizeLimit {
		return fmt.Errorf("msg size illegal, should <= %v, got %v", common.MessageSizeLimit, l)
	}
	return nil
}

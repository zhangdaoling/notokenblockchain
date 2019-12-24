package message

import (
	"fmt"

	"github.com/zhangdaoling/simplechain/common"

	"github.com/golang/protobuf/proto"
	"github.com/zhangdaoling/simplechain/pb"
)

type Message struct {
	BaseHash  []byte
	FullHash  []byte
	PbMessage *pb.Message
}

func NewMessage() *Message {
	return nil
}

func (m *Message) String() string {
	return ""
}

func (m *Message) Encode() ([]byte, error) {
	return m.FullByte()
}

func (m *Message) Decode(b []byte) error {
	msg := &pb.Message{}
	err := proto.Unmarshal(b, msg)
	if err != nil {
		return err
	}
	m.PbMessage = msg
	_, err = m.CalculateFullHash()
	if err != nil {
		return err
	}
	_, err = m.CalculateBaseHash()
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) Hash() ([]byte, error){
	return m.CalculateBaseHash()
}

func (m *Message) BaseByte() ([]byte, error) {
	buf, err := proto.Marshal(m.PbMessage.Unsignmessage)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (m *Message) FullByte() ([]byte, error) {
	buf, err := proto.Marshal(m.PbMessage)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (m *Message) CalculateBaseHash() (buf []byte, err error) {
	if m.BaseHash == nil {
		tmp, err := m.BaseByte()
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

func (m *Message) CalculateFullHash() (buf []byte, err error) {
	if m.FullHash == nil {
		tmp, err := m.FullByte()
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

func (m *Message) Sign() {
	return
}

func (m *Message) Verify() error {
	return nil
}

func (m *Message) IsExpired(ct int64) bool {
	if m.PbMessage.Unsignmessage.Time <= ct {
		return true
	}
	if ct-m.PbMessage.Unsignmessage.Time > common.MessageMaxExpiration {
		return true
	}
	return false
}

func (m *Message) CheckSize() error {
	b, err := m.FullByte()
	if err != nil {
		return err
	}

	l := len(b)
	if l > common.MessageSizeLimit {
		return fmt.Errorf("tx size illegal, should <= %v, got %v", common.MessageSizeLimit, l)
	}
	return nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Signature            *Signature     `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	UnsignMessage        *UnSignMessage `protobuf:"bytes,2,opt,name=unsign_message,json=unsignMessage,proto3" json:"unsign_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_edeb343296ceb6bc, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetSignature() *Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Message) GetUnsignMessage() *UnSignMessage {
	if m != nil {
		return m.UnsignMessage
	}
	return nil
}

type UnSignMessage struct {
	ChainId              int64    `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Nonce                int64    `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Time                 int64    `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	ReferenceBlock       string   `protobuf:"bytes,4,opt,name=reference_block,json=referenceBlock,proto3" json:"reference_block,omitempty"`
	Sender               string   `protobuf:"bytes,5,opt,name=sender,proto3" json:"sender,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnSignMessage) Reset()         { *m = UnSignMessage{} }
func (m *UnSignMessage) String() string { return proto.CompactTextString(m) }
func (*UnSignMessage) ProtoMessage()    {}
func (*UnSignMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_edeb343296ceb6bc, []int{1}
}
func (m *UnSignMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnSignMessage.Unmarshal(m, b)
}
func (m *UnSignMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnSignMessage.Marshal(b, m, deterministic)
}
func (dst *UnSignMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnSignMessage.Merge(dst, src)
}
func (m *UnSignMessage) XXX_Size() int {
	return xxx_messageInfo_UnSignMessage.Size(m)
}
func (m *UnSignMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_UnSignMessage.DiscardUnknown(m)
}

var xxx_messageInfo_UnSignMessage proto.InternalMessageInfo

func (m *UnSignMessage) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *UnSignMessage) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *UnSignMessage) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *UnSignMessage) GetReferenceBlock() string {
	if m != nil {
		return m.ReferenceBlock
	}
	return ""
}

func (m *UnSignMessage) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "pb.Message")
	proto.RegisterType((*UnSignMessage)(nil), "pb.UnSignMessage")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_message_edeb343296ceb6bc) }

var fileDescriptor_message_edeb343296ceb6bc = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x4d, 0x4b, 0x03, 0x31,
	0x10, 0x86, 0xd9, 0x6e, 0x3f, 0xec, 0xc8, 0xb6, 0x38, 0x88, 0x44, 0x4f, 0xa5, 0x17, 0x0b, 0xc2,
	0x1e, 0xf4, 0xe2, 0xd9, 0x9b, 0x07, 0x2f, 0x11, 0xcf, 0xcb, 0x26, 0x3b, 0xd6, 0xa0, 0x9d, 0x84,
	0x64, 0xfb, 0x53, 0xfc, 0xbf, 0x92, 0x0f, 0x2b, 0xde, 0xf2, 0x3c, 0xf3, 0xbe, 0x33, 0x10, 0x68,
	0x0e, 0x14, 0x42, 0xbf, 0xa7, 0xd6, 0x79, 0x3b, 0x5a, 0x9c, 0x38, 0x75, 0xb3, 0x0e, 0x66, 0xcf,
	0xfd, 0x78, 0xf4, 0x45, 0x6e, 0x1d, 0x2c, 0x5e, 0x72, 0x0a, 0xef, 0x60, 0x79, 0x9a, 0x8a, 0x6a,
	0x53, 0xed, 0xce, 0xef, 0x9b, 0xd6, 0xa9, 0xf6, 0xf5, 0x57, 0xca, 0xbf, 0x39, 0x3e, 0xc2, 0xea,
	0xc8, 0x11, 0xbb, 0x72, 0x44, 0x4c, 0x52, 0xe3, 0x22, 0x36, 0xde, 0x38, 0x76, 0xca, 0x5e, 0xd9,
	0xe4, 0x60, 0xc1, 0xed, 0x77, 0x05, 0xcd, 0xbf, 0x00, 0x5e, 0xc3, 0x99, 0xfe, 0xe8, 0x0d, 0x77,
	0x66, 0x48, 0x77, 0x6b, 0xb9, 0x48, 0xfc, 0x3c, 0xe0, 0x25, 0xcc, 0xd8, 0xb2, 0xce, 0xdb, 0x6b,
	0x99, 0x01, 0x11, 0xa6, 0xa3, 0x39, 0x90, 0xa8, 0x93, 0x4c, 0x6f, 0xbc, 0x85, 0xb5, 0xa7, 0x77,
	0xf2, 0xc4, 0x9a, 0x3a, 0xf5, 0x65, 0xf5, 0xa7, 0x98, 0x6e, 0xaa, 0xdd, 0x52, 0xae, 0x4e, 0xfa,
	0x29, 0x5a, 0xbc, 0x82, 0x79, 0x20, 0x1e, 0xc8, 0x8b, 0x59, 0x9a, 0x17, 0x52, 0xf3, 0xf4, 0x21,
	0x0f, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1e, 0xf5, 0xbf, 0x12, 0x36, 0x01, 0x00, 0x00,
}
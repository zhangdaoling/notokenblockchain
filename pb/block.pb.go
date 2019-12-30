// Code generated by protoc-gen-go. DO NOT EDIT.
// source: block.proto

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

type Block struct {
	Head                 *BlockHead `protobuf:"bytes,1,opt,name=head,proto3" json:"head,omitempty"`
	Sign                 *Signature `protobuf:"bytes,2,opt,name=sign,proto3" json:"sign,omitempty"`
	MessageHashes        [][]byte   `protobuf:"bytes,3,rep,name=message_hashes,json=messageHashes,proto3" json:"message_hashes,omitempty"`
	Messages             []*Message `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_block_7491cb4fc9a03e4e, []int{0}
}
func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (dst *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(dst, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetHead() *BlockHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *Block) GetSign() *Signature {
	if m != nil {
		return m.Sign
	}
	return nil
}

func (m *Block) GetMessageHashes() [][]byte {
	if m != nil {
		return m.MessageHashes
	}
	return nil
}

func (m *Block) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

type BlockHead struct {
	Version              int64    `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	ChainId              int64    `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Height               int64    `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Weight               int64    `protobuf:"varint,4,opt,name=weight,proto3" json:"weight,omitempty"`
	Time                 int64    `protobuf:"varint,5,opt,name=time,proto3" json:"time,omitempty"`
	ParentHash           []byte   `protobuf:"bytes,6,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	MessageMerkleHash    []byte   `protobuf:"bytes,7,opt,name=message_merkle_hash,json=messageMerkleHash,proto3" json:"message_merkle_hash,omitempty"`
	Blocker              string   `protobuf:"bytes,8,opt,name=blocker,proto3" json:"blocker,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHead) Reset()         { *m = BlockHead{} }
func (m *BlockHead) String() string { return proto.CompactTextString(m) }
func (*BlockHead) ProtoMessage()    {}
func (*BlockHead) Descriptor() ([]byte, []int) {
	return fileDescriptor_block_7491cb4fc9a03e4e, []int{1}
}
func (m *BlockHead) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockHead.Unmarshal(m, b)
}
func (m *BlockHead) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockHead.Marshal(b, m, deterministic)
}
func (dst *BlockHead) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHead.Merge(dst, src)
}
func (m *BlockHead) XXX_Size() int {
	return xxx_messageInfo_BlockHead.Size(m)
}
func (m *BlockHead) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHead.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHead proto.InternalMessageInfo

func (m *BlockHead) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *BlockHead) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *BlockHead) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHead) GetWeight() int64 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *BlockHead) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *BlockHead) GetParentHash() []byte {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *BlockHead) GetMessageMerkleHash() []byte {
	if m != nil {
		return m.MessageMerkleHash
	}
	return nil
}

func (m *BlockHead) GetBlocker() string {
	if m != nil {
		return m.Blocker
	}
	return ""
}

func init() {
	proto.RegisterType((*Block)(nil), "pb.Block")
	proto.RegisterType((*BlockHead)(nil), "pb.BlockHead")
}

func init() { proto.RegisterFile("block.proto", fileDescriptor_block_7491cb4fc9a03e4e) }

var fileDescriptor_block_7491cb4fc9a03e4e = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x91, 0x4f, 0x4e, 0xf3, 0x30,
	0x10, 0xc5, 0xe5, 0x3a, 0x6d, 0xd2, 0x49, 0xfb, 0x7d, 0xc2, 0x48, 0xc8, 0x74, 0x43, 0xa8, 0x84,
	0xc8, 0x2a, 0x8b, 0x72, 0x03, 0x56, 0x65, 0xd1, 0x8d, 0x39, 0x40, 0xe5, 0x34, 0x56, 0x6c, 0xb5,
	0xf9, 0x23, 0x3b, 0xc0, 0x69, 0x38, 0x23, 0x57, 0x40, 0x1e, 0x27, 0x65, 0x97, 0xf7, 0x7b, 0x2f,
	0xe3, 0x79, 0x1a, 0x48, 0xcb, 0x4b, 0x77, 0x3a, 0x17, 0xbd, 0xed, 0x86, 0x8e, 0xcd, 0xfa, 0x72,
	0xf3, 0xdf, 0x99, 0xba, 0x95, 0xc3, 0x87, 0x55, 0x01, 0x6e, 0xd6, 0x8d, 0x72, 0x4e, 0xd6, 0xa3,
	0xdc, 0x7e, 0x13, 0x98, 0xbf, 0xfa, 0x7f, 0xd8, 0x23, 0x44, 0x5a, 0xc9, 0x8a, 0x93, 0x8c, 0xe4,
	0xe9, 0x6e, 0x5d, 0xf4, 0x65, 0x81, 0xc6, 0x5e, 0xc9, 0x4a, 0xa0, 0xe5, 0x23, 0x7e, 0x1c, 0x9f,
	0xfd, 0x45, 0xde, 0xa7, 0xf1, 0x02, 0x2d, 0xf6, 0x04, 0xff, 0xc6, 0x07, 0x8e, 0x5a, 0x3a, 0xad,
	0x1c, 0xa7, 0x19, 0xcd, 0x57, 0x62, 0x7a, 0x76, 0x8f, 0x90, 0x3d, 0x43, 0x32, 0x02, 0xc7, 0xa3,
	0x8c, 0xe6, 0xe9, 0x2e, 0xf5, 0xd3, 0x0e, 0x81, 0x89, 0xab, 0xb9, 0xfd, 0x21, 0xb0, 0xbc, 0xae,
	0xc1, 0x38, 0xc4, 0x9f, 0xca, 0x3a, 0xd3, 0xb5, 0xb8, 0x26, 0x15, 0x93, 0x64, 0xf7, 0x90, 0x9c,
	0xb4, 0x34, 0xed, 0xd1, 0x54, 0xb8, 0x1e, 0x15, 0x31, 0xea, 0xb7, 0x8a, 0xdd, 0xc1, 0x42, 0x2b,
	0x53, 0xeb, 0x81, 0x53, 0x34, 0x46, 0xe5, 0xf9, 0x57, 0xe0, 0x51, 0xe0, 0x41, 0x31, 0x06, 0xd1,
	0x60, 0x1a, 0xc5, 0xe7, 0x48, 0xf1, 0x9b, 0x3d, 0x40, 0xda, 0x4b, 0xab, 0xda, 0x01, 0x5b, 0xf1,
	0x45, 0x46, 0xf2, 0x95, 0x80, 0x80, 0x7c, 0x25, 0x56, 0xc0, 0xed, 0xd4, 0xbb, 0x51, 0xf6, 0x7c,
	0x09, 0xf5, 0x79, 0x8c, 0xc1, 0x9b, 0xd1, 0x3a, 0xa0, 0x83, 0x79, 0x0e, 0x31, 0x9e, 0x4a, 0x59,
	0x9e, 0x64, 0x24, 0x5f, 0x8a, 0x49, 0x96, 0x0b, 0x3c, 0xcc, 0xcb, 0x6f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x23, 0x18, 0x57, 0xb1, 0xcb, 0x01, 0x00, 0x00,
}
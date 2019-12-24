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
	Messages             []*Message `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
	MessageHashes        [][]byte   `protobuf:"bytes,4,rep,name=messageHashes,proto3" json:"messageHashes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_block_8f75e8e8cd1b3446, []int{0}
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

func (m *Block) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *Block) GetMessageHashes() [][]byte {
	if m != nil {
		return m.MessageHashes
	}
	return nil
}

type BlockHead struct {
	Version              int64    `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Number               int64    `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
	ChainId              int64    `protobuf:"varint,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Time                 int64    `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
	ParentHash           []byte   `protobuf:"bytes,5,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	MessageMerkleHash    []byte   `protobuf:"bytes,6,opt,name=message_merkle_hash,json=messageMerkleHash,proto3" json:"message_merkle_hash,omitempty"`
	Blocker              string   `protobuf:"bytes,7,opt,name=blocker,proto3" json:"blocker,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHead) Reset()         { *m = BlockHead{} }
func (m *BlockHead) String() string { return proto.CompactTextString(m) }
func (*BlockHead) ProtoMessage()    {}
func (*BlockHead) Descriptor() ([]byte, []int) {
	return fileDescriptor_block_8f75e8e8cd1b3446, []int{1}
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

func (m *BlockHead) GetNumber() int64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *BlockHead) GetChainId() int64 {
	if m != nil {
		return m.ChainId
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

func init() { proto.RegisterFile("block.proto", fileDescriptor_block_8f75e8e8cd1b3446) }

var fileDescriptor_block_8f75e8e8cd1b3446 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x4b, 0x6e, 0xb3, 0x30,
	0x10, 0xc7, 0xe5, 0x98, 0x40, 0x32, 0x24, 0xfa, 0xf4, 0x4d, 0xa5, 0xca, 0xcd, 0xa6, 0x34, 0xaa,
	0x54, 0x56, 0x2c, 0xd2, 0x1b, 0x74, 0x95, 0x2e, 0xb2, 0x71, 0x0f, 0x80, 0x4c, 0x18, 0x05, 0x94,
	0xf0, 0x90, 0x4d, 0x7a, 0x9a, 0xde, 0xab, 0xd7, 0xa9, 0x18, 0x20, 0x55, 0x77, 0xfe, 0x3f, 0x66,
	0xfc, 0x93, 0x0d, 0x61, 0x76, 0x69, 0x8e, 0xe7, 0xa4, 0xb5, 0x4d, 0xd7, 0xe0, 0xac, 0xcd, 0x36,
	0xff, 0x5c, 0x79, 0xaa, 0x4d, 0x77, 0xb5, 0x34, 0x98, 0x9b, 0x75, 0x45, 0xce, 0x99, 0xd3, 0x28,
	0xb7, 0x5f, 0x02, 0xe6, 0x6f, 0xfd, 0x0c, 0x3e, 0x81, 0x57, 0x90, 0xc9, 0x95, 0x88, 0x44, 0x1c,
	0xee, 0xd6, 0x49, 0x9b, 0x25, 0x1c, 0xec, 0xc9, 0xe4, 0x9a, 0xa3, 0xbe, 0xd2, 0xaf, 0x53, 0xb3,
	0xdf, 0xca, 0xc7, 0xb4, 0x5e, 0x73, 0x84, 0x2f, 0xb0, 0x18, 0x2f, 0x70, 0x4a, 0x46, 0x32, 0x0e,
	0x77, 0x61, 0x5f, 0x3b, 0x0c, 0x9e, 0xbe, 0x85, 0xf8, 0x0c, 0x13, 0xc9, 0xde, 0xb8, 0x82, 0x9c,
	0xf2, 0x22, 0x19, 0xaf, 0xf4, 0x5f, 0x73, 0xfb, 0x2d, 0x60, 0x79, 0xa3, 0x40, 0x05, 0xc1, 0x27,
	0x59, 0x57, 0x36, 0x35, 0x53, 0x4a, 0x3d, 0x49, 0xbc, 0x07, 0xbf, 0xbe, 0x56, 0x19, 0x59, 0x66,
	0x93, 0x7a, 0x54, 0xf8, 0x00, 0x8b, 0x63, 0x61, 0xca, 0x3a, 0x2d, 0x73, 0x25, 0x87, 0x11, 0xd6,
	0xef, 0x39, 0x22, 0x78, 0x5d, 0x59, 0x91, 0xf2, 0xd8, 0xe6, 0x33, 0x3e, 0x42, 0xd8, 0x1a, 0x4b,
	0x75, 0x97, 0x16, 0xc6, 0x15, 0x6a, 0x1e, 0x89, 0x78, 0xa5, 0x61, 0xb0, 0x7a, 0x22, 0x4c, 0xe0,
	0x6e, 0x04, 0x4c, 0x2b, 0xb2, 0xe7, 0x0b, 0x0d, 0x45, 0x9f, 0x8b, 0xff, 0xc7, 0xe8, 0xc0, 0x09,
	0xf7, 0x15, 0x04, 0xfc, 0x23, 0x64, 0x55, 0x10, 0x89, 0x78, 0xa9, 0x27, 0x99, 0xf9, 0xfc, 0xfe,
	0xaf, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x17, 0x2f, 0xfe, 0x39, 0xb2, 0x01, 0x00, 0x00,
}

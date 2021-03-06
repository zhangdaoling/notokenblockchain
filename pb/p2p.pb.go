// Code generated by protoc-gen-go. DO NOT EDIT.
// source: p2p.proto

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

type RoutingQuery struct {
	Ids                  []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoutingQuery) Reset()         { *m = RoutingQuery{} }
func (m *RoutingQuery) String() string { return proto.CompactTextString(m) }
func (*RoutingQuery) ProtoMessage()    {}
func (*RoutingQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_969a50c227bce861, []int{0}
}
func (m *RoutingQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutingQuery.Unmarshal(m, b)
}
func (m *RoutingQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutingQuery.Marshal(b, m, deterministic)
}
func (dst *RoutingQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutingQuery.Merge(dst, src)
}
func (m *RoutingQuery) XXX_Size() int {
	return xxx_messageInfo_RoutingQuery.Size(m)
}
func (m *RoutingQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_RoutingQuery.DiscardUnknown(m)
}

var xxx_messageInfo_RoutingQuery proto.InternalMessageInfo

func (m *RoutingQuery) GetIds() []string {
	if m != nil {
		return m.Ids
	}
	return nil
}

type PeerInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addrs                []string `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerInfo) Reset()         { *m = PeerInfo{} }
func (m *PeerInfo) String() string { return proto.CompactTextString(m) }
func (*PeerInfo) ProtoMessage()    {}
func (*PeerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_969a50c227bce861, []int{1}
}
func (m *PeerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerInfo.Unmarshal(m, b)
}
func (m *PeerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerInfo.Marshal(b, m, deterministic)
}
func (dst *PeerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerInfo.Merge(dst, src)
}
func (m *PeerInfo) XXX_Size() int {
	return xxx_messageInfo_PeerInfo.Size(m)
}
func (m *PeerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PeerInfo proto.InternalMessageInfo

func (m *PeerInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PeerInfo) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

type RoutingResponse struct {
	Peers                []*PeerInfo `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RoutingResponse) Reset()         { *m = RoutingResponse{} }
func (m *RoutingResponse) String() string { return proto.CompactTextString(m) }
func (*RoutingResponse) ProtoMessage()    {}
func (*RoutingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_969a50c227bce861, []int{2}
}
func (m *RoutingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutingResponse.Unmarshal(m, b)
}
func (m *RoutingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutingResponse.Marshal(b, m, deterministic)
}
func (dst *RoutingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutingResponse.Merge(dst, src)
}
func (m *RoutingResponse) XXX_Size() int {
	return xxx_messageInfo_RoutingResponse.Size(m)
}
func (m *RoutingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RoutingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RoutingResponse proto.InternalMessageInfo

func (m *RoutingResponse) GetPeers() []*PeerInfo {
	if m != nil {
		return m.Peers
	}
	return nil
}

func init() {
	proto.RegisterType((*RoutingQuery)(nil), "pb.RoutingQuery")
	proto.RegisterType((*PeerInfo)(nil), "pb.PeerInfo")
	proto.RegisterType((*RoutingResponse)(nil), "pb.RoutingResponse")
}

func init() { proto.RegisterFile("p2p.proto", fileDescriptor_p2p_969a50c227bce861) }

var fileDescriptor_p2p_969a50c227bce861 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x30, 0x2a, 0xd0,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x52, 0xe0, 0xe2, 0x09, 0xca, 0x2f,
	0x2d, 0xc9, 0xcc, 0x4b, 0x0f, 0x2c, 0x4d, 0x2d, 0xaa, 0x14, 0x12, 0xe0, 0x62, 0xce, 0x4c, 0x29,
	0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x0c, 0x02, 0x31, 0x95, 0x0c, 0xb8, 0x38, 0x02, 0x52, 0x53,
	0x8b, 0x3c, 0xf3, 0xd2, 0xf2, 0x85, 0xf8, 0xb8, 0x98, 0x32, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35,
	0x38, 0x83, 0x98, 0x32, 0x53, 0x84, 0x44, 0xb8, 0x58, 0x13, 0x53, 0x52, 0x8a, 0x8a, 0x25, 0x98,
	0xc0, 0xea, 0x21, 0x1c, 0x25, 0x53, 0x2e, 0x7e, 0xa8, 0x99, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79,
	0xc5, 0xa9, 0x42, 0x4a, 0x5c, 0xac, 0x05, 0xa9, 0xa9, 0x45, 0x10, 0x83, 0xb9, 0x8d, 0x78, 0xf4,
	0x0a, 0x92, 0xf4, 0x60, 0xa6, 0x06, 0x41, 0xa4, 0x92, 0xd8, 0xc0, 0xae, 0x32, 0x06, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x61, 0xf8, 0x32, 0xfe, 0xa2, 0x00, 0x00, 0x00,
}

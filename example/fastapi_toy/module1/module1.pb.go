// Code generated by protoc-gen-go.
// source: module1.proto
// DO NOT EDIT!

/*
Package module1 is a generated protocol buffer package.

It is generated from these files:
	module1.proto

It has these top-level messages:
	AddReq
	AddRsp
*/
package module1

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

type AddReq struct {
	A int32 `protobuf:"varint,1,opt,name=A,json=a" json:"A,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=B,json=b" json:"B,omitempty"`
}

func (m *AddReq) Reset()                    { *m = AddReq{} }
func (m *AddReq) String() string            { return proto.CompactTextString(m) }
func (*AddReq) ProtoMessage()               {}
func (*AddReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AddReq) GetA() int32 {
	if m != nil {
		return m.A
	}
	return 0
}

func (m *AddReq) GetB() int32 {
	if m != nil {
		return m.B
	}
	return 0
}

type AddRsp struct {
	C int32 `protobuf:"varint,1,opt,name=C,json=c" json:"C,omitempty"`
}

func (m *AddRsp) Reset()                    { *m = AddRsp{} }
func (m *AddRsp) String() string            { return proto.CompactTextString(m) }
func (*AddRsp) ProtoMessage()               {}
func (*AddRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AddRsp) GetC() int32 {
	if m != nil {
		return m.C
	}
	return 0
}

func init() {
	proto.RegisterType((*AddReq)(nil), "module1.AddReq")
	proto.RegisterType((*AddRsp)(nil), "module1.AddRsp")
}

func init() { proto.RegisterFile("module1.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 91 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xcd, 0x4f, 0x29,
	0xcd, 0x49, 0x35, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x54, 0xb8,
	0xd8, 0x1c, 0x53, 0x52, 0x82, 0x52, 0x0b, 0x85, 0x78, 0xb8, 0x18, 0x1d, 0x25, 0x18, 0x15, 0x18,
	0x35, 0x58, 0x83, 0x18, 0x13, 0x41, 0x3c, 0x27, 0x09, 0x26, 0x08, 0x2f, 0x49, 0x49, 0x0c, 0xa2,
	0xaa, 0xb8, 0x00, 0x24, 0xee, 0x0c, 0x53, 0x95, 0x9c, 0xc4, 0x06, 0x36, 0xcd, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0x0b, 0x12, 0xe2, 0x9b, 0x5e, 0x00, 0x00, 0x00,
}

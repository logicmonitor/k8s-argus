// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/common/value.proto

package common

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// A generic data container.
type Value struct {
	// A value.
	//
	// Types that are valid to be assigned to Value:
	//	*Value_BooleanValue
	//	*Value_Int64Value
	//	*Value_FloatValue
	//	*Value_DoubleValue
	//	*Value_StringValue
	Value                isValue_Value `protobuf_oneof:"value"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_a73d72db515a577a, []int{0}
}

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

type isValue_Value interface {
	isValue_Value()
}

type Value_BooleanValue struct {
	BooleanValue bool `protobuf:"varint,1,opt,name=boolean_value,json=booleanValue,proto3,oneof"`
}

type Value_Int64Value struct {
	Int64Value int64 `protobuf:"varint,2,opt,name=int64_value,json=int64Value,proto3,oneof"`
}

type Value_FloatValue struct {
	FloatValue float32 `protobuf:"fixed32,3,opt,name=float_value,json=floatValue,proto3,oneof"`
}

type Value_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,4,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

type Value_StringValue struct {
	StringValue string `protobuf:"bytes,5,opt,name=string_value,json=stringValue,proto3,oneof"`
}

func (*Value_BooleanValue) isValue_Value() {}

func (*Value_Int64Value) isValue_Value() {}

func (*Value_FloatValue) isValue_Value() {}

func (*Value_DoubleValue) isValue_Value() {}

func (*Value_StringValue) isValue_Value() {}

func (m *Value) GetValue() isValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Value) GetBooleanValue() bool {
	if x, ok := m.GetValue().(*Value_BooleanValue); ok {
		return x.BooleanValue
	}
	return false
}

func (m *Value) GetInt64Value() int64 {
	if x, ok := m.GetValue().(*Value_Int64Value); ok {
		return x.Int64Value
	}
	return 0
}

func (m *Value) GetFloatValue() float32 {
	if x, ok := m.GetValue().(*Value_FloatValue); ok {
		return x.FloatValue
	}
	return 0
}

func (m *Value) GetDoubleValue() float64 {
	if x, ok := m.GetValue().(*Value_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

func (m *Value) GetStringValue() string {
	if x, ok := m.GetValue().(*Value_StringValue); ok {
		return x.StringValue
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Value) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Value_BooleanValue)(nil),
		(*Value_Int64Value)(nil),
		(*Value_FloatValue)(nil),
		(*Value_DoubleValue)(nil),
		(*Value_StringValue)(nil),
	}
}

func init() {
	proto.RegisterType((*Value)(nil), "google.ads.googleads.v3.common.Value")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/common/value.proto", fileDescriptor_a73d72db515a577a)
}

var fileDescriptor_a73d72db515a577a = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0xd1, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0x06, 0xf0, 0xa5, 0x73, 0xfe, 0xc9, 0xe6, 0x65, 0x27, 0x11, 0x19, 0x75, 0x22, 0x14, 0x0f,
	0xe9, 0xa1, 0xe2, 0x21, 0x9e, 0x3a, 0x85, 0xed, 0x38, 0x76, 0xe8, 0x41, 0x0a, 0x92, 0xad, 0x35,
	0x14, 0xb2, 0xbc, 0x63, 0xc9, 0xf6, 0x81, 0x3c, 0xfa, 0x45, 0x04, 0xbf, 0x87, 0x17, 0x3f, 0x85,
	0x24, 0x6f, 0xd6, 0x9b, 0x9e, 0xfa, 0xf2, 0xf4, 0x97, 0x87, 0xf6, 0x0d, 0xbd, 0x93, 0x00, 0x52,
	0xd5, 0xa9, 0xa8, 0x4c, 0x8a, 0xa3, 0x9b, 0xf6, 0x59, 0xba, 0x82, 0xf5, 0x1a, 0x74, 0xba, 0x17,
	0x6a, 0x57, 0xb3, 0xcd, 0x16, 0x2c, 0x0c, 0x47, 0x08, 0x98, 0xa8, 0x0c, 0x6b, 0x2d, 0xdb, 0x67,
	0x0c, 0xed, 0xe5, 0xd5, 0xa1, 0x6b, 0xd3, 0xa4, 0x42, 0x6b, 0xb0, 0xc2, 0x36, 0xa0, 0x0d, 0x9e,
	0x1e, 0x7f, 0x12, 0xda, 0x2b, 0x5c, 0xdb, 0xf0, 0x96, 0x9e, 0x2f, 0x01, 0x54, 0x2d, 0xf4, 0xab,
	0xaf, 0xbf, 0x20, 0x31, 0x49, 0x4e, 0x67, 0x9d, 0xc5, 0x20, 0xc4, 0xc8, 0xae, 0x69, 0xbf, 0xd1,
	0xf6, 0xe1, 0x3e, 0xa0, 0x28, 0x26, 0x49, 0x77, 0xd6, 0x59, 0x50, 0x1f, 0xb6, 0xe4, 0x4d, 0x81,
	0xb0, 0x81, 0x74, 0x63, 0x92, 0x44, 0x8e, 0xf8, 0x10, 0xc9, 0x0d, 0x1d, 0x54, 0xb0, 0x5b, 0xaa,
	0x3a, 0x98, 0xa3, 0x98, 0x24, 0x64, 0xd6, 0x59, 0xf4, 0x31, 0x6d, 0x91, 0xb1, 0xdb, 0x46, 0xcb,
	0x80, 0x7a, 0x31, 0x49, 0xce, 0x1c, 0xc2, 0xd4, 0xa3, 0xc9, 0x09, 0xed, 0xf9, 0xb7, 0x93, 0x6f,
	0x42, 0xc7, 0x2b, 0x58, 0xb3, 0xff, 0xd7, 0x31, 0xa1, 0xfe, 0xd8, 0xdc, 0xfd, 0xfc, 0x9c, 0xbc,
	0x3c, 0x07, 0x2d, 0x41, 0x09, 0x2d, 0x19, 0x6c, 0x65, 0x2a, 0x6b, 0xed, 0x57, 0x73, 0x58, 0xfc,
	0xa6, 0x31, 0x7f, 0xdd, 0xc3, 0x23, 0x3e, 0xde, 0xa3, 0xee, 0x34, 0xcf, 0x3f, 0xa2, 0xd1, 0x14,
	0xcb, 0xf2, 0xca, 0x30, 0x1c, 0xdd, 0x54, 0x64, 0xec, 0xc9, 0xb3, 0xaf, 0x03, 0x28, 0xf3, 0xca,
	0x94, 0x2d, 0x28, 0x8b, 0xac, 0x44, 0xf0, 0x13, 0x8d, 0x31, 0xe5, 0x3c, 0xaf, 0x0c, 0xe7, 0x2d,
	0xe1, 0xbc, 0xc8, 0x38, 0x47, 0xb4, 0x3c, 0xf6, 0x5f, 0x97, 0xfd, 0x06, 0x00, 0x00, 0xff, 0xff,
	0x12, 0x9e, 0x55, 0x7b, 0x24, 0x02, 0x00, 0x00,
}
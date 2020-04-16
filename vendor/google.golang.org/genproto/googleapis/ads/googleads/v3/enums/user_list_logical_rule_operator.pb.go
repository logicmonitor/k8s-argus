// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/enums/user_list_logical_rule_operator.proto

package enums

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

// Enum describing possible user list logical rule operators.
type UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator int32

const (
	// Not specified.
	UserListLogicalRuleOperatorEnum_UNSPECIFIED UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator = 0
	// Used for return value only. Represents value unknown in this version.
	UserListLogicalRuleOperatorEnum_UNKNOWN UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator = 1
	// And - all of the operands.
	UserListLogicalRuleOperatorEnum_ALL UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator = 2
	// Or - at least one of the operands.
	UserListLogicalRuleOperatorEnum_ANY UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator = 3
	// Not - none of the operands.
	UserListLogicalRuleOperatorEnum_NONE UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator = 4
)

var UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "ALL",
	3: "ANY",
	4: "NONE",
}

var UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator_value = map[string]int32{
	"UNSPECIFIED": 0,
	"UNKNOWN":     1,
	"ALL":         2,
	"ANY":         3,
	"NONE":        4,
}

func (x UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator) String() string {
	return proto.EnumName(UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator_name, int32(x))
}

func (UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1f4edaf0f110fb85, []int{0, 0}
}

// The logical operator of the rule.
type UserListLogicalRuleOperatorEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserListLogicalRuleOperatorEnum) Reset()         { *m = UserListLogicalRuleOperatorEnum{} }
func (m *UserListLogicalRuleOperatorEnum) String() string { return proto.CompactTextString(m) }
func (*UserListLogicalRuleOperatorEnum) ProtoMessage()    {}
func (*UserListLogicalRuleOperatorEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_1f4edaf0f110fb85, []int{0}
}

func (m *UserListLogicalRuleOperatorEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserListLogicalRuleOperatorEnum.Unmarshal(m, b)
}
func (m *UserListLogicalRuleOperatorEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserListLogicalRuleOperatorEnum.Marshal(b, m, deterministic)
}
func (m *UserListLogicalRuleOperatorEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserListLogicalRuleOperatorEnum.Merge(m, src)
}
func (m *UserListLogicalRuleOperatorEnum) XXX_Size() int {
	return xxx_messageInfo_UserListLogicalRuleOperatorEnum.Size(m)
}
func (m *UserListLogicalRuleOperatorEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_UserListLogicalRuleOperatorEnum.DiscardUnknown(m)
}

var xxx_messageInfo_UserListLogicalRuleOperatorEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v3.enums.UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator", UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator_name, UserListLogicalRuleOperatorEnum_UserListLogicalRuleOperator_value)
	proto.RegisterType((*UserListLogicalRuleOperatorEnum)(nil), "google.ads.googleads.v3.enums.UserListLogicalRuleOperatorEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/enums/user_list_logical_rule_operator.proto", fileDescriptor_1f4edaf0f110fb85)
}

var fileDescriptor_1f4edaf0f110fb85 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x50, 0xb1, 0x6e, 0xea, 0x30,
	0x14, 0x7d, 0x04, 0xf4, 0x78, 0x32, 0xc3, 0x8b, 0x32, 0xb6, 0x45, 0x2d, 0x7c, 0x80, 0x33, 0x64,
	0x73, 0xa7, 0x40, 0x53, 0x84, 0x1a, 0x19, 0xd4, 0x0a, 0x50, 0xab, 0x48, 0x91, 0x4b, 0x2c, 0x2b,
	0x92, 0xb1, 0x23, 0xdf, 0x84, 0xa1, 0x9f, 0xd3, 0xb1, 0x9f, 0xd2, 0x4f, 0xe9, 0xde, 0xbd, 0x8a,
	0x0d, 0x6c, 0x65, 0xb1, 0x8e, 0x7c, 0xce, 0x3d, 0xe7, 0xde, 0x83, 0xa6, 0x42, 0x6b, 0x21, 0x79,
	0xc8, 0x0a, 0x08, 0x1d, 0x6c, 0xd1, 0x3e, 0x0a, 0xb9, 0x6a, 0x76, 0x10, 0x36, 0xc0, 0x4d, 0x2e,
	0x4b, 0xa8, 0x73, 0xa9, 0x45, 0xb9, 0x65, 0x32, 0x37, 0x8d, 0xe4, 0xb9, 0xae, 0xb8, 0x61, 0xb5,
	0x36, 0xb8, 0x32, 0xba, 0xd6, 0xc1, 0xd0, 0x4d, 0x62, 0x56, 0x00, 0x3e, 0x99, 0xe0, 0x7d, 0x84,
	0xad, 0xc9, 0xc5, 0xd5, 0x31, 0xa3, 0x2a, 0x43, 0xa6, 0x94, 0xae, 0x59, 0x5d, 0x6a, 0x05, 0x6e,
	0x78, 0xfc, 0x86, 0xae, 0x57, 0xc0, 0x4d, 0x5a, 0x42, 0x9d, 0xba, 0x8c, 0xc7, 0x46, 0xf2, 0xc5,
	0x21, 0x21, 0x51, 0xcd, 0x6e, 0xbc, 0x41, 0x97, 0x67, 0x24, 0xc1, 0x7f, 0x34, 0x58, 0xd1, 0xa7,
	0x65, 0x32, 0x9d, 0xdf, 0xcf, 0x93, 0x3b, 0xff, 0x4f, 0x30, 0x40, 0xfd, 0x15, 0x7d, 0xa0, 0x8b,
	0x0d, 0xf5, 0x3b, 0x41, 0x1f, 0x75, 0xe3, 0x34, 0xf5, 0x3d, 0x0b, 0xe8, 0xb3, 0xdf, 0x0d, 0xfe,
	0xa1, 0x1e, 0x5d, 0xd0, 0xc4, 0xef, 0x4d, 0xbe, 0x3b, 0x68, 0xb4, 0xd5, 0x3b, 0x7c, 0x76, 0xff,
	0xc9, 0xcd, 0x99, 0xf0, 0x65, 0x7b, 0xc3, 0xb2, 0xf3, 0x32, 0x39, 0x58, 0x08, 0x2d, 0x99, 0x12,
	0x58, 0x1b, 0x11, 0x0a, 0xae, 0xec, 0x85, 0xc7, 0x5e, 0xab, 0x12, 0x7e, 0xa9, 0xf9, 0xd6, 0xbe,
	0xef, 0x5e, 0x77, 0x16, 0xc7, 0x1f, 0xde, 0x70, 0xe6, 0xac, 0xe2, 0x02, 0xb0, 0x83, 0x2d, 0x5a,
	0x47, 0xb8, 0xed, 0x02, 0x3e, 0x8f, 0x7c, 0x16, 0x17, 0x90, 0x9d, 0xf8, 0x6c, 0x1d, 0x65, 0x96,
	0xff, 0xf2, 0x46, 0xee, 0x93, 0x90, 0xb8, 0x00, 0x42, 0x4e, 0x0a, 0x42, 0xd6, 0x11, 0x21, 0x56,
	0xf3, 0xfa, 0xd7, 0x2e, 0x16, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xcd, 0xbb, 0x94, 0xfe,
	0x01, 0x00, 0x00,
}

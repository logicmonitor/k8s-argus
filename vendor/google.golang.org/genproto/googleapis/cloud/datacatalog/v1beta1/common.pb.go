// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/datacatalog/v1beta1/common.proto

package datacatalog

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

// This enum describes all the possible systems that Data Catalog integrates
// with.
type IntegratedSystem int32

const (
	// Default unknown system.
	IntegratedSystem_INTEGRATED_SYSTEM_UNSPECIFIED IntegratedSystem = 0
	// BigQuery.
	IntegratedSystem_BIGQUERY IntegratedSystem = 1
	// Cloud Pub/Sub.
	IntegratedSystem_CLOUD_PUBSUB IntegratedSystem = 2
)

var IntegratedSystem_name = map[int32]string{
	0: "INTEGRATED_SYSTEM_UNSPECIFIED",
	1: "BIGQUERY",
	2: "CLOUD_PUBSUB",
}

var IntegratedSystem_value = map[string]int32{
	"INTEGRATED_SYSTEM_UNSPECIFIED": 0,
	"BIGQUERY":                      1,
	"CLOUD_PUBSUB":                  2,
}

func (x IntegratedSystem) String() string {
	return proto.EnumName(IntegratedSystem_name, int32(x))
}

func (IntegratedSystem) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c02185d1e0d17d2d, []int{0}
}

func init() {
	proto.RegisterEnum("google.cloud.datacatalog.v1beta1.IntegratedSystem", IntegratedSystem_name, IntegratedSystem_value)
}

func init() {
	proto.RegisterFile("google/cloud/datacatalog/v1beta1/common.proto", fileDescriptor_c02185d1e0d17d2d)
}

var fileDescriptor_c02185d1e0d17d2d = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xcf, 0xb1, 0x4b, 0x03, 0x31,
	0x14, 0x06, 0x70, 0xcf, 0x41, 0x24, 0x74, 0x08, 0xb7, 0x0b, 0x0a, 0x4e, 0x82, 0x09, 0xc5, 0xd1,
	0xc9, 0xeb, 0xc5, 0x12, 0xd4, 0x7a, 0x36, 0x97, 0xa1, 0x2e, 0xc7, 0x6b, 0x2e, 0x3c, 0x84, 0x24,
	0xaf, 0x5c, 0xa3, 0xe8, 0x7f, 0xee, 0x28, 0x5c, 0x3b, 0x74, 0x91, 0xae, 0xef, 0xfd, 0x1e, 0xef,
	0xfb, 0xd8, 0x2d, 0x12, 0x61, 0xf0, 0xd2, 0x05, 0xfa, 0xec, 0x65, 0x0f, 0x19, 0x1c, 0x64, 0x08,
	0x84, 0xf2, 0x6b, 0xba, 0xf6, 0x19, 0xa6, 0xd2, 0x51, 0x8c, 0x94, 0xc4, 0x66, 0xa0, 0x4c, 0xe5,
	0xe5, 0x8e, 0x8b, 0x91, 0x8b, 0x03, 0x2e, 0xf6, 0xfc, 0xc6, 0x32, 0xae, 0x53, 0xf6, 0x38, 0x40,
	0xf6, 0xbd, 0xf9, 0xd9, 0x66, 0x1f, 0xcb, 0x2b, 0x76, 0xa1, 0x17, 0xad, 0x9a, 0x2f, 0x1f, 0x5a,
	0x55, 0x77, 0x66, 0x65, 0x5a, 0xf5, 0xd2, 0xd9, 0x85, 0x69, 0xd4, 0x4c, 0x3f, 0x6a, 0x55, 0xf3,
	0x93, 0x72, 0xc2, 0xce, 0x2b, 0x3d, 0x7f, 0xb3, 0x6a, 0xb9, 0xe2, 0x45, 0xc9, 0xd9, 0x64, 0xf6,
	0xfc, 0x6a, 0xeb, 0xae, 0xb1, 0x95, 0xb1, 0x15, 0x3f, 0xad, 0xbe, 0xd9, 0xb5, 0xa3, 0x28, 0x8e,
	0xbd, 0x6f, 0x8a, 0xf7, 0xa7, 0xbd, 0x41, 0x0a, 0x90, 0x50, 0xd0, 0x80, 0x12, 0x7d, 0x1a, 0xe3,
	0xcb, 0xdd, 0x0a, 0x36, 0x1f, 0xdb, 0xff, 0x0b, 0xdf, 0x1f, 0xcc, 0x7e, 0x8b, 0x62, 0x7d, 0x36,
	0x9e, 0xde, 0xfd, 0x05, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x96, 0x3f, 0x0f, 0x2a, 0x01, 0x00, 0x00,
}

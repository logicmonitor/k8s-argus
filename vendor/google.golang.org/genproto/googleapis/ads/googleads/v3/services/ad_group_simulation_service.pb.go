// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/services/ad_group_simulation_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v3/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// Request message for [AdGroupSimulationService.GetAdGroupSimulation][google.ads.googleads.v3.services.AdGroupSimulationService.GetAdGroupSimulation].
type GetAdGroupSimulationRequest struct {
	// Required. The resource name of the ad group simulation to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdGroupSimulationRequest) Reset()         { *m = GetAdGroupSimulationRequest{} }
func (m *GetAdGroupSimulationRequest) String() string { return proto.CompactTextString(m) }
func (*GetAdGroupSimulationRequest) ProtoMessage()    {}
func (*GetAdGroupSimulationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0549cf58468ce402, []int{0}
}

func (m *GetAdGroupSimulationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdGroupSimulationRequest.Unmarshal(m, b)
}
func (m *GetAdGroupSimulationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdGroupSimulationRequest.Marshal(b, m, deterministic)
}
func (m *GetAdGroupSimulationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdGroupSimulationRequest.Merge(m, src)
}
func (m *GetAdGroupSimulationRequest) XXX_Size() int {
	return xxx_messageInfo_GetAdGroupSimulationRequest.Size(m)
}
func (m *GetAdGroupSimulationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdGroupSimulationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdGroupSimulationRequest proto.InternalMessageInfo

func (m *GetAdGroupSimulationRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetAdGroupSimulationRequest)(nil), "google.ads.googleads.v3.services.GetAdGroupSimulationRequest")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/services/ad_group_simulation_service.proto", fileDescriptor_0549cf58468ce402)
}

var fileDescriptor_0549cf58468ce402 = []byte{
	// 415 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x3f, 0x8b, 0xd4, 0x40,
	0x1c, 0x25, 0x39, 0x10, 0x0c, 0xda, 0x04, 0xc1, 0x25, 0xa7, 0xb8, 0x1c, 0x57, 0x1c, 0x57, 0xcc,
	0x80, 0x39, 0x2c, 0xe6, 0xb8, 0x62, 0xd2, 0xc4, 0x42, 0xe4, 0xb8, 0x83, 0x2d, 0x24, 0x10, 0xe6,
	0x32, 0x63, 0x1c, 0x48, 0x32, 0x71, 0x7e, 0x49, 0x1a, 0xb1, 0x11, 0xbf, 0xc1, 0x7e, 0x03, 0x4b,
	0x3f, 0xca, 0xb6, 0x56, 0x5a, 0x59, 0x58, 0xf9, 0x29, 0x24, 0x3b, 0x99, 0xec, 0xae, 0x6b, 0xdc,
	0xee, 0x31, 0xef, 0xfd, 0xde, 0xfb, 0xfd, 0x19, 0x2f, 0xca, 0x95, 0xca, 0x0b, 0x81, 0x19, 0x07,
	0x6c, 0x60, 0x8f, 0xba, 0x10, 0x83, 0xd0, 0x9d, 0xcc, 0x04, 0x60, 0xc6, 0xd3, 0x5c, 0xab, 0xb6,
	0x4e, 0x41, 0x96, 0x6d, 0xc1, 0x1a, 0xa9, 0xaa, 0x74, 0x20, 0x51, 0xad, 0x55, 0xa3, 0xfc, 0xb9,
	0x29, 0x44, 0x8c, 0x03, 0x1a, 0x3d, 0x50, 0x17, 0x22, 0xeb, 0x11, 0x5c, 0x4e, 0xa5, 0x68, 0x01,
	0xaa, 0xd5, 0x13, 0x31, 0xc6, 0x3e, 0x78, 0x62, 0x8b, 0x6b, 0x89, 0x59, 0x55, 0xa9, 0x66, 0x4d,
	0xc2, 0xc0, 0x3e, 0xde, 0x62, 0xb3, 0x42, 0x8a, 0xaa, 0x19, 0x88, 0x67, 0x5b, 0xc4, 0x5b, 0x29,
	0x0a, 0x9e, 0xde, 0x89, 0x77, 0xac, 0x93, 0x4a, 0x1b, 0xc1, 0x49, 0xec, 0x1d, 0xc7, 0xa2, 0xa1,
	0x3c, 0xee, 0x63, 0x6f, 0xc7, 0xd4, 0x1b, 0xf1, 0xbe, 0x15, 0xd0, 0xf8, 0x67, 0xde, 0x43, 0xdb,
	0x5d, 0x5a, 0xb1, 0x52, 0xcc, 0x9c, 0xb9, 0x73, 0x76, 0x3f, 0x3a, 0xfa, 0x49, 0xdd, 0x9b, 0x07,
	0x96, 0x79, 0xcd, 0x4a, 0xf1, 0x7c, 0xe9, 0x7a, 0xb3, 0x3d, 0x9b, 0x5b, 0x33, 0xbb, 0xff, 0xdd,
	0xf1, 0x1e, 0xfd, 0x2b, 0xc6, 0xbf, 0x42, 0x87, 0xd6, 0x86, 0xfe, 0xd3, 0x5e, 0x70, 0x31, 0x59,
	0x3e, 0xee, 0x14, 0xed, 0x15, 0x9f, 0xbc, 0xfa, 0x41, 0x77, 0xa7, 0xfa, 0xf4, 0xed, 0xd7, 0xd2,
	0x7d, 0xe1, 0x5f, 0xf4, 0xc7, 0xf8, 0xb0, 0xc3, 0x5c, 0x65, 0x2d, 0x34, 0xaa, 0x14, 0x1a, 0xf0,
	0x39, 0x66, 0x7f, 0x3b, 0x01, 0x3e, 0xff, 0x18, 0x1c, 0xaf, 0xe8, 0x6c, 0x13, 0x3d, 0xa0, 0x5a,
	0x02, 0xca, 0x54, 0x19, 0x7d, 0x76, 0xbd, 0xd3, 0x4c, 0x95, 0x07, 0xa7, 0x8c, 0x9e, 0x4e, 0xed,
	0xee, 0xba, 0x3f, 0xd3, 0xb5, 0xf3, 0xe6, 0xe5, 0x60, 0x91, 0xab, 0x82, 0x55, 0x39, 0x52, 0x3a,
	0xc7, 0xb9, 0xa8, 0xd6, 0x47, 0xc4, 0x9b, 0xd0, 0xe9, 0x2f, 0x7c, 0x69, 0xc1, 0x17, 0xf7, 0x28,
	0xa6, 0xf4, 0xab, 0x3b, 0x8f, 0x8d, 0x21, 0xe5, 0x80, 0x0c, 0xec, 0xd1, 0x22, 0x44, 0x43, 0x30,
	0xac, 0xac, 0x24, 0xa1, 0x1c, 0x92, 0x51, 0x92, 0x2c, 0xc2, 0xc4, 0x4a, 0x7e, 0xbb, 0xa7, 0xe6,
	0x9d, 0x10, 0xca, 0x81, 0x90, 0x51, 0x44, 0xc8, 0x22, 0x24, 0xc4, 0xca, 0xee, 0xee, 0xad, 0xfb,
	0x0c, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x35, 0x6d, 0x4b, 0x69, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdGroupSimulationServiceClient is the client API for AdGroupSimulationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdGroupSimulationServiceClient interface {
	// Returns the requested ad group simulation in full detail.
	GetAdGroupSimulation(ctx context.Context, in *GetAdGroupSimulationRequest, opts ...grpc.CallOption) (*resources.AdGroupSimulation, error)
}

type adGroupSimulationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdGroupSimulationServiceClient(cc grpc.ClientConnInterface) AdGroupSimulationServiceClient {
	return &adGroupSimulationServiceClient{cc}
}

func (c *adGroupSimulationServiceClient) GetAdGroupSimulation(ctx context.Context, in *GetAdGroupSimulationRequest, opts ...grpc.CallOption) (*resources.AdGroupSimulation, error) {
	out := new(resources.AdGroupSimulation)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v3.services.AdGroupSimulationService/GetAdGroupSimulation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdGroupSimulationServiceServer is the server API for AdGroupSimulationService service.
type AdGroupSimulationServiceServer interface {
	// Returns the requested ad group simulation in full detail.
	GetAdGroupSimulation(context.Context, *GetAdGroupSimulationRequest) (*resources.AdGroupSimulation, error)
}

// UnimplementedAdGroupSimulationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAdGroupSimulationServiceServer struct {
}

func (*UnimplementedAdGroupSimulationServiceServer) GetAdGroupSimulation(ctx context.Context, req *GetAdGroupSimulationRequest) (*resources.AdGroupSimulation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdGroupSimulation not implemented")
}

func RegisterAdGroupSimulationServiceServer(s *grpc.Server, srv AdGroupSimulationServiceServer) {
	s.RegisterService(&_AdGroupSimulationService_serviceDesc, srv)
}

func _AdGroupSimulationService_GetAdGroupSimulation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdGroupSimulationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdGroupSimulationServiceServer).GetAdGroupSimulation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v3.services.AdGroupSimulationService/GetAdGroupSimulation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdGroupSimulationServiceServer).GetAdGroupSimulation(ctx, req.(*GetAdGroupSimulationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdGroupSimulationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v3.services.AdGroupSimulationService",
	HandlerType: (*AdGroupSimulationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAdGroupSimulation",
			Handler:    _AdGroupSimulationService_GetAdGroupSimulation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v3/services/ad_group_simulation_service.proto",
}

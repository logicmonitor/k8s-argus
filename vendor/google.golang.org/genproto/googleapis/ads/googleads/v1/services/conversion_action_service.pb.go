// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/services/conversion_action_service.proto

package services

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	resources "google.golang.org/genproto/googleapis/ads/googleads/v1/resources"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	status "google.golang.org/genproto/googleapis/rpc/status"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status1 "google.golang.org/grpc/status"
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

// Request message for [ConversionActionService.GetConversionAction][google.ads.googleads.v1.services.ConversionActionService.GetConversionAction].
type GetConversionActionRequest struct {
	// The resource name of the conversion action to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetConversionActionRequest) Reset()         { *m = GetConversionActionRequest{} }
func (m *GetConversionActionRequest) String() string { return proto.CompactTextString(m) }
func (*GetConversionActionRequest) ProtoMessage()    {}
func (*GetConversionActionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c1381831f293da1, []int{0}
}

func (m *GetConversionActionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetConversionActionRequest.Unmarshal(m, b)
}
func (m *GetConversionActionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetConversionActionRequest.Marshal(b, m, deterministic)
}
func (m *GetConversionActionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetConversionActionRequest.Merge(m, src)
}
func (m *GetConversionActionRequest) XXX_Size() int {
	return xxx_messageInfo_GetConversionActionRequest.Size(m)
}
func (m *GetConversionActionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetConversionActionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetConversionActionRequest proto.InternalMessageInfo

func (m *GetConversionActionRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

// Request message for [ConversionActionService.MutateConversionActions][google.ads.googleads.v1.services.ConversionActionService.MutateConversionActions].
type MutateConversionActionsRequest struct {
	// The ID of the customer whose conversion actions are being modified.
	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	// The list of operations to perform on individual conversion actions.
	Operations []*ConversionActionOperation `protobuf:"bytes,2,rep,name=operations,proto3" json:"operations,omitempty"`
	// If true, successful operations will be carried out and invalid
	// operations will return errors. If false, all operations will be carried
	// out in one transaction if and only if they are all valid.
	// Default is false.
	PartialFailure bool `protobuf:"varint,3,opt,name=partial_failure,json=partialFailure,proto3" json:"partial_failure,omitempty"`
	// If true, the request is validated but not executed. Only errors are
	// returned, not results.
	ValidateOnly         bool     `protobuf:"varint,4,opt,name=validate_only,json=validateOnly,proto3" json:"validate_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateConversionActionsRequest) Reset()         { *m = MutateConversionActionsRequest{} }
func (m *MutateConversionActionsRequest) String() string { return proto.CompactTextString(m) }
func (*MutateConversionActionsRequest) ProtoMessage()    {}
func (*MutateConversionActionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c1381831f293da1, []int{1}
}

func (m *MutateConversionActionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateConversionActionsRequest.Unmarshal(m, b)
}
func (m *MutateConversionActionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateConversionActionsRequest.Marshal(b, m, deterministic)
}
func (m *MutateConversionActionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateConversionActionsRequest.Merge(m, src)
}
func (m *MutateConversionActionsRequest) XXX_Size() int {
	return xxx_messageInfo_MutateConversionActionsRequest.Size(m)
}
func (m *MutateConversionActionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateConversionActionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MutateConversionActionsRequest proto.InternalMessageInfo

func (m *MutateConversionActionsRequest) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *MutateConversionActionsRequest) GetOperations() []*ConversionActionOperation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *MutateConversionActionsRequest) GetPartialFailure() bool {
	if m != nil {
		return m.PartialFailure
	}
	return false
}

func (m *MutateConversionActionsRequest) GetValidateOnly() bool {
	if m != nil {
		return m.ValidateOnly
	}
	return false
}

// A single operation (create, update, remove) on a conversion action.
type ConversionActionOperation struct {
	// FieldMask that determines which resource fields are modified in an update.
	UpdateMask *field_mask.FieldMask `protobuf:"bytes,4,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	// The mutate operation.
	//
	// Types that are valid to be assigned to Operation:
	//	*ConversionActionOperation_Create
	//	*ConversionActionOperation_Update
	//	*ConversionActionOperation_Remove
	Operation            isConversionActionOperation_Operation `protobuf_oneof:"operation"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *ConversionActionOperation) Reset()         { *m = ConversionActionOperation{} }
func (m *ConversionActionOperation) String() string { return proto.CompactTextString(m) }
func (*ConversionActionOperation) ProtoMessage()    {}
func (*ConversionActionOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c1381831f293da1, []int{2}
}

func (m *ConversionActionOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConversionActionOperation.Unmarshal(m, b)
}
func (m *ConversionActionOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConversionActionOperation.Marshal(b, m, deterministic)
}
func (m *ConversionActionOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConversionActionOperation.Merge(m, src)
}
func (m *ConversionActionOperation) XXX_Size() int {
	return xxx_messageInfo_ConversionActionOperation.Size(m)
}
func (m *ConversionActionOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_ConversionActionOperation.DiscardUnknown(m)
}

var xxx_messageInfo_ConversionActionOperation proto.InternalMessageInfo

func (m *ConversionActionOperation) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type isConversionActionOperation_Operation interface {
	isConversionActionOperation_Operation()
}

type ConversionActionOperation_Create struct {
	Create *resources.ConversionAction `protobuf:"bytes,1,opt,name=create,proto3,oneof"`
}

type ConversionActionOperation_Update struct {
	Update *resources.ConversionAction `protobuf:"bytes,2,opt,name=update,proto3,oneof"`
}

type ConversionActionOperation_Remove struct {
	Remove string `protobuf:"bytes,3,opt,name=remove,proto3,oneof"`
}

func (*ConversionActionOperation_Create) isConversionActionOperation_Operation() {}

func (*ConversionActionOperation_Update) isConversionActionOperation_Operation() {}

func (*ConversionActionOperation_Remove) isConversionActionOperation_Operation() {}

func (m *ConversionActionOperation) GetOperation() isConversionActionOperation_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *ConversionActionOperation) GetCreate() *resources.ConversionAction {
	if x, ok := m.GetOperation().(*ConversionActionOperation_Create); ok {
		return x.Create
	}
	return nil
}

func (m *ConversionActionOperation) GetUpdate() *resources.ConversionAction {
	if x, ok := m.GetOperation().(*ConversionActionOperation_Update); ok {
		return x.Update
	}
	return nil
}

func (m *ConversionActionOperation) GetRemove() string {
	if x, ok := m.GetOperation().(*ConversionActionOperation_Remove); ok {
		return x.Remove
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ConversionActionOperation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ConversionActionOperation_Create)(nil),
		(*ConversionActionOperation_Update)(nil),
		(*ConversionActionOperation_Remove)(nil),
	}
}

// Response message for [ConversionActionService.MutateConversionActions][google.ads.googleads.v1.services.ConversionActionService.MutateConversionActions].
type MutateConversionActionsResponse struct {
	// Errors that pertain to operation failures in the partial failure mode.
	// Returned only when partial_failure = true and all errors occur inside the
	// operations. If any errors occur outside the operations (e.g. auth errors),
	// we return an RPC level error.
	PartialFailureError *status.Status `protobuf:"bytes,3,opt,name=partial_failure_error,json=partialFailureError,proto3" json:"partial_failure_error,omitempty"`
	// All results for the mutate.
	Results              []*MutateConversionActionResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *MutateConversionActionsResponse) Reset()         { *m = MutateConversionActionsResponse{} }
func (m *MutateConversionActionsResponse) String() string { return proto.CompactTextString(m) }
func (*MutateConversionActionsResponse) ProtoMessage()    {}
func (*MutateConversionActionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c1381831f293da1, []int{3}
}

func (m *MutateConversionActionsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateConversionActionsResponse.Unmarshal(m, b)
}
func (m *MutateConversionActionsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateConversionActionsResponse.Marshal(b, m, deterministic)
}
func (m *MutateConversionActionsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateConversionActionsResponse.Merge(m, src)
}
func (m *MutateConversionActionsResponse) XXX_Size() int {
	return xxx_messageInfo_MutateConversionActionsResponse.Size(m)
}
func (m *MutateConversionActionsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateConversionActionsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MutateConversionActionsResponse proto.InternalMessageInfo

func (m *MutateConversionActionsResponse) GetPartialFailureError() *status.Status {
	if m != nil {
		return m.PartialFailureError
	}
	return nil
}

func (m *MutateConversionActionsResponse) GetResults() []*MutateConversionActionResult {
	if m != nil {
		return m.Results
	}
	return nil
}

// The result for the conversion action mutate.
type MutateConversionActionResult struct {
	// Returned for successful operations.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MutateConversionActionResult) Reset()         { *m = MutateConversionActionResult{} }
func (m *MutateConversionActionResult) String() string { return proto.CompactTextString(m) }
func (*MutateConversionActionResult) ProtoMessage()    {}
func (*MutateConversionActionResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c1381831f293da1, []int{4}
}

func (m *MutateConversionActionResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MutateConversionActionResult.Unmarshal(m, b)
}
func (m *MutateConversionActionResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MutateConversionActionResult.Marshal(b, m, deterministic)
}
func (m *MutateConversionActionResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MutateConversionActionResult.Merge(m, src)
}
func (m *MutateConversionActionResult) XXX_Size() int {
	return xxx_messageInfo_MutateConversionActionResult.Size(m)
}
func (m *MutateConversionActionResult) XXX_DiscardUnknown() {
	xxx_messageInfo_MutateConversionActionResult.DiscardUnknown(m)
}

var xxx_messageInfo_MutateConversionActionResult proto.InternalMessageInfo

func (m *MutateConversionActionResult) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetConversionActionRequest)(nil), "google.ads.googleads.v1.services.GetConversionActionRequest")
	proto.RegisterType((*MutateConversionActionsRequest)(nil), "google.ads.googleads.v1.services.MutateConversionActionsRequest")
	proto.RegisterType((*ConversionActionOperation)(nil), "google.ads.googleads.v1.services.ConversionActionOperation")
	proto.RegisterType((*MutateConversionActionsResponse)(nil), "google.ads.googleads.v1.services.MutateConversionActionsResponse")
	proto.RegisterType((*MutateConversionActionResult)(nil), "google.ads.googleads.v1.services.MutateConversionActionResult")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/services/conversion_action_service.proto", fileDescriptor_7c1381831f293da1)
}

var fileDescriptor_7c1381831f293da1 = []byte{
	// 728 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0x4f, 0x4f, 0xd4, 0x4c,
	0x18, 0xc0, 0xdf, 0x76, 0xdf, 0xf0, 0xbe, 0xcc, 0xa2, 0x26, 0x43, 0x0c, 0x6b, 0x25, 0xb2, 0xa9,
	0x24, 0x92, 0x3d, 0xb4, 0xd9, 0x25, 0x9a, 0xd8, 0x15, 0x43, 0x21, 0x02, 0x1e, 0x10, 0x52, 0x12,
	0x62, 0x74, 0x93, 0x66, 0x68, 0x87, 0x4d, 0x43, 0xdb, 0xa9, 0x33, 0xd3, 0x4d, 0x08, 0xe1, 0xa2,
	0x1f, 0xc1, 0xb3, 0x17, 0x8f, 0x7e, 0x0d, 0xe3, 0x85, 0xab, 0x9f, 0x40, 0xe3, 0xc9, 0xf8, 0x21,
	0xcc, 0x74, 0x3a, 0x0b, 0x2c, 0xd4, 0x35, 0x70, 0xda, 0xa7, 0xf3, 0x3c, 0xf3, 0x7b, 0xfe, 0xce,
	0xb3, 0x60, 0xb9, 0x4f, 0x48, 0x3f, 0xc6, 0x36, 0x0a, 0x99, 0x2d, 0x45, 0x21, 0x0d, 0xda, 0x36,
	0xc3, 0x74, 0x10, 0x05, 0x98, 0xd9, 0x01, 0x49, 0x07, 0x98, 0xb2, 0x88, 0xa4, 0x3e, 0x0a, 0xb8,
	0xf8, 0x29, 0x55, 0x56, 0x46, 0x09, 0x27, 0xb0, 0x29, 0xaf, 0x59, 0x28, 0x64, 0xd6, 0x90, 0x60,
	0x0d, 0xda, 0x96, 0x22, 0x18, 0x8f, 0xab, 0x7c, 0x50, 0xcc, 0x48, 0x4e, 0x2f, 0x75, 0x22, 0xe1,
	0xc6, 0xac, 0xba, 0x9a, 0x45, 0x36, 0x4a, 0x53, 0xc2, 0x91, 0x50, 0xb2, 0x52, 0x5b, 0xba, 0xb6,
	0x8b, 0xaf, 0xbd, 0x7c, 0xdf, 0xde, 0x8f, 0x70, 0x1c, 0xfa, 0x09, 0x62, 0x07, 0xa5, 0xc5, 0x4c,
	0x69, 0x41, 0xb3, 0xc0, 0x66, 0x1c, 0xf1, 0x9c, 0x8d, 0x28, 0x04, 0x38, 0x88, 0x23, 0x9c, 0x72,
	0xa9, 0x30, 0x5d, 0x60, 0xac, 0x63, 0xbe, 0x3a, 0x8c, 0xc7, 0x2d, 0xc2, 0xf1, 0xf0, 0x9b, 0x1c,
	0x33, 0x0e, 0xef, 0x83, 0x1b, 0x2a, 0x68, 0x3f, 0x45, 0x09, 0x6e, 0x68, 0x4d, 0x6d, 0x61, 0xd2,
	0x9b, 0x52, 0x87, 0x2f, 0x50, 0x82, 0xcd, 0x5f, 0x1a, 0xb8, 0xb7, 0x99, 0x73, 0xc4, 0xf1, 0x28,
	0x86, 0x29, 0xce, 0x1c, 0xa8, 0x07, 0x39, 0xe3, 0x24, 0xc1, 0xd4, 0x8f, 0xc2, 0x92, 0x02, 0xd4,
	0xd1, 0xf3, 0x10, 0xbe, 0x06, 0x80, 0x64, 0x98, 0xca, 0x74, 0x1b, 0x7a, 0xb3, 0xb6, 0x50, 0xef,
	0x74, 0xad, 0x71, 0xa5, 0xb6, 0x46, 0x1d, 0x6e, 0x29, 0x86, 0x77, 0x06, 0x07, 0x1f, 0x80, 0x5b,
	0x19, 0xa2, 0x3c, 0x42, 0xb1, 0xbf, 0x8f, 0xa2, 0x38, 0xa7, 0xb8, 0x51, 0x6b, 0x6a, 0x0b, 0xff,
	0x7b, 0x37, 0xcb, 0xe3, 0x35, 0x79, 0x2a, 0xd2, 0x1d, 0xa0, 0x38, 0x0a, 0x11, 0xc7, 0x3e, 0x49,
	0xe3, 0xc3, 0xc6, 0xbf, 0x85, 0xd9, 0x94, 0x3a, 0xdc, 0x4a, 0xe3, 0x43, 0xf3, 0x83, 0x0e, 0xee,
	0x54, 0xfa, 0x85, 0x5d, 0x50, 0xcf, 0xb3, 0x02, 0x20, 0xda, 0x52, 0x00, 0xea, 0x1d, 0x43, 0x65,
	0xa2, 0x3a, 0x67, 0xad, 0x89, 0xce, 0x6d, 0x22, 0x76, 0xe0, 0x01, 0x69, 0x2e, 0x64, 0xb8, 0x09,
	0x26, 0x02, 0x8a, 0x11, 0x97, 0x75, 0xae, 0x77, 0x16, 0x2b, 0x2b, 0x30, 0x1c, 0xa5, 0x0b, 0x25,
	0xd8, 0xf8, 0xc7, 0x2b, 0x21, 0x02, 0x27, 0xe1, 0x0d, 0xfd, 0x5a, 0x38, 0x09, 0x81, 0x0d, 0x30,
	0x41, 0x71, 0x42, 0x06, 0xb2, 0x7a, 0x93, 0x42, 0x23, 0xbf, 0x57, 0xea, 0x60, 0x72, 0x58, 0x6e,
	0xf3, 0xb3, 0x06, 0xe6, 0x2a, 0xc7, 0x81, 0x65, 0x24, 0x65, 0x18, 0xae, 0x81, 0xdb, 0x23, 0x1d,
	0xf1, 0x31, 0xa5, 0x84, 0x16, 0xe4, 0x7a, 0x07, 0xaa, 0x40, 0x69, 0x16, 0x58, 0x3b, 0xc5, 0x1c,
	0x7b, 0xd3, 0xe7, 0x7b, 0xf5, 0x4c, 0x98, 0xc3, 0x97, 0xe0, 0x3f, 0x8a, 0x59, 0x1e, 0x73, 0x35,
	0x33, 0x4f, 0xc7, 0xcf, 0xcc, 0xe5, 0xb1, 0x79, 0x05, 0xc6, 0x53, 0x38, 0x73, 0x15, 0xcc, 0xfe,
	0xc9, 0xf0, 0xaf, 0x5e, 0x46, 0xe7, 0x7b, 0x0d, 0xcc, 0x8c, 0xde, 0xdf, 0x91, 0x71, 0xc0, 0x2f,
	0x1a, 0x98, 0xbe, 0xe4, 0xe5, 0xc1, 0x27, 0xe3, 0x33, 0xa8, 0x7e, 0xb0, 0xc6, 0x55, 0x5a, 0x6c,
	0x76, 0xdf, 0x7e, 0xfd, 0xf1, 0x5e, 0x7f, 0x08, 0x17, 0xc5, 0x92, 0x3a, 0x3a, 0x97, 0xd6, 0x92,
	0x7a, 0xa3, 0xcc, 0x6e, 0x9d, 0xd9, 0x5a, 0x65, 0x3f, 0xed, 0xd6, 0x31, 0xfc, 0xa6, 0x81, 0x99,
	0x8a, 0x76, 0xc3, 0xe5, 0xab, 0x76, 0x43, 0x2d, 0x0e, 0xc3, 0xbd, 0x06, 0x41, 0xce, 0x9a, 0xe9,
	0x16, 0xd9, 0x75, 0xcd, 0x47, 0x22, 0xbb, 0xd3, 0x74, 0x8e, 0xce, 0x2c, 0xa4, 0xa5, 0xd6, 0xf1,
	0xc5, 0xe4, 0x9c, 0xa4, 0x00, 0x3b, 0x5a, 0xcb, 0xb8, 0x7b, 0xe2, 0x36, 0x4e, 0x9d, 0x97, 0x52,
	0x16, 0x31, 0x2b, 0x20, 0xc9, 0xca, 0x3b, 0x1d, 0xcc, 0x07, 0x24, 0x19, 0x1b, 0xe8, 0xca, 0x6c,
	0xc5, 0x28, 0x6c, 0x8b, 0xa5, 0xb0, 0xad, 0xbd, 0xda, 0x28, 0x09, 0x7d, 0x12, 0xa3, 0xb4, 0x6f,
	0x11, 0xda, 0xb7, 0xfb, 0x38, 0x2d, 0x56, 0x86, 0x7d, 0xea, 0xb3, 0xfa, 0xaf, 0xab, 0xab, 0x84,
	0x8f, 0x7a, 0x6d, 0xdd, 0x75, 0x3f, 0xe9, 0xcd, 0x75, 0x09, 0x74, 0x43, 0x66, 0x49, 0x51, 0x48,
	0xbb, 0x6d, 0xab, 0x74, 0xcc, 0x4e, 0x94, 0x49, 0xcf, 0x0d, 0x59, 0x6f, 0x68, 0xd2, 0xdb, 0x6d,
	0xf7, 0x94, 0xc9, 0x4f, 0x7d, 0x5e, 0x9e, 0x3b, 0x8e, 0x1b, 0x32, 0xc7, 0x19, 0x1a, 0x39, 0xce,
	0x6e, 0xdb, 0x71, 0x94, 0xd9, 0xde, 0x44, 0x11, 0xe7, 0xe2, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xec, 0x7a, 0x1f, 0xe7, 0x61, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ConversionActionServiceClient is the client API for ConversionActionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConversionActionServiceClient interface {
	// Returns the requested conversion action.
	GetConversionAction(ctx context.Context, in *GetConversionActionRequest, opts ...grpc.CallOption) (*resources.ConversionAction, error)
	// Creates, updates or removes conversion actions. Operation statuses are
	// returned.
	MutateConversionActions(ctx context.Context, in *MutateConversionActionsRequest, opts ...grpc.CallOption) (*MutateConversionActionsResponse, error)
}

type conversionActionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConversionActionServiceClient(cc grpc.ClientConnInterface) ConversionActionServiceClient {
	return &conversionActionServiceClient{cc}
}

func (c *conversionActionServiceClient) GetConversionAction(ctx context.Context, in *GetConversionActionRequest, opts ...grpc.CallOption) (*resources.ConversionAction, error) {
	out := new(resources.ConversionAction)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.ConversionActionService/GetConversionAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *conversionActionServiceClient) MutateConversionActions(ctx context.Context, in *MutateConversionActionsRequest, opts ...grpc.CallOption) (*MutateConversionActionsResponse, error) {
	out := new(MutateConversionActionsResponse)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v1.services.ConversionActionService/MutateConversionActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConversionActionServiceServer is the server API for ConversionActionService service.
type ConversionActionServiceServer interface {
	// Returns the requested conversion action.
	GetConversionAction(context.Context, *GetConversionActionRequest) (*resources.ConversionAction, error)
	// Creates, updates or removes conversion actions. Operation statuses are
	// returned.
	MutateConversionActions(context.Context, *MutateConversionActionsRequest) (*MutateConversionActionsResponse, error)
}

// UnimplementedConversionActionServiceServer can be embedded to have forward compatible implementations.
type UnimplementedConversionActionServiceServer struct {
}

func (*UnimplementedConversionActionServiceServer) GetConversionAction(ctx context.Context, req *GetConversionActionRequest) (*resources.ConversionAction, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method GetConversionAction not implemented")
}
func (*UnimplementedConversionActionServiceServer) MutateConversionActions(ctx context.Context, req *MutateConversionActionsRequest) (*MutateConversionActionsResponse, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method MutateConversionActions not implemented")
}

func RegisterConversionActionServiceServer(s *grpc.Server, srv ConversionActionServiceServer) {
	s.RegisterService(&_ConversionActionService_serviceDesc, srv)
}

func _ConversionActionService_GetConversionAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConversionActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversionActionServiceServer).GetConversionAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.ConversionActionService/GetConversionAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversionActionServiceServer).GetConversionAction(ctx, req.(*GetConversionActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConversionActionService_MutateConversionActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateConversionActionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversionActionServiceServer).MutateConversionActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v1.services.ConversionActionService/MutateConversionActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversionActionServiceServer).MutateConversionActions(ctx, req.(*MutateConversionActionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ConversionActionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v1.services.ConversionActionService",
	HandlerType: (*ConversionActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConversionAction",
			Handler:    _ConversionActionService_GetConversionAction_Handler,
		},
		{
			MethodName: "MutateConversionActions",
			Handler:    _ConversionActionService_MutateConversionActions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v1/services/conversion_action_service.proto",
}

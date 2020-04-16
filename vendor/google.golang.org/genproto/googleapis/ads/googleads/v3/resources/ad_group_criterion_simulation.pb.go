// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v3/resources/ad_group_criterion_simulation.proto

package resources

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	common "google.golang.org/genproto/googleapis/ads/googleads/v3/common"
	enums "google.golang.org/genproto/googleapis/ads/googleads/v3/enums"
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

// An ad group criterion simulation. Supported combinations of advertising
// channel type, criterion type, simulation type, and simulation modification
// method are detailed below respectively.
//
// 1. DISPLAY - KEYWORD - CPC_BID - UNIFORM
// 2. SEARCH - KEYWORD - CPC_BID - UNIFORM
// 3. SHOPPING - LISTING_GROUP - CPC_BID - UNIFORM
type AdGroupCriterionSimulation struct {
	// The resource name of the ad group criterion simulation.
	// Ad group criterion simulation resource names have the form:
	//
	// `customers/{customer_id}/adGroupCriterionSimulations/{ad_group_id}~{criterion_id}~{type}~{modification_method}~{start_date}~{end_date}`
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	// AdGroup ID of the simulation.
	AdGroupId *wrappers.Int64Value `protobuf:"bytes,2,opt,name=ad_group_id,json=adGroupId,proto3" json:"ad_group_id,omitempty"`
	// Criterion ID of the simulation.
	CriterionId *wrappers.Int64Value `protobuf:"bytes,3,opt,name=criterion_id,json=criterionId,proto3" json:"criterion_id,omitempty"`
	// The field that the simulation modifies.
	Type enums.SimulationTypeEnum_SimulationType `protobuf:"varint,4,opt,name=type,proto3,enum=google.ads.googleads.v3.enums.SimulationTypeEnum_SimulationType" json:"type,omitempty"`
	// How the simulation modifies the field.
	ModificationMethod enums.SimulationModificationMethodEnum_SimulationModificationMethod `protobuf:"varint,5,opt,name=modification_method,json=modificationMethod,proto3,enum=google.ads.googleads.v3.enums.SimulationModificationMethodEnum_SimulationModificationMethod" json:"modification_method,omitempty"`
	// First day on which the simulation is based, in YYYY-MM-DD format.
	StartDate *wrappers.StringValue `protobuf:"bytes,6,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	// Last day on which the simulation is based, in YYYY-MM-DD format.
	EndDate *wrappers.StringValue `protobuf:"bytes,7,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	// List of simulation points.
	//
	// Types that are valid to be assigned to PointList:
	//	*AdGroupCriterionSimulation_CpcBidPointList
	PointList            isAdGroupCriterionSimulation_PointList `protobuf_oneof:"point_list"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *AdGroupCriterionSimulation) Reset()         { *m = AdGroupCriterionSimulation{} }
func (m *AdGroupCriterionSimulation) String() string { return proto.CompactTextString(m) }
func (*AdGroupCriterionSimulation) ProtoMessage()    {}
func (*AdGroupCriterionSimulation) Descriptor() ([]byte, []int) {
	return fileDescriptor_79b2fc18b4f92e57, []int{0}
}

func (m *AdGroupCriterionSimulation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdGroupCriterionSimulation.Unmarshal(m, b)
}
func (m *AdGroupCriterionSimulation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdGroupCriterionSimulation.Marshal(b, m, deterministic)
}
func (m *AdGroupCriterionSimulation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdGroupCriterionSimulation.Merge(m, src)
}
func (m *AdGroupCriterionSimulation) XXX_Size() int {
	return xxx_messageInfo_AdGroupCriterionSimulation.Size(m)
}
func (m *AdGroupCriterionSimulation) XXX_DiscardUnknown() {
	xxx_messageInfo_AdGroupCriterionSimulation.DiscardUnknown(m)
}

var xxx_messageInfo_AdGroupCriterionSimulation proto.InternalMessageInfo

func (m *AdGroupCriterionSimulation) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func (m *AdGroupCriterionSimulation) GetAdGroupId() *wrappers.Int64Value {
	if m != nil {
		return m.AdGroupId
	}
	return nil
}

func (m *AdGroupCriterionSimulation) GetCriterionId() *wrappers.Int64Value {
	if m != nil {
		return m.CriterionId
	}
	return nil
}

func (m *AdGroupCriterionSimulation) GetType() enums.SimulationTypeEnum_SimulationType {
	if m != nil {
		return m.Type
	}
	return enums.SimulationTypeEnum_UNSPECIFIED
}

func (m *AdGroupCriterionSimulation) GetModificationMethod() enums.SimulationModificationMethodEnum_SimulationModificationMethod {
	if m != nil {
		return m.ModificationMethod
	}
	return enums.SimulationModificationMethodEnum_UNSPECIFIED
}

func (m *AdGroupCriterionSimulation) GetStartDate() *wrappers.StringValue {
	if m != nil {
		return m.StartDate
	}
	return nil
}

func (m *AdGroupCriterionSimulation) GetEndDate() *wrappers.StringValue {
	if m != nil {
		return m.EndDate
	}
	return nil
}

type isAdGroupCriterionSimulation_PointList interface {
	isAdGroupCriterionSimulation_PointList()
}

type AdGroupCriterionSimulation_CpcBidPointList struct {
	CpcBidPointList *common.CpcBidSimulationPointList `protobuf:"bytes,8,opt,name=cpc_bid_point_list,json=cpcBidPointList,proto3,oneof"`
}

func (*AdGroupCriterionSimulation_CpcBidPointList) isAdGroupCriterionSimulation_PointList() {}

func (m *AdGroupCriterionSimulation) GetPointList() isAdGroupCriterionSimulation_PointList {
	if m != nil {
		return m.PointList
	}
	return nil
}

func (m *AdGroupCriterionSimulation) GetCpcBidPointList() *common.CpcBidSimulationPointList {
	if x, ok := m.GetPointList().(*AdGroupCriterionSimulation_CpcBidPointList); ok {
		return x.CpcBidPointList
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AdGroupCriterionSimulation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AdGroupCriterionSimulation_CpcBidPointList)(nil),
	}
}

func init() {
	proto.RegisterType((*AdGroupCriterionSimulation)(nil), "google.ads.googleads.v3.resources.AdGroupCriterionSimulation")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v3/resources/ad_group_criterion_simulation.proto", fileDescriptor_79b2fc18b4f92e57)
}

var fileDescriptor_79b2fc18b4f92e57 = []byte{
	// 605 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdf, 0x4e, 0xd4, 0x4e,
	0x14, 0xfe, 0x75, 0xe1, 0xc7, 0x9f, 0x01, 0x35, 0xa9, 0x37, 0x15, 0x89, 0x82, 0x86, 0x84, 0xab,
	0x99, 0x84, 0x1a, 0x8d, 0x25, 0x31, 0x76, 0x91, 0x20, 0x46, 0xcc, 0xa6, 0x90, 0xbd, 0x30, 0x9b,
	0x34, 0x43, 0x67, 0x28, 0x93, 0x6c, 0x67, 0x9a, 0x99, 0x29, 0x86, 0x10, 0xae, 0xbd, 0x50, 0xdf,
	0xc1, 0x78, 0xe9, 0xa3, 0xf8, 0x28, 0x3c, 0x85, 0xe9, 0xb4, 0x9d, 0x12, 0xb0, 0xcb, 0xde, 0x9d,
	0x9e, 0xf3, 0x7d, 0xdf, 0x39, 0xe7, 0x3b, 0xbb, 0x03, 0x76, 0x53, 0x21, 0xd2, 0x31, 0x45, 0x98,
	0x28, 0x54, 0x85, 0x65, 0x74, 0xe6, 0x23, 0x49, 0x95, 0x28, 0x64, 0x42, 0x15, 0xc2, 0x24, 0x4e,
	0xa5, 0x28, 0xf2, 0x38, 0x91, 0x4c, 0x53, 0xc9, 0x04, 0x8f, 0x15, 0xcb, 0x8a, 0x31, 0xd6, 0x4c,
	0x70, 0x98, 0x4b, 0xa1, 0x85, 0xbb, 0x5e, 0x71, 0x21, 0x26, 0x0a, 0x5a, 0x19, 0x78, 0xe6, 0x43,
	0x2b, 0xb3, 0x82, 0xba, 0x3a, 0x25, 0x22, 0xcb, 0x04, 0x47, 0x37, 0x35, 0x57, 0xfa, 0x5d, 0x04,
	0xca, 0x8b, 0x4c, 0x5d, 0xc3, 0xc7, 0x99, 0x20, 0xec, 0x84, 0x25, 0xf5, 0x07, 0xd5, 0xa7, 0x82,
	0xd4, 0x1a, 0xfe, 0xd4, 0x1a, 0xfa, 0x3c, 0xa7, 0x35, 0xe9, 0x51, 0x43, 0xca, 0x99, 0xb5, 0xa1,
	0x2e, 0x3d, 0xa9, 0x4b, 0xe6, 0xeb, 0xb8, 0x38, 0x41, 0x5f, 0x24, 0xce, 0x73, 0x2a, 0x55, 0x5d,
	0x5f, 0xbd, 0x46, 0xc5, 0x9c, 0x0b, 0x6d, 0xc4, 0xeb, 0xea, 0xb3, 0x9f, 0x73, 0x60, 0x25, 0x24,
	0x7b, 0xa5, 0x99, 0x3b, 0x8d, 0x97, 0x87, 0x76, 0x04, 0xf7, 0x39, 0xb8, 0xd7, 0xb4, 0x8b, 0x39,
	0xce, 0xa8, 0xe7, 0xac, 0x39, 0x9b, 0x8b, 0xd1, 0x72, 0x93, 0xfc, 0x84, 0x33, 0xea, 0x6e, 0x83,
	0x25, 0x7b, 0x10, 0x46, 0xbc, 0xde, 0x9a, 0xb3, 0xb9, 0xb4, 0xf5, 0xb8, 0x36, 0x1d, 0x36, 0x73,
	0xc1, 0x7d, 0xae, 0x5f, 0xbe, 0x18, 0xe2, 0x71, 0x41, 0xa3, 0x45, 0x5c, 0xb5, 0xdc, 0x27, 0xee,
	0x1b, 0xb0, 0xdc, 0x1e, 0x91, 0x11, 0x6f, 0xe6, 0x6e, 0xf6, 0x92, 0x25, 0xec, 0x13, 0xf7, 0x08,
	0xcc, 0x96, 0x3e, 0x79, 0xb3, 0x6b, 0xce, 0xe6, 0xfd, 0xad, 0xb7, 0xb0, 0xeb, 0xea, 0xc6, 0x5d,
	0xd8, 0xae, 0x76, 0x74, 0x9e, 0xd3, 0x5d, 0x5e, 0x64, 0x37, 0x52, 0x91, 0x51, 0x73, 0x7f, 0x38,
	0xe0, 0xe1, 0x3f, 0x4e, 0xe8, 0xfd, 0x6f, 0xba, 0x8c, 0xa6, 0xee, 0x72, 0x70, 0x4d, 0xe3, 0xc0,
	0x48, 0xdc, 0xe8, 0x79, 0x1b, 0x10, 0xb9, 0xd9, 0xad, 0x9c, 0xbb, 0x0d, 0x80, 0xd2, 0x58, 0xea,
	0x98, 0x60, 0x4d, 0xbd, 0x39, 0xe3, 0xd1, 0xea, 0x2d, 0x8f, 0x0e, 0xb5, 0x64, 0x3c, 0xad, 0x2d,
	0x36, 0xf8, 0x77, 0x58, 0x53, 0xf7, 0x15, 0x58, 0xa0, 0x9c, 0x54, 0xd4, 0xf9, 0x29, 0xa8, 0xf3,
	0x94, 0x13, 0x43, 0x3c, 0x05, 0x6e, 0x92, 0x27, 0xf1, 0x31, 0x23, 0x71, 0x2e, 0x18, 0xd7, 0xf1,
	0x98, 0x29, 0xed, 0x2d, 0x18, 0x89, 0xd7, 0x9d, 0x1e, 0x54, 0x7f, 0x1e, 0xb8, 0x93, 0x27, 0x7d,
	0x46, 0xda, 0x4d, 0x07, 0xa5, 0xc2, 0x47, 0xa6, 0xf4, 0xfb, 0xff, 0xa2, 0x07, 0x89, 0x29, 0xda,
	0x54, 0xf0, 0xcd, 0xb9, 0x0a, 0xbf, 0x3a, 0xc0, 0x6f, 0x85, 0xea, 0x28, 0x67, 0xaa, 0x14, 0x44,
	0x13, 0x7e, 0xa3, 0x83, 0xa4, 0x50, 0x5a, 0x64, 0x54, 0x2a, 0x74, 0xd1, 0x84, 0x97, 0x08, 0x77,
	0x12, 0x14, 0xba, 0x98, 0xf8, 0x7e, 0x5c, 0xf6, 0x97, 0x01, 0x68, 0xf7, 0xed, 0x7f, 0xef, 0x81,
	0x8d, 0x44, 0x64, 0xf0, 0xce, 0xf7, 0xa4, 0xff, 0xb4, 0x7b, 0xca, 0x41, 0x69, 0xf4, 0xc0, 0xf9,
	0xfc, 0xa1, 0x56, 0x49, 0xc5, 0x18, 0xf3, 0x14, 0x0a, 0x99, 0xa2, 0x94, 0x72, 0x73, 0x06, 0xd4,
	0xae, 0x3c, 0xe1, 0xed, 0xdb, 0xb6, 0xd1, 0xaf, 0xde, 0xcc, 0x5e, 0x18, 0xfe, 0xee, 0xad, 0xef,
	0x55, 0x92, 0x21, 0x51, 0xb0, 0x0a, 0xcb, 0x68, 0xe8, 0xc3, 0xa8, 0x41, 0xfe, 0x69, 0x30, 0xa3,
	0x90, 0xa8, 0x91, 0xc5, 0x8c, 0x86, 0xfe, 0xc8, 0x62, 0xae, 0x7a, 0x1b, 0x55, 0x21, 0x08, 0x42,
	0xa2, 0x82, 0xc0, 0xa2, 0x82, 0x60, 0xe8, 0x07, 0x81, 0xc5, 0x1d, 0xcf, 0x99, 0x61, 0xfd, 0xbf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x16, 0x8a, 0x0a, 0x2e, 0xa7, 0x05, 0x00, 0x00,
}

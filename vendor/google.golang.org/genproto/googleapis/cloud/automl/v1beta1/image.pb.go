// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/automl/v1beta1/image.proto

package automl

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

// Dataset metadata that is specific to image classification.
type ImageClassificationDatasetMetadata struct {
	// Required. Type of the classification problem.
	ClassificationType   ClassificationType `protobuf:"varint,1,opt,name=classification_type,json=classificationType,proto3,enum=google.cloud.automl.v1beta1.ClassificationType" json:"classification_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ImageClassificationDatasetMetadata) Reset()         { *m = ImageClassificationDatasetMetadata{} }
func (m *ImageClassificationDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*ImageClassificationDatasetMetadata) ProtoMessage()    {}
func (*ImageClassificationDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{0}
}

func (m *ImageClassificationDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageClassificationDatasetMetadata.Unmarshal(m, b)
}
func (m *ImageClassificationDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageClassificationDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *ImageClassificationDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageClassificationDatasetMetadata.Merge(m, src)
}
func (m *ImageClassificationDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageClassificationDatasetMetadata.Size(m)
}
func (m *ImageClassificationDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageClassificationDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageClassificationDatasetMetadata proto.InternalMessageInfo

func (m *ImageClassificationDatasetMetadata) GetClassificationType() ClassificationType {
	if m != nil {
		return m.ClassificationType
	}
	return ClassificationType_CLASSIFICATION_TYPE_UNSPECIFIED
}

// Dataset metadata specific to image object detection.
type ImageObjectDetectionDatasetMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageObjectDetectionDatasetMetadata) Reset()         { *m = ImageObjectDetectionDatasetMetadata{} }
func (m *ImageObjectDetectionDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*ImageObjectDetectionDatasetMetadata) ProtoMessage()    {}
func (*ImageObjectDetectionDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{1}
}

func (m *ImageObjectDetectionDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageObjectDetectionDatasetMetadata.Unmarshal(m, b)
}
func (m *ImageObjectDetectionDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageObjectDetectionDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *ImageObjectDetectionDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageObjectDetectionDatasetMetadata.Merge(m, src)
}
func (m *ImageObjectDetectionDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageObjectDetectionDatasetMetadata.Size(m)
}
func (m *ImageObjectDetectionDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageObjectDetectionDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageObjectDetectionDatasetMetadata proto.InternalMessageInfo

// Model metadata for image classification.
type ImageClassificationModelMetadata struct {
	// Optional. The ID of the `base` model. If it is specified, the new model
	// will be created based on the `base` model. Otherwise, the new model will be
	// created from scratch. The `base` model must be in the same
	// `project` and `location` as the new model to create, and have the same
	// `model_type`.
	BaseModelId string `protobuf:"bytes,1,opt,name=base_model_id,json=baseModelId,proto3" json:"base_model_id,omitempty"`
	// Required. The train budget of creating this model, expressed in hours. The
	// actual `train_cost` will be equal or less than this value.
	TrainBudget int64 `protobuf:"varint,2,opt,name=train_budget,json=trainBudget,proto3" json:"train_budget,omitempty"`
	// Output only. The actual train cost of creating this model, expressed in
	// hours. If this model is created from a `base` model, the train cost used
	// to create the `base` model are not included.
	TrainCost int64 `protobuf:"varint,3,opt,name=train_cost,json=trainCost,proto3" json:"train_cost,omitempty"`
	// Output only. The reason that this create model operation stopped,
	// e.g. `BUDGET_REACHED`, `MODEL_CONVERGED`.
	StopReason string `protobuf:"bytes,5,opt,name=stop_reason,json=stopReason,proto3" json:"stop_reason,omitempty"`
	// Optional. Type of the model. The available values are:
	// *   `cloud` - Model to be used via prediction calls to AutoML API.
	//               This is the default value.
	// *   `mobile-low-latency-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards. Expected to have low latency, but
	//               may have lower prediction quality than other models.
	// *   `mobile-versatile-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards.
	// *   `mobile-high-accuracy-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards.  Expected to have a higher
	//               latency, but should also have a higher prediction quality
	//               than other models.
	// *   `mobile-core-ml-low-latency-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile device with Core
	//               ML afterwards. Expected to have low latency, but may have
	//               lower prediction quality than other models.
	// *   `mobile-core-ml-versatile-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile device with Core
	//               ML afterwards.
	// *   `mobile-core-ml-high-accuracy-1` - A model that, in addition to
	//               providing prediction via AutoML API, can also be exported
	//               (see [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile device with
	//               Core ML afterwards.  Expected to have a higher latency, but
	//               should also have a higher prediction quality than other
	//               models.
	ModelType string `protobuf:"bytes,7,opt,name=model_type,json=modelType,proto3" json:"model_type,omitempty"`
	// Output only. An approximate number of online prediction QPS that can
	// be supported by this model per each node on which it is deployed.
	NodeQps float64 `protobuf:"fixed64,13,opt,name=node_qps,json=nodeQps,proto3" json:"node_qps,omitempty"`
	// Output only. The number of nodes this model is deployed on. A node is an
	// abstraction of a machine resource, which can handle online prediction QPS
	// as given in the node_qps field.
	NodeCount            int64    `protobuf:"varint,14,opt,name=node_count,json=nodeCount,proto3" json:"node_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageClassificationModelMetadata) Reset()         { *m = ImageClassificationModelMetadata{} }
func (m *ImageClassificationModelMetadata) String() string { return proto.CompactTextString(m) }
func (*ImageClassificationModelMetadata) ProtoMessage()    {}
func (*ImageClassificationModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{2}
}

func (m *ImageClassificationModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageClassificationModelMetadata.Unmarshal(m, b)
}
func (m *ImageClassificationModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageClassificationModelMetadata.Marshal(b, m, deterministic)
}
func (m *ImageClassificationModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageClassificationModelMetadata.Merge(m, src)
}
func (m *ImageClassificationModelMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageClassificationModelMetadata.Size(m)
}
func (m *ImageClassificationModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageClassificationModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageClassificationModelMetadata proto.InternalMessageInfo

func (m *ImageClassificationModelMetadata) GetBaseModelId() string {
	if m != nil {
		return m.BaseModelId
	}
	return ""
}

func (m *ImageClassificationModelMetadata) GetTrainBudget() int64 {
	if m != nil {
		return m.TrainBudget
	}
	return 0
}

func (m *ImageClassificationModelMetadata) GetTrainCost() int64 {
	if m != nil {
		return m.TrainCost
	}
	return 0
}

func (m *ImageClassificationModelMetadata) GetStopReason() string {
	if m != nil {
		return m.StopReason
	}
	return ""
}

func (m *ImageClassificationModelMetadata) GetModelType() string {
	if m != nil {
		return m.ModelType
	}
	return ""
}

func (m *ImageClassificationModelMetadata) GetNodeQps() float64 {
	if m != nil {
		return m.NodeQps
	}
	return 0
}

func (m *ImageClassificationModelMetadata) GetNodeCount() int64 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

// Model metadata specific to image object detection.
type ImageObjectDetectionModelMetadata struct {
	// Optional. Type of the model. The available values are:
	// *   `cloud-high-accuracy-1` - (default) A model to be used via prediction
	//               calls to AutoML API. Expected to have a higher latency, but
	//               should also have a higher prediction quality than other
	//               models.
	// *   `cloud-low-latency-1` -  A model to be used via prediction
	//               calls to AutoML API. Expected to have low latency, but may
	//               have lower prediction quality than other models.
	// *   `mobile-low-latency-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards. Expected to have low latency, but
	//               may have lower prediction quality than other models.
	// *   `mobile-versatile-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards.
	// *   `mobile-high-accuracy-1` - A model that, in addition to providing
	//               prediction via AutoML API, can also be exported (see
	//               [AutoMl.ExportModel][google.cloud.automl.v1beta1.AutoMl.ExportModel]) and used on a mobile or edge device
	//               with TensorFlow afterwards.  Expected to have a higher
	//               latency, but should also have a higher prediction quality
	//               than other models.
	ModelType string `protobuf:"bytes,1,opt,name=model_type,json=modelType,proto3" json:"model_type,omitempty"`
	// Output only. The number of nodes this model is deployed on. A node is an
	// abstraction of a machine resource, which can handle online prediction QPS
	// as given in the qps_per_node field.
	NodeCount int64 `protobuf:"varint,3,opt,name=node_count,json=nodeCount,proto3" json:"node_count,omitempty"`
	// Output only. An approximate number of online prediction QPS that can
	// be supported by this model per each node on which it is deployed.
	NodeQps float64 `protobuf:"fixed64,4,opt,name=node_qps,json=nodeQps,proto3" json:"node_qps,omitempty"`
	// Output only. The reason that this create model operation stopped,
	// e.g. `BUDGET_REACHED`, `MODEL_CONVERGED`.
	StopReason string `protobuf:"bytes,5,opt,name=stop_reason,json=stopReason,proto3" json:"stop_reason,omitempty"`
	// The train budget of creating this model, expressed in milli node
	// hours i.e. 1,000 value in this field means 1 node hour. The actual
	// `train_cost` will be equal or less than this value. If further model
	// training ceases to provide any improvements, it will stop without using
	// full budget and the stop_reason will be `MODEL_CONVERGED`.
	// Note, node_hour  = actual_hour * number_of_nodes_invovled.
	// For model type `cloud-high-accuracy-1`(default) and `cloud-low-latency-1`,
	// the train budget must be between 20,000 and 900,000 milli node hours,
	// inclusive. The default value is 216, 000 which represents one day in
	// wall time.
	// For model type `mobile-low-latency-1`, `mobile-versatile-1`,
	// `mobile-high-accuracy-1`, `mobile-core-ml-low-latency-1`,
	// `mobile-core-ml-versatile-1`, `mobile-core-ml-high-accuracy-1`, the train
	// budget must be between 1,000 and 100,000 milli node hours, inclusive.
	// The default value is 24, 000 which represents one day in wall time.
	TrainBudgetMilliNodeHours int64 `protobuf:"varint,6,opt,name=train_budget_milli_node_hours,json=trainBudgetMilliNodeHours,proto3" json:"train_budget_milli_node_hours,omitempty"`
	// Output only. The actual train cost of creating this model, expressed in
	// milli node hours, i.e. 1,000 value in this field means 1 node hour.
	// Guaranteed to not exceed the train budget.
	TrainCostMilliNodeHours int64    `protobuf:"varint,7,opt,name=train_cost_milli_node_hours,json=trainCostMilliNodeHours,proto3" json:"train_cost_milli_node_hours,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *ImageObjectDetectionModelMetadata) Reset()         { *m = ImageObjectDetectionModelMetadata{} }
func (m *ImageObjectDetectionModelMetadata) String() string { return proto.CompactTextString(m) }
func (*ImageObjectDetectionModelMetadata) ProtoMessage()    {}
func (*ImageObjectDetectionModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{3}
}

func (m *ImageObjectDetectionModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageObjectDetectionModelMetadata.Unmarshal(m, b)
}
func (m *ImageObjectDetectionModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageObjectDetectionModelMetadata.Marshal(b, m, deterministic)
}
func (m *ImageObjectDetectionModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageObjectDetectionModelMetadata.Merge(m, src)
}
func (m *ImageObjectDetectionModelMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageObjectDetectionModelMetadata.Size(m)
}
func (m *ImageObjectDetectionModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageObjectDetectionModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageObjectDetectionModelMetadata proto.InternalMessageInfo

func (m *ImageObjectDetectionModelMetadata) GetModelType() string {
	if m != nil {
		return m.ModelType
	}
	return ""
}

func (m *ImageObjectDetectionModelMetadata) GetNodeCount() int64 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

func (m *ImageObjectDetectionModelMetadata) GetNodeQps() float64 {
	if m != nil {
		return m.NodeQps
	}
	return 0
}

func (m *ImageObjectDetectionModelMetadata) GetStopReason() string {
	if m != nil {
		return m.StopReason
	}
	return ""
}

func (m *ImageObjectDetectionModelMetadata) GetTrainBudgetMilliNodeHours() int64 {
	if m != nil {
		return m.TrainBudgetMilliNodeHours
	}
	return 0
}

func (m *ImageObjectDetectionModelMetadata) GetTrainCostMilliNodeHours() int64 {
	if m != nil {
		return m.TrainCostMilliNodeHours
	}
	return 0
}

// Model deployment metadata specific to Image Classification.
type ImageClassificationModelDeploymentMetadata struct {
	// Input only. The number of nodes to deploy the model on. A node is an
	// abstraction of a machine resource, which can handle online prediction QPS
	// as given in the model's
	//
	// [node_qps][google.cloud.automl.v1beta1.ImageClassificationModelMetadata.node_qps].
	// Must be between 1 and 100, inclusive on both ends.
	NodeCount            int64    `protobuf:"varint,1,opt,name=node_count,json=nodeCount,proto3" json:"node_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageClassificationModelDeploymentMetadata) Reset() {
	*m = ImageClassificationModelDeploymentMetadata{}
}
func (m *ImageClassificationModelDeploymentMetadata) String() string {
	return proto.CompactTextString(m)
}
func (*ImageClassificationModelDeploymentMetadata) ProtoMessage() {}
func (*ImageClassificationModelDeploymentMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{4}
}

func (m *ImageClassificationModelDeploymentMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageClassificationModelDeploymentMetadata.Unmarshal(m, b)
}
func (m *ImageClassificationModelDeploymentMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageClassificationModelDeploymentMetadata.Marshal(b, m, deterministic)
}
func (m *ImageClassificationModelDeploymentMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageClassificationModelDeploymentMetadata.Merge(m, src)
}
func (m *ImageClassificationModelDeploymentMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageClassificationModelDeploymentMetadata.Size(m)
}
func (m *ImageClassificationModelDeploymentMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageClassificationModelDeploymentMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageClassificationModelDeploymentMetadata proto.InternalMessageInfo

func (m *ImageClassificationModelDeploymentMetadata) GetNodeCount() int64 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

// Model deployment metadata specific to Image Object Detection.
type ImageObjectDetectionModelDeploymentMetadata struct {
	// Input only. The number of nodes to deploy the model on. A node is an
	// abstraction of a machine resource, which can handle online prediction QPS
	// as given in the model's
	//
	// [qps_per_node][google.cloud.automl.v1beta1.ImageObjectDetectionModelMetadata.qps_per_node].
	// Must be between 1 and 100, inclusive on both ends.
	NodeCount            int64    `protobuf:"varint,1,opt,name=node_count,json=nodeCount,proto3" json:"node_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageObjectDetectionModelDeploymentMetadata) Reset() {
	*m = ImageObjectDetectionModelDeploymentMetadata{}
}
func (m *ImageObjectDetectionModelDeploymentMetadata) String() string {
	return proto.CompactTextString(m)
}
func (*ImageObjectDetectionModelDeploymentMetadata) ProtoMessage() {}
func (*ImageObjectDetectionModelDeploymentMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_29b9f2bc900da869, []int{5}
}

func (m *ImageObjectDetectionModelDeploymentMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata.Unmarshal(m, b)
}
func (m *ImageObjectDetectionModelDeploymentMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata.Marshal(b, m, deterministic)
}
func (m *ImageObjectDetectionModelDeploymentMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata.Merge(m, src)
}
func (m *ImageObjectDetectionModelDeploymentMetadata) XXX_Size() int {
	return xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata.Size(m)
}
func (m *ImageObjectDetectionModelDeploymentMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ImageObjectDetectionModelDeploymentMetadata proto.InternalMessageInfo

func (m *ImageObjectDetectionModelDeploymentMetadata) GetNodeCount() int64 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

func init() {
	proto.RegisterType((*ImageClassificationDatasetMetadata)(nil), "google.cloud.automl.v1beta1.ImageClassificationDatasetMetadata")
	proto.RegisterType((*ImageObjectDetectionDatasetMetadata)(nil), "google.cloud.automl.v1beta1.ImageObjectDetectionDatasetMetadata")
	proto.RegisterType((*ImageClassificationModelMetadata)(nil), "google.cloud.automl.v1beta1.ImageClassificationModelMetadata")
	proto.RegisterType((*ImageObjectDetectionModelMetadata)(nil), "google.cloud.automl.v1beta1.ImageObjectDetectionModelMetadata")
	proto.RegisterType((*ImageClassificationModelDeploymentMetadata)(nil), "google.cloud.automl.v1beta1.ImageClassificationModelDeploymentMetadata")
	proto.RegisterType((*ImageObjectDetectionModelDeploymentMetadata)(nil), "google.cloud.automl.v1beta1.ImageObjectDetectionModelDeploymentMetadata")
}

func init() {
	proto.RegisterFile("google/cloud/automl/v1beta1/image.proto", fileDescriptor_29b9f2bc900da869)
}

var fileDescriptor_29b9f2bc900da869 = []byte{
	// 574 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x14, 0x94, 0x13, 0x68, 0xc8, 0x86, 0xf6, 0x60, 0x0e, 0x38, 0x0d, 0x25, 0xa9, 0x11, 0x22, 0x02,
	0xc9, 0x26, 0x70, 0x33, 0x1c, 0x68, 0x12, 0x09, 0x2a, 0x1a, 0x28, 0x11, 0xe2, 0x80, 0x22, 0x99,
	0xb5, 0xbd, 0x35, 0x46, 0xb6, 0x9f, 0xf1, 0x3e, 0x23, 0xe5, 0xca, 0x81, 0xbf, 0xe0, 0x23, 0xf8,
	0x0d, 0x3e, 0x85, 0xaf, 0x40, 0xfb, 0x6c, 0xb5, 0x71, 0x9a, 0x06, 0x71, 0xdc, 0x79, 0x33, 0x6f,
	0x66, 0x27, 0x59, 0xb3, 0x07, 0x21, 0x40, 0x18, 0x0b, 0xdb, 0x8f, 0xa1, 0x08, 0x6c, 0x5e, 0x20,
	0x24, 0xb1, 0xfd, 0x6d, 0xe4, 0x09, 0xe4, 0x23, 0x3b, 0x4a, 0x78, 0x28, 0xac, 0x2c, 0x07, 0x04,
	0xbd, 0x57, 0x12, 0x2d, 0x22, 0x5a, 0x25, 0xd1, 0xaa, 0x88, 0xfb, 0xdd, 0x6a, 0x0b, 0xcf, 0x22,
	0x3b, 0x17, 0x12, 0x8a, 0xdc, 0xaf, 0x74, 0xfb, 0xa3, 0x6d, 0x06, 0x3c, 0x4d, 0x01, 0x39, 0x46,
	0x90, 0xba, 0x32, 0x13, 0x7e, 0x25, 0x79, 0xbc, 0x4d, 0xe2, 0xc7, 0x5c, 0xca, 0xe8, 0x2c, 0xf2,
	0x49, 0x56, 0x29, 0xfa, 0x95, 0x82, 0x4e, 0x5e, 0x71, 0x66, 0x63, 0x94, 0x08, 0x89, 0x3c, 0xc9,
	0x2a, 0xc2, 0x9d, 0x95, 0x80, 0x17, 0xa6, 0xb2, 0x9c, 0x9a, 0x3f, 0x34, 0x66, 0x1e, 0xab, 0xbb,
	0x4e, 0x6a, 0xcb, 0xa7, 0x1c, 0xb9, 0x14, 0x38, 0x13, 0xc8, 0x03, 0x8e, 0x5c, 0xff, 0xc4, 0x6e,
	0xd5, 0xdd, 0x5d, 0x5c, 0x66, 0xc2, 0xd0, 0x06, 0xda, 0x70, 0xef, 0x89, 0x6d, 0x6d, 0x29, 0xc8,
	0xaa, 0x2f, 0x7e, 0xbf, 0xcc, 0xc4, 0x5c, 0xf7, 0x2f, 0x61, 0xe6, 0x7d, 0x76, 0x8f, 0x72, 0xbc,
	0xf5, 0xbe, 0x08, 0x1f, 0xa7, 0x02, 0x85, 0xbf, 0x21, 0x88, 0xf9, 0xbd, 0xc1, 0x06, 0x1b, 0xf2,
	0xce, 0x20, 0x10, 0xf1, 0x79, 0x5a, 0x93, 0xed, 0x7a, 0x5c, 0x0a, 0x37, 0x51, 0xa8, 0x1b, 0x05,
	0x94, 0xb3, 0x3d, 0xef, 0x28, 0x90, 0x98, 0xc7, 0x81, 0x7e, 0xc8, 0x6e, 0x62, 0xce, 0xa3, 0xd4,
	0xf5, 0x8a, 0x20, 0x14, 0x68, 0x34, 0x06, 0xda, 0xb0, 0x39, 0xef, 0x10, 0x36, 0x26, 0x48, 0x3f,
	0x60, 0xac, 0xa4, 0xf8, 0x20, 0xd1, 0x68, 0x12, 0xa1, 0x4d, 0xc8, 0x04, 0x24, 0xea, 0x7d, 0xd6,
	0x91, 0x08, 0x99, 0x9b, 0x0b, 0x2e, 0x21, 0x35, 0xae, 0x93, 0x07, 0x53, 0xd0, 0x9c, 0x10, 0xa5,
	0x2f, 0x13, 0x50, 0x57, 0x2d, 0x9a, 0xb7, 0x09, 0x51, 0x37, 0xd6, 0xbb, 0xec, 0x46, 0x0a, 0x81,
	0x70, 0xbf, 0x66, 0xd2, 0xd8, 0x1d, 0x68, 0x43, 0x6d, 0xde, 0x52, 0xe7, 0x77, 0x99, 0x54, 0x4a,
	0x1a, 0xf9, 0x50, 0xa4, 0x68, 0xec, 0x95, 0xce, 0x0a, 0x99, 0x28, 0xc0, 0xfc, 0xd9, 0x60, 0x87,
	0x9b, 0xca, 0xaa, 0xb7, 0x50, 0xb7, 0xd7, 0xd6, 0xed, 0xeb, 0x1e, 0xcd, 0x35, 0x8f, 0x5a, 0xba,
	0x6b, 0xf5, 0x74, 0xff, 0xbc, 0xf8, 0x0b, 0x76, 0xb0, 0xda, 0xad, 0x9b, 0x44, 0x71, 0x1c, 0xb9,
	0xb4, 0xee, 0x33, 0x14, 0xb9, 0x34, 0x76, 0xc8, 0xad, 0xbb, 0x52, 0xf6, 0x4c, 0x51, 0xde, 0x40,
	0x20, 0x5e, 0x29, 0x82, 0xfe, 0x9c, 0xf5, 0x2e, 0xaa, 0xbf, 0xac, 0x6f, 0x91, 0xfe, 0xf6, 0xf9,
	0x6f, 0x51, 0x57, 0x9b, 0xaf, 0xd9, 0xc3, 0xab, 0xfe, 0x23, 0x53, 0x91, 0xc5, 0xb0, 0x4c, 0x44,
	0x8a, 0xab, 0x3d, 0xad, 0x14, 0xa1, 0xad, 0x97, 0x7d, 0xc2, 0x1e, 0x5d, 0xd9, 0xf5, 0x7f, 0x6f,
	0x1b, 0xff, 0xd2, 0x58, 0xdf, 0x87, 0x64, 0xdb, 0x8b, 0x19, 0x33, 0xf2, 0x3b, 0x55, 0xef, 0xf3,
	0x54, 0xfb, 0x78, 0x54, 0x51, 0x43, 0x88, 0x79, 0x1a, 0x5a, 0x90, 0x87, 0x76, 0x28, 0x52, 0x7a,
	0xbd, 0x76, 0x39, 0xe2, 0x59, 0x24, 0x37, 0x7e, 0x31, 0x9e, 0x95, 0xc7, 0xdf, 0x8d, 0xde, 0x4b,
	0x22, 0x2e, 0x26, 0x8a, 0xb4, 0x38, 0x2a, 0x10, 0x66, 0xf1, 0xe2, 0x43, 0x49, 0xfa, 0xd3, 0xb8,
	0x5b, 0x4e, 0x1d, 0x87, 0xc6, 0x8e, 0x43, 0xf3, 0x13, 0xc7, 0xa9, 0x08, 0xde, 0x0e, 0x99, 0x3d,
	0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xd3, 0xc2, 0x5d, 0x90, 0x31, 0x05, 0x00, 0x00,
}

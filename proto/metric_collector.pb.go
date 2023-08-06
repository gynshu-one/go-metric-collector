// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: metric_collector.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	MType string  `protobuf:"bytes,2,opt,name=MType,proto3" json:"MType,omitempty"`
	Value float64 `protobuf:"fixed64,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Delta int64   `protobuf:"varint,4,opt,name=Delta,proto3" json:"Delta,omitempty"`
	Hash  string  `protobuf:"bytes,5,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Metric) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type LiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LiveRequest) Reset() {
	*x = LiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LiveRequest) ProtoMessage() {}

func (x *LiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LiveRequest.ProtoReflect.Descriptor instead.
func (*LiveRequest) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{1}
}

type LiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *LiveResponse) Reset() {
	*x = LiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LiveResponse) ProtoMessage() {}

func (x *LiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LiveResponse.ProtoReflect.Descriptor instead.
func (*LiveResponse) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{2}
}

func (x *LiveResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ValueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetricName string `protobuf:"bytes,1,opt,name=metric_name,json=metricName,proto3" json:"metric_name,omitempty"`
	MetricType string `protobuf:"bytes,2,opt,name=metric_type,json=metricType,proto3" json:"metric_type,omitempty"`
}

func (x *ValueRequest) Reset() {
	*x = ValueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueRequest) ProtoMessage() {}

func (x *ValueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueRequest.ProtoReflect.Descriptor instead.
func (*ValueRequest) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{3}
}

func (x *ValueRequest) GetMetricName() string {
	if x != nil {
		return x.MetricName
	}
	return ""
}

func (x *ValueRequest) GetMetricType() string {
	if x != nil {
		return x.MetricType
	}
	return ""
}

type UpdateMetricsJSONRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateMetricsJSONRequest) Reset() {
	*x = UpdateMetricsJSONRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsJSONRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsJSONRequest) ProtoMessage() {}

func (x *UpdateMetricsJSONRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsJSONRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsJSONRequest) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateMetricsJSONRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetricName  string `protobuf:"bytes,1,opt,name=metric_name,json=metricName,proto3" json:"metric_name,omitempty"`
	MetricType  string `protobuf:"bytes,2,opt,name=metric_type,json=metricType,proto3" json:"metric_type,omitempty"`
	MetricValue string `protobuf:"bytes,3,opt,name=metric_value,json=metricValue,proto3" json:"metric_value,omitempty"`
}

func (x *UpdateMetricRequest) Reset() {
	*x = UpdateMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricRequest) ProtoMessage() {}

func (x *UpdateMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricRequest) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMetricRequest) GetMetricName() string {
	if x != nil {
		return x.MetricName
	}
	return ""
}

func (x *UpdateMetricRequest) GetMetricType() string {
	if x != nil {
		return x.MetricType
	}
	return ""
}

func (x *UpdateMetricRequest) GetMetricValue() string {
	if x != nil {
		return x.MetricValue
	}
	return ""
}

type BulkUpdateJSONRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *BulkUpdateJSONRequest) Reset() {
	*x = BulkUpdateJSONRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BulkUpdateJSONRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BulkUpdateJSONRequest) ProtoMessage() {}

func (x *BulkUpdateJSONRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BulkUpdateJSONRequest.ProtoReflect.Descriptor instead.
func (*BulkUpdateJSONRequest) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{6}
}

func (x *BulkUpdateJSONRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type ValueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ValueResponse) Reset() {
	*x = ValueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueResponse) ProtoMessage() {}

func (x *ValueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueResponse.ProtoReflect.Descriptor instead.
func (*ValueResponse) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{7}
}

func (x *ValueResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type MetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *MetricResponse) Reset() {
	*x = MetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricResponse) ProtoMessage() {}

func (x *MetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricResponse.ProtoReflect.Descriptor instead.
func (*MetricResponse) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{8}
}

func (x *MetricResponse) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type BulkUpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *BulkUpdateResponse) Reset() {
	*x = BulkUpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BulkUpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BulkUpdateResponse) ProtoMessage() {}

func (x *BulkUpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BulkUpdateResponse.ProtoReflect.Descriptor instead.
func (*BulkUpdateResponse) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{9}
}

func (x *BulkUpdateResponse) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type PingDBResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PingDBResponse) Reset() {
	*x = PingDBResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metric_collector_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingDBResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingDBResponse) ProtoMessage() {}

func (x *PingDBResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metric_collector_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingDBResponse.ProtoReflect.Descriptor instead.
func (*PingDBResponse) Descriptor() ([]byte, []int) {
	return file_metric_collector_proto_rawDescGZIP(), []int{10}
}

func (x *PingDBResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_metric_collector_proto protoreflect.FileDescriptor

var file_metric_collector_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6e, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x14, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x44,
	0x65, 0x6c, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74,
	0x61, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x22, 0x0d, 0x0a, 0x0b, 0x4c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x4c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x50,
	0x0a, 0x0c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x3b, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x4a, 0x53, 0x4f, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x7a, 0x0a,
	0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x0a, 0x15, 0x42, 0x75, 0x6c,
	0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x25, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x31, 0x0a, 0x0e,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22,
	0x37, 0x0a, 0x12, 0x42, 0x75, 0x6c, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x2a, 0x0a, 0x0e, 0x50, 0x69, 0x6e, 0x67,
	0x44, 0x42, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0xfd, 0x02, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x4c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4a, 0x53,
	0x4f, 0x4e, 0x12, 0x0d, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0d, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x11, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x4a, 0x53, 0x4f, 0x4e, 0x12,
	0x19, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x4a,
	0x53, 0x4f, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x42, 0x75, 0x6c, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4a, 0x53, 0x4f, 0x4e, 0x12, 0x16, 0x2e, 0x42, 0x75, 0x6c, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4a, 0x53, 0x4f, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x42,
	0x75, 0x6c, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x31, 0x0a, 0x06, 0x50, 0x69, 0x6e, 0x67, 0x44, 0x42, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x44, 0x42, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x79, 0x6e, 0x73, 0x68, 0x75, 0x2d, 0x6f, 0x6e, 0x65, 0x2f, 0x67, 0x6f,
	0x2d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2d, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metric_collector_proto_rawDescOnce sync.Once
	file_metric_collector_proto_rawDescData = file_metric_collector_proto_rawDesc
)

func file_metric_collector_proto_rawDescGZIP() []byte {
	file_metric_collector_proto_rawDescOnce.Do(func() {
		file_metric_collector_proto_rawDescData = protoimpl.X.CompressGZIP(file_metric_collector_proto_rawDescData)
	})
	return file_metric_collector_proto_rawDescData
}

var file_metric_collector_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_metric_collector_proto_goTypes = []interface{}{
	(*Metric)(nil),                   // 0: Metric
	(*LiveRequest)(nil),              // 1: LiveRequest
	(*LiveResponse)(nil),             // 2: LiveResponse
	(*ValueRequest)(nil),             // 3: ValueRequest
	(*UpdateMetricsJSONRequest)(nil), // 4: UpdateMetricsJSONRequest
	(*UpdateMetricRequest)(nil),      // 5: UpdateMetricRequest
	(*BulkUpdateJSONRequest)(nil),    // 6: BulkUpdateJSONRequest
	(*ValueResponse)(nil),            // 7: ValueResponse
	(*MetricResponse)(nil),           // 8: MetricResponse
	(*BulkUpdateResponse)(nil),       // 9: BulkUpdateResponse
	(*PingDBResponse)(nil),           // 10: PingDBResponse
	(*emptypb.Empty)(nil),            // 11: google.protobuf.Empty
}
var file_metric_collector_proto_depIdxs = []int32{
	0,  // 0: UpdateMetricsJSONRequest.metric:type_name -> Metric
	0,  // 1: BulkUpdateJSONRequest.metrics:type_name -> Metric
	0,  // 2: MetricResponse.metric:type_name -> Metric
	0,  // 3: BulkUpdateResponse.metrics:type_name -> Metric
	11, // 4: MetricService.Live:input_type -> google.protobuf.Empty
	3,  // 5: MetricService.ValueJSON:input_type -> ValueRequest
	3,  // 6: MetricService.Value:input_type -> ValueRequest
	4,  // 7: MetricService.UpdateMetricsJSON:input_type -> UpdateMetricsJSONRequest
	5,  // 8: MetricService.UpdateMetric:input_type -> UpdateMetricRequest
	6,  // 9: MetricService.BulkUpdateJSON:input_type -> BulkUpdateJSONRequest
	11, // 10: MetricService.PingDB:input_type -> google.protobuf.Empty
	2,  // 11: MetricService.Live:output_type -> LiveResponse
	8,  // 12: MetricService.ValueJSON:output_type -> MetricResponse
	7,  // 13: MetricService.Value:output_type -> ValueResponse
	8,  // 14: MetricService.UpdateMetricsJSON:output_type -> MetricResponse
	8,  // 15: MetricService.UpdateMetric:output_type -> MetricResponse
	9,  // 16: MetricService.BulkUpdateJSON:output_type -> BulkUpdateResponse
	10, // 17: MetricService.PingDB:output_type -> PingDBResponse
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_metric_collector_proto_init() }
func file_metric_collector_proto_init() {
	if File_metric_collector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metric_collector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LiveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LiveResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValueRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricsJSONRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMetricRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BulkUpdateJSONRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValueResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BulkUpdateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metric_collector_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingDBResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metric_collector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metric_collector_proto_goTypes,
		DependencyIndexes: file_metric_collector_proto_depIdxs,
		MessageInfos:      file_metric_collector_proto_msgTypes,
	}.Build()
	File_metric_collector_proto = out.File
	file_metric_collector_proto_rawDesc = nil
	file_metric_collector_proto_goTypes = nil
	file_metric_collector_proto_depIdxs = nil
}

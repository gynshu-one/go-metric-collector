// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: metric_collector.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MetricService_Live_FullMethodName              = "/MetricService/Live"
	MetricService_ValueJSON_FullMethodName         = "/MetricService/ValueJSON"
	MetricService_Value_FullMethodName             = "/MetricService/Value"
	MetricService_UpdateMetricsJSON_FullMethodName = "/MetricService/UpdateMetricsJSON"
	MetricService_UpdateMetric_FullMethodName      = "/MetricService/UpdateMetric"
	MetricService_BulkUpdateJSON_FullMethodName    = "/MetricService/BulkUpdateJSON"
	MetricService_PingDB_FullMethodName            = "/MetricService/PingDB"
)

// MetricServiceClient is the client API for MetricService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricServiceClient interface {
	Live(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*LiveResponse, error)
	ValueJSON(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*MetricResponse, error)
	Value(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error)
	UpdateMetricsJSON(ctx context.Context, in *UpdateMetricsJSONRequest, opts ...grpc.CallOption) (*MetricResponse, error)
	UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*MetricResponse, error)
	BulkUpdateJSON(ctx context.Context, in *BulkUpdateJSONRequest, opts ...grpc.CallOption) (*BulkUpdateResponse, error)
	PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingDBResponse, error)
}

type metricServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricServiceClient(cc grpc.ClientConnInterface) MetricServiceClient {
	return &metricServiceClient{cc}
}

func (c *metricServiceClient) Live(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*LiveResponse, error) {
	out := new(LiveResponse)
	err := c.cc.Invoke(ctx, MetricService_Live_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) ValueJSON(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*MetricResponse, error) {
	out := new(MetricResponse)
	err := c.cc.Invoke(ctx, MetricService_ValueJSON_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) Value(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, MetricService_Value_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) UpdateMetricsJSON(ctx context.Context, in *UpdateMetricsJSONRequest, opts ...grpc.CallOption) (*MetricResponse, error) {
	out := new(MetricResponse)
	err := c.cc.Invoke(ctx, MetricService_UpdateMetricsJSON_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*MetricResponse, error) {
	out := new(MetricResponse)
	err := c.cc.Invoke(ctx, MetricService_UpdateMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) BulkUpdateJSON(ctx context.Context, in *BulkUpdateJSONRequest, opts ...grpc.CallOption) (*BulkUpdateResponse, error) {
	out := new(BulkUpdateResponse)
	err := c.cc.Invoke(ctx, MetricService_BulkUpdateJSON_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingDBResponse, error) {
	out := new(PingDBResponse)
	err := c.cc.Invoke(ctx, MetricService_PingDB_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricServiceServer is the server API for MetricService service.
// All implementations must embed UnimplementedMetricServiceServer
// for forward compatibility
type MetricServiceServer interface {
	Live(context.Context, *emptypb.Empty) (*LiveResponse, error)
	ValueJSON(context.Context, *ValueRequest) (*MetricResponse, error)
	Value(context.Context, *ValueRequest) (*ValueResponse, error)
	UpdateMetricsJSON(context.Context, *UpdateMetricsJSONRequest) (*MetricResponse, error)
	UpdateMetric(context.Context, *UpdateMetricRequest) (*MetricResponse, error)
	BulkUpdateJSON(context.Context, *BulkUpdateJSONRequest) (*BulkUpdateResponse, error)
	PingDB(context.Context, *emptypb.Empty) (*PingDBResponse, error)
	mustEmbedUnimplementedMetricServiceServer()
}

// UnimplementedMetricServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetricServiceServer struct {
}

func (UnimplementedMetricServiceServer) Live(context.Context, *emptypb.Empty) (*LiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Live not implemented")
}
func (UnimplementedMetricServiceServer) ValueJSON(context.Context, *ValueRequest) (*MetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValueJSON not implemented")
}
func (UnimplementedMetricServiceServer) Value(context.Context, *ValueRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Value not implemented")
}
func (UnimplementedMetricServiceServer) UpdateMetricsJSON(context.Context, *UpdateMetricsJSONRequest) (*MetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetricsJSON not implemented")
}
func (UnimplementedMetricServiceServer) UpdateMetric(context.Context, *UpdateMetricRequest) (*MetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMetricServiceServer) BulkUpdateJSON(context.Context, *BulkUpdateJSONRequest) (*BulkUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkUpdateJSON not implemented")
}
func (UnimplementedMetricServiceServer) PingDB(context.Context, *emptypb.Empty) (*PingDBResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingDB not implemented")
}
func (UnimplementedMetricServiceServer) mustEmbedUnimplementedMetricServiceServer() {}

// UnsafeMetricServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricServiceServer will
// result in compilation errors.
type UnsafeMetricServiceServer interface {
	mustEmbedUnimplementedMetricServiceServer()
}

func RegisterMetricServiceServer(s grpc.ServiceRegistrar, srv MetricServiceServer) {
	s.RegisterService(&MetricService_ServiceDesc, srv)
}

func _MetricService_Live_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).Live(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_Live_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).Live(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_ValueJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).ValueJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_ValueJSON_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).ValueJSON(ctx, req.(*ValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_Value_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).Value(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_Value_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).Value(ctx, req.(*ValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_UpdateMetricsJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricsJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).UpdateMetricsJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_UpdateMetricsJSON_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).UpdateMetricsJSON(ctx, req.(*UpdateMetricsJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).UpdateMetric(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_BulkUpdateJSON_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkUpdateJSONRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).BulkUpdateJSON(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_BulkUpdateJSON_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).BulkUpdateJSON(ctx, req.(*BulkUpdateJSONRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_PingDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).PingDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricService_PingDB_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).PingDB(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricService_ServiceDesc is the grpc.ServiceDesc for MetricService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MetricService",
	HandlerType: (*MetricServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Live",
			Handler:    _MetricService_Live_Handler,
		},
		{
			MethodName: "ValueJSON",
			Handler:    _MetricService_ValueJSON_Handler,
		},
		{
			MethodName: "Value",
			Handler:    _MetricService_Value_Handler,
		},
		{
			MethodName: "UpdateMetricsJSON",
			Handler:    _MetricService_UpdateMetricsJSON_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _MetricService_UpdateMetric_Handler,
		},
		{
			MethodName: "BulkUpdateJSON",
			Handler:    _MetricService_BulkUpdateJSON_Handler,
		},
		{
			MethodName: "PingDB",
			Handler:    _MetricService_PingDB_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metric_collector.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: unicom/api/v1/service.proto

package apiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UnicomClient is the client API for Unicom service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UnicomClient interface {
	SendSync(ctx context.Context, in *SendSyncRequest, opts ...grpc.CallOption) (*SendResponse, error)
	SendAsync(ctx context.Context, in *SendAsyncRequest, opts ...grpc.CallOption) (*SendResponse, error)
	GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error)
}

type unicomClient struct {
	cc grpc.ClientConnInterface
}

func NewUnicomClient(cc grpc.ClientConnInterface) UnicomClient {
	return &unicomClient{cc}
}

func (c *unicomClient) SendSync(ctx context.Context, in *SendSyncRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/unicom.api.v1.Unicom/SendSync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unicomClient) SendAsync(ctx context.Context, in *SendAsyncRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/unicom.api.v1.Unicom/SendAsync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unicomClient) GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error) {
	out := new(GetStatusResponse)
	err := c.cc.Invoke(ctx, "/unicom.api.v1.Unicom/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UnicomServer is the server API for Unicom service.
// All implementations should embed UnimplementedUnicomServer
// for forward compatibility
type UnicomServer interface {
	SendSync(context.Context, *SendSyncRequest) (*SendResponse, error)
	SendAsync(context.Context, *SendAsyncRequest) (*SendResponse, error)
	GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error)
}

// UnimplementedUnicomServer should be embedded to have forward compatible implementations.
type UnimplementedUnicomServer struct {
}

func (UnimplementedUnicomServer) SendSync(context.Context, *SendSyncRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSync not implemented")
}
func (UnimplementedUnicomServer) SendAsync(context.Context, *SendAsyncRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendAsync not implemented")
}
func (UnimplementedUnicomServer) GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}

// UnsafeUnicomServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UnicomServer will
// result in compilation errors.
type UnsafeUnicomServer interface {
	mustEmbedUnimplementedUnicomServer()
}

func RegisterUnicomServer(s grpc.ServiceRegistrar, srv UnicomServer) {
	s.RegisterService(&Unicom_ServiceDesc, srv)
}

func _Unicom_SendSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendSyncRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnicomServer).SendSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/unicom.api.v1.Unicom/SendSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnicomServer).SendSync(ctx, req.(*SendSyncRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Unicom_SendAsync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendAsyncRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnicomServer).SendAsync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/unicom.api.v1.Unicom/SendAsync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnicomServer).SendAsync(ctx, req.(*SendAsyncRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Unicom_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnicomServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/unicom.api.v1.Unicom/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnicomServer).GetStatus(ctx, req.(*GetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Unicom_ServiceDesc is the grpc.ServiceDesc for Unicom service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Unicom_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "unicom.api.v1.Unicom",
	HandlerType: (*UnicomServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSync",
			Handler:    _Unicom_SendSync_Handler,
		},
		{
			MethodName: "SendAsync",
			Handler:    _Unicom_SendAsync_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _Unicom_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "unicom/api/v1/service.proto",
}
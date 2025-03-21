// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: authorization_sercvice.proto

package authorizationpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Authorization_TryAuthorization_FullMethodName = "/authorization.Authorization/TryAuthorization"
)

// AuthorizationClient is the client API for Authorization service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorizationClient interface {
	TryAuthorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error)
}

type authorizationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorizationClient(cc grpc.ClientConnInterface) AuthorizationClient {
	return &authorizationClient{cc}
}

func (c *authorizationClient) TryAuthorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizationResponse)
	err := c.cc.Invoke(ctx, Authorization_TryAuthorization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorizationServer is the server API for Authorization service.
// All implementations must embed UnimplementedAuthorizationServer
// for forward compatibility.
type AuthorizationServer interface {
	TryAuthorization(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error)
	mustEmbedUnimplementedAuthorizationServer()
}

// UnimplementedAuthorizationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthorizationServer struct{}

func (UnimplementedAuthorizationServer) TryAuthorization(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryAuthorization not implemented")
}
func (UnimplementedAuthorizationServer) mustEmbedUnimplementedAuthorizationServer() {}
func (UnimplementedAuthorizationServer) testEmbeddedByValue()                       {}

// UnsafeAuthorizationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorizationServer will
// result in compilation errors.
type UnsafeAuthorizationServer interface {
	mustEmbedUnimplementedAuthorizationServer()
}

func RegisterAuthorizationServer(s grpc.ServiceRegistrar, srv AuthorizationServer) {
	// If the following call pancis, it indicates UnimplementedAuthorizationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Authorization_ServiceDesc, srv)
}

func _Authorization_TryAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).TryAuthorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Authorization_TryAuthorization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).TryAuthorization(ctx, req.(*AuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Authorization_ServiceDesc is the grpc.ServiceDesc for Authorization service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authorization_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authorization.Authorization",
	HandlerType: (*AuthorizationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TryAuthorization",
			Handler:    _Authorization_TryAuthorization_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authorization_sercvice.proto",
}

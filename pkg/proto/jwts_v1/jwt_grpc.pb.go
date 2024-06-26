// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: jwts_v1/jwt.proto

package jwts_v1

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

// JwtClient is the client API for Jwt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JwtClient interface {
	Create(ctx context.Context, in *JwtCreateReq, opts ...grpc.CallOption) (*JwtCreateRep, error)
	Validate(ctx context.Context, in *JwtValidateReq, opts ...grpc.CallOption) (*JwtValidateRep, error)
}

type jwtClient struct {
	cc grpc.ClientConnInterface
}

func NewJwtClient(cc grpc.ClientConnInterface) JwtClient {
	return &jwtClient{cc}
}

func (c *jwtClient) Create(ctx context.Context, in *JwtCreateReq, opts ...grpc.CallOption) (*JwtCreateRep, error) {
	out := new(JwtCreateRep)
	err := c.cc.Invoke(ctx, "/jwts_v1.Jwt/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jwtClient) Validate(ctx context.Context, in *JwtValidateReq, opts ...grpc.CallOption) (*JwtValidateRep, error) {
	out := new(JwtValidateRep)
	err := c.cc.Invoke(ctx, "/jwts_v1.Jwt/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JwtServer is the server API for Jwt service.
// All implementations must embed UnimplementedJwtServer
// for forward compatibility
type JwtServer interface {
	Create(context.Context, *JwtCreateReq) (*JwtCreateRep, error)
	Validate(context.Context, *JwtValidateReq) (*JwtValidateRep, error)
	mustEmbedUnimplementedJwtServer()
}

// UnimplementedJwtServer must be embedded to have forward compatible implementations.
type UnimplementedJwtServer struct {
}

func (UnimplementedJwtServer) Create(context.Context, *JwtCreateReq) (*JwtCreateRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedJwtServer) Validate(context.Context, *JwtValidateReq) (*JwtValidateRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedJwtServer) mustEmbedUnimplementedJwtServer() {}

// UnsafeJwtServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JwtServer will
// result in compilation errors.
type UnsafeJwtServer interface {
	mustEmbedUnimplementedJwtServer()
}

func RegisterJwtServer(s grpc.ServiceRegistrar, srv JwtServer) {
	s.RegisterService(&Jwt_ServiceDesc, srv)
}

func _Jwt_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JwtCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JwtServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jwts_v1.Jwt/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JwtServer).Create(ctx, req.(*JwtCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jwt_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JwtValidateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JwtServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jwts_v1.Jwt/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JwtServer).Validate(ctx, req.(*JwtValidateReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Jwt_ServiceDesc is the grpc.ServiceDesc for Jwt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Jwt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jwts_v1.Jwt",
	HandlerType: (*JwtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Jwt_Create_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Jwt_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "jwts_v1/jwt.proto",
}

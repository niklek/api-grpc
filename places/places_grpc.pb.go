// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// PlacesClient is the client API for Places service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlacesClient interface {
	GetById(ctx context.Context, in *PlaceIdRequest, opts ...grpc.CallOption) (*PlaceResponse, error)
}

type placesClient struct {
	cc grpc.ClientConnInterface
}

func NewPlacesClient(cc grpc.ClientConnInterface) PlacesClient {
	return &placesClient{cc}
}

func (c *placesClient) GetById(ctx context.Context, in *PlaceIdRequest, opts ...grpc.CallOption) (*PlaceResponse, error) {
	out := new(PlaceResponse)
	err := c.cc.Invoke(ctx, "/places.Places/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlacesServer is the server API for Places service.
// All implementations must embed UnimplementedPlacesServer
// for forward compatibility
type PlacesServer interface {
	GetById(context.Context, *PlaceIdRequest) (*PlaceResponse, error)
	mustEmbedUnimplementedPlacesServer()
}

// UnimplementedPlacesServer must be embedded to have forward compatible implementations.
type UnimplementedPlacesServer struct {
}

func (UnimplementedPlacesServer) GetById(context.Context, *PlaceIdRequest) (*PlaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedPlacesServer) mustEmbedUnimplementedPlacesServer() {}

// UnsafePlacesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlacesServer will
// result in compilation errors.
type UnsafePlacesServer interface {
	mustEmbedUnimplementedPlacesServer()
}

func RegisterPlacesServer(s grpc.ServiceRegistrar, srv PlacesServer) {
	s.RegisterService(&_Places_serviceDesc, srv)
}

func _Places_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlacesServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/places.Places/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlacesServer).GetById(ctx, req.(*PlaceIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Places_serviceDesc = grpc.ServiceDesc{
	ServiceName: "places.Places",
	HandlerType: (*PlacesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetById",
			Handler:    _Places_GetById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "places.proto",
}

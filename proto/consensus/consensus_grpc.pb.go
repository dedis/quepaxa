// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package consensus

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

// ConsensusClient is the client API for Consensus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsensusClient interface {
	ESP(ctx context.Context, in *ProposerMessage, opts ...grpc.CallOption) (*RecorderResponse, error)
}

type consensusClient struct {
	cc grpc.ClientConnInterface
}

func NewConsensusClient(cc grpc.ClientConnInterface) ConsensusClient {
	return &consensusClient{cc}
}

func (c *consensusClient) ESP(ctx context.Context, in *ProposerMessage, opts ...grpc.CallOption) (*RecorderResponse, error) {
	out := new(RecorderResponse)
	err := c.cc.Invoke(ctx, "/Consensus/ESP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusServer is the server API for Consensus service.
// All implementations must embed UnimplementedConsensusServer
// for forward compatibility
type ConsensusServer interface {
	ESP(context.Context, *ProposerMessage) (*RecorderResponse, error)
	mustEmbedUnimplementedConsensusServer()
}

// UnimplementedConsensusServer must be embedded to have forward compatible implementations.
type UnimplementedConsensusServer struct {
}

func (UnimplementedConsensusServer) ESP(context.Context, *ProposerMessage) (*RecorderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ESP not implemented")
}
func (UnimplementedConsensusServer) mustEmbedUnimplementedConsensusServer() {}

// UnsafeConsensusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsensusServer will
// result in compilation errors.
type UnsafeConsensusServer interface {
	mustEmbedUnimplementedConsensusServer()
}

func RegisterConsensusServer(s grpc.ServiceRegistrar, srv ConsensusServer) {
	s.RegisterService(&Consensus_ServiceDesc, srv)
}

func _Consensus_ESP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProposerMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).ESP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus/ESP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).ESP(ctx, req.(*ProposerMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Consensus_ServiceDesc is the grpc.ServiceDesc for Consensus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Consensus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Consensus",
	HandlerType: (*ConsensusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ESP",
			Handler:    _Consensus_ESP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/consensus/consensus.proto",
}
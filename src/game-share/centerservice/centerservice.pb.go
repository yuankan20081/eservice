// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/game-share/centerservice/centerservice.proto

/*
Package centerservice is a generated protocol buffer package.

It is generated from these files:
	src/game-share/centerservice/centerservice.proto

It has these top-level messages:
	AgentAuthRequest
	AgentAuthReply
*/
package centerservice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AgentAuthReply_AuthCode int32

const (
	AgentAuthReply_SUCCESS AgentAuthReply_AuthCode = 0
	AgentAuthReply_FAIL    AgentAuthReply_AuthCode = 1
)

var AgentAuthReply_AuthCode_name = map[int32]string{
	0: "SUCCESS",
	1: "FAIL",
}
var AgentAuthReply_AuthCode_value = map[string]int32{
	"SUCCESS": 0,
	"FAIL":    1,
}

func (x AgentAuthReply_AuthCode) String() string {
	return proto.EnumName(AgentAuthReply_AuthCode_name, int32(x))
}
func (AgentAuthReply_AuthCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type AgentAuthRequest struct {
	Ticket string `protobuf:"bytes,1,opt,name=ticket" json:"ticket,omitempty"`
	Ip     string `protobuf:"bytes,2,opt,name=ip" json:"ip,omitempty"`
}

func (m *AgentAuthRequest) Reset()                    { *m = AgentAuthRequest{} }
func (m *AgentAuthRequest) String() string            { return proto.CompactTextString(m) }
func (*AgentAuthRequest) ProtoMessage()               {}
func (*AgentAuthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AgentAuthRequest) GetTicket() string {
	if m != nil {
		return m.Ticket
	}
	return ""
}

func (m *AgentAuthRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type AgentAuthReply struct {
	Token  string                  `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Server string                  `protobuf:"bytes,2,opt,name=server" json:"server,omitempty"`
	Code   AgentAuthReply_AuthCode `protobuf:"varint,4,opt,name=code,enum=centerservice.AgentAuthReply_AuthCode" json:"code,omitempty"`
}

func (m *AgentAuthReply) Reset()                    { *m = AgentAuthReply{} }
func (m *AgentAuthReply) String() string            { return proto.CompactTextString(m) }
func (*AgentAuthReply) ProtoMessage()               {}
func (*AgentAuthReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AgentAuthReply) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AgentAuthReply) GetServer() string {
	if m != nil {
		return m.Server
	}
	return ""
}

func (m *AgentAuthReply) GetCode() AgentAuthReply_AuthCode {
	if m != nil {
		return m.Code
	}
	return AgentAuthReply_SUCCESS
}

func init() {
	proto.RegisterType((*AgentAuthRequest)(nil), "centerservice.AgentAuthRequest")
	proto.RegisterType((*AgentAuthReply)(nil), "centerservice.AgentAuthReply")
	proto.RegisterEnum("centerservice.AgentAuthReply_AuthCode", AgentAuthReply_AuthCode_name, AgentAuthReply_AuthCode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CenterService service

type CenterServiceClient interface {
	AgentAuth(ctx context.Context, in *AgentAuthRequest, opts ...grpc.CallOption) (*AgentAuthReply, error)
}

type centerServiceClient struct {
	cc *grpc.ClientConn
}

func NewCenterServiceClient(cc *grpc.ClientConn) CenterServiceClient {
	return &centerServiceClient{cc}
}

func (c *centerServiceClient) AgentAuth(ctx context.Context, in *AgentAuthRequest, opts ...grpc.CallOption) (*AgentAuthReply, error) {
	out := new(AgentAuthReply)
	err := grpc.Invoke(ctx, "/centerservice.CenterService/AgentAuth", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CenterService service

type CenterServiceServer interface {
	AgentAuth(context.Context, *AgentAuthRequest) (*AgentAuthReply, error)
}

func RegisterCenterServiceServer(s *grpc.Server, srv CenterServiceServer) {
	s.RegisterService(&_CenterService_serviceDesc, srv)
}

func _CenterService_AgentAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CenterServiceServer).AgentAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/centerservice.CenterService/AgentAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CenterServiceServer).AgentAuth(ctx, req.(*AgentAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CenterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "centerservice.CenterService",
	HandlerType: (*CenterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AgentAuth",
			Handler:    _CenterService_AgentAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/game-share/centerservice/centerservice.proto",
}

func init() { proto.RegisterFile("src/game-share/centerservice/centerservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x28, 0x2e, 0x4a, 0xd6,
	0x4f, 0x4f, 0xcc, 0x4d, 0xd5, 0x2d, 0xce, 0x48, 0x2c, 0x4a, 0xd5, 0x4f, 0x4e, 0xcd, 0x2b, 0x49,
	0x2d, 0x2a, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x46, 0xe3, 0xe9, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b,
	0xf1, 0xa2, 0x08, 0x2a, 0x59, 0x71, 0x09, 0x38, 0xa6, 0xa7, 0xe6, 0x95, 0x38, 0x96, 0x96, 0x64,
	0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x89, 0x71, 0xb1, 0x95, 0x64, 0x26, 0x67, 0xa7,
	0x96, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x79, 0x42, 0x7c, 0x5c, 0x4c, 0x99, 0x05,
	0x12, 0x4c, 0x60, 0x31, 0xa6, 0xcc, 0x02, 0xa5, 0xb9, 0x8c, 0x5c, 0x7c, 0x48, 0x9a, 0x0b, 0x72,
	0x2a, 0x85, 0x44, 0xb8, 0x58, 0x4b, 0xf2, 0xb3, 0x53, 0xf3, 0xa0, 0x3a, 0x21, 0x1c, 0x90, 0x81,
	0x20, 0xfb, 0x52, 0x8b, 0xa0, 0x9a, 0xa1, 0x3c, 0x21, 0x2b, 0x2e, 0x96, 0xe4, 0xfc, 0x94, 0x54,
	0x09, 0x16, 0x05, 0x46, 0x0d, 0x3e, 0x23, 0x35, 0x3d, 0x54, 0xf7, 0xa2, 0x1a, 0xad, 0x07, 0x62,
	0x39, 0xe7, 0xa7, 0xa4, 0x06, 0x81, 0xf5, 0x28, 0x29, 0x72, 0x71, 0xc0, 0x44, 0x84, 0xb8, 0xb9,
	0xd8, 0x83, 0x43, 0x9d, 0x9d, 0x5d, 0x83, 0x83, 0x05, 0x18, 0x84, 0x38, 0xb8, 0x58, 0xdc, 0x1c,
	0x3d, 0x7d, 0x04, 0x18, 0x8d, 0xe2, 0xb8, 0x78, 0x9d, 0xc1, 0x26, 0x06, 0x43, 0x4c, 0x14, 0xf2,
	0xe5, 0xe2, 0x84, 0x1b, 0x2a, 0x24, 0x8f, 0xdb, 0x3a, 0x70, 0x30, 0x48, 0xc9, 0xe2, 0x75, 0x8f,
	0x12, 0x43, 0x12, 0x1b, 0x38, 0x44, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcd, 0xa7, 0xd2,
	0x08, 0x85, 0x01, 0x00, 0x00,
}
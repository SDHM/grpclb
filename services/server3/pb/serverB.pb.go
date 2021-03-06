// Code generated by protoc-gen-go. DO NOT EDIT.
// source: serverB.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	serverB.proto

It has these top-level messages:
	ServierBRequest
	ServierBReply
*/
package pb

import (
	fmt "fmt"
	"grpclb/tracing"

	proto "github.com/golang/protobuf/proto"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"

	math "math"

	context "golang.org/x/net/context"

	"github.com/opentracing/opentracing-go/ext"
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

// The request message containing the user's name.
type ServierBRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *ServierBRequest) Reset()                    { *m = ServierBRequest{} }
func (m *ServierBRequest) String() string            { return proto.CompactTextString(m) }
func (*ServierBRequest) ProtoMessage()               {}
func (*ServierBRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServierBRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type ServierBReply struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *ServierBReply) Reset()                    { *m = ServierBReply{} }
func (m *ServierBReply) String() string            { return proto.CompactTextString(m) }
func (*ServierBReply) ProtoMessage()               {}
func (*ServierBReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServierBReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ServierBRequest)(nil), "pb.ServierBRequest")
	proto.RegisterType((*ServierBReply)(nil), "pb.ServierBReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

//hello this is my test This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ServerB service

type ServerBClient interface {
	//   Sends a greeting
	ServerBFunc(ctx context.Context, in *ServierBRequest, opts ...grpc.CallOption) (*ServierBReply, error)
}

type serverBClient struct {
	cc     *grpc.ClientConn
	tracer opentracing.Tracer
}

func NewServerBClient(tracer opentracing.Tracer, cc *grpc.ClientConn) ServerBClient {
	return &serverBClient{
		cc:     cc,
		tracer: tracer, // self-add
	}
}

func (c *serverBClient) ServerBFunc(ctx context.Context, in *ServierBRequest, opts ...grpc.CallOption) (*ServierBReply, error) {
	out := new(ServierBReply)

	// self-add
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); nil == span {
		span, ctx = opentracing.StartSpanFromContext(ctx, "client call /pb.Greeter/SayHello")
		ext.SpanKindRPCServer.Set(span)
	}

	if nil != span {
		toGRPCFunc := tracing.ToGRPCRequest(c.tracer)
		md := metadata.Pairs()
		ctx = toGRPCFunc(ctx, &md)
		ctx = metadata.NewContext(ctx, md)
		defer span.Finish()
	}
	// self-add

	err := grpc.Invoke(ctx, "/pb.ServerB/ServerBFunc", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ServerB service

type ServerBServer interface {
	//   Sends a greeting
	ServerBFunc(context.Context, *ServierBRequest) (*ServierBReply, error)
}

func RegisterServerBServer(s *grpc.Server, srv ServerBServer) {
	s.RegisterService(&_ServerB_serviceDesc, srv)
}

func _ServerB_ServerBFunc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServierBRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerBServer).ServerBFunc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ServerB/ServerBFunc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerBServer).ServerBFunc(ctx, req.(*ServierBRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServerB_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ServerB",
	HandlerType: (*ServerBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServerBFunc",
			Handler:    _ServerB_ServerBFunc_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "serverB.proto",
}

func init() { proto.RegisterFile("serverB.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0x72, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x52, 0xe5,
	0xe2, 0x0f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x2d, 0x72, 0x0a, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e,
	0x11, 0x12, 0xe2, 0x62, 0xc9, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02,
	0xb3, 0x95, 0x34, 0xb9, 0x78, 0x11, 0xca, 0x0a, 0x72, 0x2a, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53,
	0x8b, 0x8b, 0x13, 0xd3, 0x61, 0xea, 0x60, 0x5c, 0x23, 0x27, 0x2e, 0xf6, 0x60, 0x88, 0x35, 0x42,
	0xe6, 0x5c, 0xdc, 0x50, 0xa6, 0x5b, 0x69, 0x5e, 0xb2, 0x90, 0xb0, 0x5e, 0x41, 0x92, 0x1e, 0x9a,
	0x6d, 0x52, 0x82, 0xa8, 0x82, 0x05, 0x39, 0x95, 0x4a, 0x0c, 0x4e, 0x3a, 0x5c, 0xa2, 0xc9, 0xf9,
	0xb9, 0x7a, 0x15, 0x99, 0xc5, 0x19, 0x99, 0x85, 0xa5, 0x7a, 0x25, 0xa9, 0xc5, 0x25, 0x7a, 0xe9,
	0x45, 0x05, 0xc9, 0x4e, 0x30, 0xa3, 0x03, 0x18, 0x17, 0x31, 0xc1, 0xd8, 0x49, 0x6c, 0x60, 0xef,
	0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x7b, 0x2d, 0xc6, 0xdf, 0x00, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/v1/wait.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type WaitRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WaitRequest) Reset()         { *m = WaitRequest{} }
func (m *WaitRequest) String() string { return proto.CompactTextString(m) }
func (*WaitRequest) ProtoMessage()    {}
func (*WaitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9d6360f924548db, []int{0}
}

func (m *WaitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitRequest.Unmarshal(m, b)
}
func (m *WaitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitRequest.Marshal(b, m, deterministic)
}
func (m *WaitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitRequest.Merge(m, src)
}
func (m *WaitRequest) XXX_Size() int {
	return xxx_messageInfo_WaitRequest.Size(m)
}
func (m *WaitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WaitRequest proto.InternalMessageInfo

func (m *WaitRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type WaitResponse struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Allocated            bool     `protobuf:"varint,2,opt,name=allocated,proto3" json:"allocated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WaitResponse) Reset()         { *m = WaitResponse{} }
func (m *WaitResponse) String() string { return proto.CompactTextString(m) }
func (*WaitResponse) ProtoMessage()    {}
func (*WaitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9d6360f924548db, []int{1}
}

func (m *WaitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WaitResponse.Unmarshal(m, b)
}
func (m *WaitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WaitResponse.Marshal(b, m, deterministic)
}
func (m *WaitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WaitResponse.Merge(m, src)
}
func (m *WaitResponse) XXX_Size() int {
	return xxx_messageInfo_WaitResponse.Size(m)
}
func (m *WaitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WaitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WaitResponse proto.InternalMessageInfo

func (m *WaitResponse) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *WaitResponse) GetAllocated() bool {
	if m != nil {
		return m.Allocated
	}
	return false
}

func init() {
	proto.RegisterType((*WaitRequest)(nil), "v1.WaitRequest")
	proto.RegisterType((*WaitResponse)(nil), "v1.WaitResponse")
}

func init() { proto.RegisterFile("api/v1/wait.proto", fileDescriptor_c9d6360f924548db) }

var fileDescriptor_c9d6360f924548db = []byte{
	// 151 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2f, 0x33, 0xd4, 0x2f, 0x4f, 0xcc, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a,
	0x33, 0x54, 0x52, 0xe4, 0xe2, 0x0e, 0x4f, 0xcc, 0x2c, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e,
	0x11, 0x12, 0xe2, 0x62, 0x29, 0x2d, 0xcd, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02,
	0xb3, 0x95, 0x1c, 0xb8, 0x78, 0x20, 0x4a, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0xb1, 0xa9, 0x11,
	0x92, 0xe1, 0xe2, 0x4c, 0xcc, 0xc9, 0xc9, 0x4f, 0x4e, 0x2c, 0x49, 0x4d, 0x91, 0x60, 0x52, 0x60,
	0xd4, 0xe0, 0x08, 0x42, 0x08, 0x18, 0x19, 0x73, 0xb1, 0x80, 0x4c, 0x10, 0xd2, 0x86, 0xd2, 0xfc,
	0x7a, 0x65, 0x86, 0x7a, 0x48, 0xd6, 0x4a, 0x09, 0x20, 0x04, 0x20, 0x96, 0x18, 0x30, 0x3a, 0x71,
	0x44, 0xb1, 0x41, 0x9c, 0x9c, 0xc4, 0x06, 0x76, 0xae, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x2b,
	0x68, 0x10, 0x39, 0xc3, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WaitClient is the client API for Wait service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WaitClient interface {
	Wait(ctx context.Context, in *WaitRequest, opts ...grpc.CallOption) (Wait_WaitClient, error)
}

type waitClient struct {
	cc grpc.ClientConnInterface
}

func NewWaitClient(cc grpc.ClientConnInterface) WaitClient {
	return &waitClient{cc}
}

func (c *waitClient) Wait(ctx context.Context, in *WaitRequest, opts ...grpc.CallOption) (Wait_WaitClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Wait_serviceDesc.Streams[0], "/v1.Wait/Wait", opts...)
	if err != nil {
		return nil, err
	}
	x := &waitWaitClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Wait_WaitClient interface {
	Recv() (*WaitResponse, error)
	grpc.ClientStream
}

type waitWaitClient struct {
	grpc.ClientStream
}

func (x *waitWaitClient) Recv() (*WaitResponse, error) {
	m := new(WaitResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WaitServer is the server API for Wait service.
type WaitServer interface {
	Wait(*WaitRequest, Wait_WaitServer) error
}

// UnimplementedWaitServer can be embedded to have forward compatible implementations.
type UnimplementedWaitServer struct {
}

func (*UnimplementedWaitServer) Wait(req *WaitRequest, srv Wait_WaitServer) error {
	return status.Errorf(codes.Unimplemented, "method Wait not implemented")
}

func RegisterWaitServer(s *grpc.Server, srv WaitServer) {
	s.RegisterService(&_Wait_serviceDesc, srv)
}

func _Wait_Wait_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WaitRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WaitServer).Wait(m, &waitWaitServer{stream})
}

type Wait_WaitServer interface {
	Send(*WaitResponse) error
	grpc.ServerStream
}

type waitWaitServer struct {
	grpc.ServerStream
}

func (x *waitWaitServer) Send(m *WaitResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Wait_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Wait",
	HandlerType: (*WaitServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Wait",
			Handler:       _Wait_Wait_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/wait.proto",
}

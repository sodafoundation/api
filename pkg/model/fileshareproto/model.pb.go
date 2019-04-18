// Code generated by protoc-gen-go. DO NOT EDIT.
// source: model.proto

package fileshareproto

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

// CreateFileShareOpts is a structure which indicates all required properties for creating a file share.
type CreateFileShareOpts struct {
	// The uuid of the file share, optional when creating.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the file share, required.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The requested capacity of the file share, required.
	Size int64 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// The description of the file share, optional.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// The locality that file share belongs to, required.
	AvailabilityZone string `protobuf:"bytes,6,opt,name=availabilityZone,proto3" json:"availabilityZone,omitempty"`
	// The service level that file share belongs to, required.
	ProfileId string `protobuf:"bytes,7,opt,name=profileId,proto3" json:"profileId,omitempty"`
	// The uuid of the pool on which file share will be created, required.
	PoolId string `protobuf:"bytes,8,opt,name=poolId,proto3" json:"poolId,omitempty"`
	// The name of the pool on which file share will be created, required.
	PoolName string `protobuf:"bytes,9,opt,name=poolName,proto3" json:"poolName,omitempty"`
	// The metadata of the file share, optional.
	Metadata map[string]string `protobuf:"bytes,10,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The storage driver type.
	DriverName string `protobuf:"bytes,11,opt,name=driverName,proto3" json:"driverName,omitempty"`
	// The Context
	Context              string   `protobuf:"bytes,12,opt,name=context,proto3" json:"context,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateFileShareOpts) Reset()         { *m = CreateFileShareOpts{} }
func (m *CreateFileShareOpts) String() string { return proto.CompactTextString(m) }
func (*CreateFileShareOpts) ProtoMessage()    {}
func (*CreateFileShareOpts) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{0}
}

func (m *CreateFileShareOpts) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateFileShareOpts.Unmarshal(m, b)
}
func (m *CreateFileShareOpts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateFileShareOpts.Marshal(b, m, deterministic)
}
func (m *CreateFileShareOpts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateFileShareOpts.Merge(m, src)
}
func (m *CreateFileShareOpts) XXX_Size() int {
	return xxx_messageInfo_CreateFileShareOpts.Size(m)
}
func (m *CreateFileShareOpts) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateFileShareOpts.DiscardUnknown(m)
}

var xxx_messageInfo_CreateFileShareOpts proto.InternalMessageInfo

func (m *CreateFileShareOpts) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CreateFileShareOpts) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateFileShareOpts) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *CreateFileShareOpts) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateFileShareOpts) GetAvailabilityZone() string {
	if m != nil {
		return m.AvailabilityZone
	}
	return ""
}

func (m *CreateFileShareOpts) GetProfileId() string {
	if m != nil {
		return m.ProfileId
	}
	return ""
}

func (m *CreateFileShareOpts) GetPoolId() string {
	if m != nil {
		return m.PoolId
	}
	return ""
}

func (m *CreateFileShareOpts) GetPoolName() string {
	if m != nil {
		return m.PoolName
	}
	return ""
}

func (m *CreateFileShareOpts) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *CreateFileShareOpts) GetDriverName() string {
	if m != nil {
		return m.DriverName
	}
	return ""
}

func (m *CreateFileShareOpts) GetContext() string {
	if m != nil {
		return m.Context
	}
	return ""
}

type HostInfo struct {
	// The platform of the host, such as "x86_64"
	Platform string `protobuf:"bytes,1,opt,name=platform,proto3" json:"platform,omitempty"`
	// The type of OS, such as "linux","windows", etc.
	OsType string `protobuf:"bytes,2,opt,name=osType,proto3" json:"osType,omitempty"`
	// The name of the host
	Host string `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
	// The ip address of the host
	Ip string `protobuf:"bytes,4,opt,name=ip,proto3" json:"ip,omitempty"`
	// The initiator infomation, such as: "iqn.2017.com.redhat:e08039b48d5c"
	Initiator            string   `protobuf:"bytes,5,opt,name=initiator,proto3" json:"initiator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HostInfo) Reset()         { *m = HostInfo{} }
func (m *HostInfo) String() string { return proto.CompactTextString(m) }
func (*HostInfo) ProtoMessage()    {}
func (*HostInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{1}
}

func (m *HostInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HostInfo.Unmarshal(m, b)
}
func (m *HostInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HostInfo.Marshal(b, m, deterministic)
}
func (m *HostInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HostInfo.Merge(m, src)
}
func (m *HostInfo) XXX_Size() int {
	return xxx_messageInfo_HostInfo.Size(m)
}
func (m *HostInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_HostInfo.DiscardUnknown(m)
}

var xxx_messageInfo_HostInfo proto.InternalMessageInfo

func (m *HostInfo) GetPlatform() string {
	if m != nil {
		return m.Platform
	}
	return ""
}

func (m *HostInfo) GetOsType() string {
	if m != nil {
		return m.OsType
	}
	return ""
}

func (m *HostInfo) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *HostInfo) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *HostInfo) GetInitiator() string {
	if m != nil {
		return m.Initiator
	}
	return ""
}

type FileShareData struct {
	Data                 map[string]string `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FileShareData) Reset()         { *m = FileShareData{} }
func (m *FileShareData) String() string { return proto.CompactTextString(m) }
func (*FileShareData) ProtoMessage()    {}
func (*FileShareData) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{2}
}

func (m *FileShareData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileShareData.Unmarshal(m, b)
}
func (m *FileShareData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileShareData.Marshal(b, m, deterministic)
}
func (m *FileShareData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileShareData.Merge(m, src)
}
func (m *FileShareData) XXX_Size() int {
	return xxx_messageInfo_FileShareData.Size(m)
}
func (m *FileShareData) XXX_DiscardUnknown() {
	xxx_messageInfo_FileShareData.DiscardUnknown(m)
}

var xxx_messageInfo_FileShareData proto.InternalMessageInfo

func (m *FileShareData) GetData() map[string]string {
	if m != nil {
		return m.Data
	}
	return nil
}

// Generic response, it return:
// 1. Return result with message when create/update resource successfully.
// 2. Return result without message when delete resource successfully.
// 3. Return Error with error code and message when operate unsuccessfully.
type GenericResponse struct {
	// Types that are valid to be assigned to Reply:
	//	*GenericResponse_Result_
	//	*GenericResponse_Error_
	Reply                isGenericResponse_Reply `protobuf_oneof:"reply"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *GenericResponse) Reset()         { *m = GenericResponse{} }
func (m *GenericResponse) String() string { return proto.CompactTextString(m) }
func (*GenericResponse) ProtoMessage()    {}
func (*GenericResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{3}
}

func (m *GenericResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenericResponse.Unmarshal(m, b)
}
func (m *GenericResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenericResponse.Marshal(b, m, deterministic)
}
func (m *GenericResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenericResponse.Merge(m, src)
}
func (m *GenericResponse) XXX_Size() int {
	return xxx_messageInfo_GenericResponse.Size(m)
}
func (m *GenericResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenericResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenericResponse proto.InternalMessageInfo

type isGenericResponse_Reply interface {
	isGenericResponse_Reply()
}

type GenericResponse_Result_ struct {
	Result *GenericResponse_Result `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type GenericResponse_Error_ struct {
	Error *GenericResponse_Error `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*GenericResponse_Result_) isGenericResponse_Reply() {}

func (*GenericResponse_Error_) isGenericResponse_Reply() {}

func (m *GenericResponse) GetReply() isGenericResponse_Reply {
	if m != nil {
		return m.Reply
	}
	return nil
}

func (m *GenericResponse) GetResult() *GenericResponse_Result {
	if x, ok := m.GetReply().(*GenericResponse_Result_); ok {
		return x.Result
	}
	return nil
}

func (m *GenericResponse) GetError() *GenericResponse_Error {
	if x, ok := m.GetReply().(*GenericResponse_Error_); ok {
		return x.Error
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*GenericResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*GenericResponse_Result_)(nil),
		(*GenericResponse_Error_)(nil),
	}
}

type GenericResponse_Result struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenericResponse_Result) Reset()         { *m = GenericResponse_Result{} }
func (m *GenericResponse_Result) String() string { return proto.CompactTextString(m) }
func (*GenericResponse_Result) ProtoMessage()    {}
func (*GenericResponse_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{3, 0}
}

func (m *GenericResponse_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenericResponse_Result.Unmarshal(m, b)
}
func (m *GenericResponse_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenericResponse_Result.Marshal(b, m, deterministic)
}
func (m *GenericResponse_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenericResponse_Result.Merge(m, src)
}
func (m *GenericResponse_Result) XXX_Size() int {
	return xxx_messageInfo_GenericResponse_Result.Size(m)
}
func (m *GenericResponse_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_GenericResponse_Result.DiscardUnknown(m)
}

var xxx_messageInfo_GenericResponse_Result proto.InternalMessageInfo

func (m *GenericResponse_Result) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GenericResponse_Error struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenericResponse_Error) Reset()         { *m = GenericResponse_Error{} }
func (m *GenericResponse_Error) String() string { return proto.CompactTextString(m) }
func (*GenericResponse_Error) ProtoMessage()    {}
func (*GenericResponse_Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{3, 1}
}

func (m *GenericResponse_Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenericResponse_Error.Unmarshal(m, b)
}
func (m *GenericResponse_Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenericResponse_Error.Marshal(b, m, deterministic)
}
func (m *GenericResponse_Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenericResponse_Error.Merge(m, src)
}
func (m *GenericResponse_Error) XXX_Size() int {
	return xxx_messageInfo_GenericResponse_Error.Size(m)
}
func (m *GenericResponse_Error) XXX_DiscardUnknown() {
	xxx_messageInfo_GenericResponse_Error.DiscardUnknown(m)
}

var xxx_messageInfo_GenericResponse_Error proto.InternalMessageInfo

func (m *GenericResponse_Error) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *GenericResponse_Error) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateFileShareOpts)(nil), "fileshareproto.CreateFileShareOpts")
	proto.RegisterMapType((map[string]string)(nil), "fileshareproto.CreateFileShareOpts.MetadataEntry")
	proto.RegisterType((*HostInfo)(nil), "fileshareproto.HostInfo")
	proto.RegisterType((*FileShareData)(nil), "fileshareproto.FileShareData")
	proto.RegisterMapType((map[string]string)(nil), "fileshareproto.FileShareData.DataEntry")
	proto.RegisterType((*GenericResponse)(nil), "fileshareproto.GenericResponse")
	proto.RegisterType((*GenericResponse_Result)(nil), "fileshareproto.GenericResponse.Result")
	proto.RegisterType((*GenericResponse_Error)(nil), "fileshareproto.GenericResponse.Error")
}

func init() { proto.RegisterFile("model.proto", fileDescriptor_4c16552f9fdb66d8) }

var fileDescriptor_4c16552f9fdb66d8 = []byte{
	// 550 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xcf, 0x6b, 0x13, 0x41,
	0x14, 0xc7, 0xb3, 0x9b, 0xdf, 0x6f, 0x6d, 0x5a, 0x46, 0x91, 0x65, 0x11, 0x0d, 0x11, 0x35, 0x78,
	0x58, 0x30, 0x1e, 0x14, 0x4b, 0x41, 0x6c, 0x5a, 0x93, 0x43, 0x55, 0x56, 0x2f, 0x7a, 0x9b, 0x66,
	0x5f, 0xda, 0x21, 0x93, 0x9d, 0x65, 0x66, 0x1a, 0x8c, 0x27, 0x2f, 0xfe, 0x45, 0xfe, 0x7f, 0x22,
	0xf3, 0xb2, 0x59, 0x9b, 0x54, 0xa8, 0x08, 0xbd, 0xbd, 0xef, 0xfb, 0x39, 0xf3, 0x79, 0x33, 0x10,
	0xcc, 0x55, 0x8a, 0x32, 0xce, 0xb5, 0xb2, 0x8a, 0x75, 0xa6, 0x42, 0xa2, 0x39, 0xe7, 0x1a, 0x49,
	0xf7, 0x7e, 0x56, 0xe1, 0xf6, 0xa1, 0x46, 0x6e, 0xf1, 0x58, 0x48, 0xfc, 0xe8, 0x02, 0xef, 0x73,
	0x6b, 0x58, 0x07, 0x7c, 0x91, 0x86, 0x5e, 0xd7, 0xeb, 0xb7, 0x13, 0x5f, 0xa4, 0x8c, 0x41, 0x2d,
	0xe3, 0x73, 0x0c, 0x7d, 0xf2, 0x90, 0xed, 0x7c, 0x46, 0x7c, 0xc3, 0xb0, 0xda, 0xf5, 0xfa, 0xd5,
	0x84, 0x6c, 0xd6, 0x85, 0x20, 0x45, 0x33, 0xd1, 0x22, 0xb7, 0x42, 0x65, 0x61, 0x8d, 0xd2, 0x2f,
	0xbb, 0xd8, 0x53, 0xd8, 0xe3, 0x0b, 0x2e, 0x24, 0x3f, 0x15, 0x52, 0xd8, 0xe5, 0x17, 0x95, 0x61,
	0xd8, 0xa0, 0xb4, 0x2b, 0x7e, 0x76, 0x0f, 0xda, 0xb9, 0x56, 0xee, 0xc8, 0xe3, 0x34, 0x6c, 0x52,
	0xd2, 0x1f, 0x07, 0xbb, 0x0b, 0x8d, 0x5c, 0x29, 0x39, 0x4e, 0xc3, 0x16, 0x85, 0x0a, 0xc5, 0x22,
	0x68, 0x39, 0xeb, 0x9d, 0x3b, 0x6f, 0x9b, 0x22, 0xa5, 0x66, 0x27, 0xd0, 0x9a, 0xa3, 0xe5, 0x29,
	0xb7, 0x3c, 0x84, 0x6e, 0xb5, 0x1f, 0x0c, 0x9e, 0xc5, 0x9b, 0x48, 0xe2, 0xbf, 0xe0, 0x88, 0x4f,
	0x8a, 0x9a, 0xa3, 0xcc, 0xea, 0x65, 0x52, 0xb6, 0x60, 0xf7, 0x01, 0x52, 0x2d, 0x16, 0xa8, 0x69,
	0x58, 0x40, 0xc3, 0x2e, 0x79, 0x58, 0x08, 0xcd, 0x89, 0xca, 0x2c, 0x7e, 0xb5, 0xe1, 0x2d, 0x0a,
	0xae, 0x65, 0xb4, 0x0f, 0x3b, 0x1b, 0x4d, 0xd9, 0x1e, 0x54, 0x67, 0xb8, 0x2c, 0x90, 0x3b, 0x93,
	0xdd, 0x81, 0xfa, 0x82, 0xcb, 0x8b, 0x35, 0xf4, 0x95, 0x78, 0xe5, 0xbf, 0xf4, 0x7a, 0xdf, 0x3d,
	0x68, 0x8d, 0x94, 0xb1, 0xe3, 0x6c, 0xaa, 0xe8, 0xba, 0x92, 0xdb, 0xa9, 0xd2, 0xf3, 0xa2, 0xba,
	0xd4, 0x0e, 0x91, 0x32, 0x9f, 0x96, 0xf9, 0xba, 0x47, 0xa1, 0xdc, 0xea, 0xce, 0x95, 0xb1, 0xb4,
	0xba, 0x76, 0x42, 0x36, 0xad, 0x3c, 0x2f, 0x36, 0xe6, 0x8b, 0xdc, 0xc1, 0x17, 0x99, 0xb0, 0x82,
	0x5b, 0xa5, 0xc3, 0xfa, 0x0a, 0x7e, 0xe9, 0xe8, 0xfd, 0xf0, 0x60, 0xa7, 0x64, 0x34, 0x74, 0x2c,
	0xf6, 0xa1, 0x46, 0x58, 0x3d, 0xc2, 0xfa, 0x64, 0x1b, 0xeb, 0x46, 0x72, 0x3c, 0x2c, 0x61, 0x52,
	0x51, 0xf4, 0x02, 0xda, 0xc3, 0xff, 0x42, 0xf1, 0xcb, 0x83, 0xdd, 0xb7, 0x98, 0xa1, 0x16, 0x93,
	0x04, 0x4d, 0xae, 0x32, 0x83, 0xec, 0x35, 0x34, 0x34, 0x9a, 0x0b, 0x69, 0xa9, 0x45, 0x30, 0x78,
	0xbc, 0x7d, 0x96, 0xad, 0x82, 0x38, 0xa1, 0xec, 0x51, 0x25, 0x29, 0xea, 0xd8, 0x01, 0xd4, 0x51,
	0x6b, 0xa5, 0x69, 0x5e, 0x30, 0x78, 0x74, 0x5d, 0x83, 0x23, 0x97, 0x3c, 0xaa, 0x24, 0xab, 0xaa,
	0xa8, 0x07, 0x8d, 0x55, 0x4b, 0xf7, 0x00, 0xe6, 0x68, 0x0c, 0x3f, 0xc3, 0xe2, 0x3a, 0x6b, 0x19,
	0x1d, 0x40, 0x9d, 0xaa, 0xdc, 0x2e, 0x26, 0x2a, 0x5d, 0xc7, 0xc9, 0xde, 0xfe, 0x46, 0xfe, 0x95,
	0x6f, 0xf4, 0xa6, 0x09, 0x75, 0x8d, 0xb9, 0x5c, 0x0e, 0xce, 0x00, 0x0e, 0x55, 0x66, 0xb5, 0x92,
	0x12, 0x35, 0xfb, 0x0c, 0xbb, 0x5b, 0xef, 0x97, 0x3d, 0xfc, 0x87, 0x07, 0x1e, 0x3d, 0xb8, 0xe6,
	0x86, 0xbd, 0xca, 0x60, 0x06, 0x9d, 0xe3, 0x0f, 0x5a, 0x2d, 0x84, 0x11, 0x2a, 0x1b, 0xaa, 0xc9,
	0xec, 0x06, 0x87, 0x9d, 0x36, 0x28, 0xf0, 0xfc, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x93, 0x3a,
	0x11, 0x5b, 0xbd, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ControllerClient is the client API for Controller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ControllerClient interface {
	// Create a file share
	CreateFileShare(ctx context.Context, in *CreateFileShareOpts, opts ...grpc.CallOption) (*GenericResponse, error)
}

type controllerClient struct {
	cc *grpc.ClientConn
}

func NewControllerClient(cc *grpc.ClientConn) ControllerClient {
	return &controllerClient{cc}
}

func (c *controllerClient) CreateFileShare(ctx context.Context, in *CreateFileShareOpts, opts ...grpc.CallOption) (*GenericResponse, error) {
	out := new(GenericResponse)
	err := c.cc.Invoke(ctx, "/fileshareproto.Controller/CreateFileShare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControllerServer is the server API for Controller service.
type ControllerServer interface {
	// Create a file share
	CreateFileShare(context.Context, *CreateFileShareOpts) (*GenericResponse, error)
}

// UnimplementedControllerServer can be embedded to have forward compatible implementations.
type UnimplementedControllerServer struct {
}

func (*UnimplementedControllerServer) CreateFileShare(ctx context.Context, req *CreateFileShareOpts) (*GenericResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFileShare not implemented")
}

func RegisterControllerServer(s *grpc.Server, srv ControllerServer) {
	s.RegisterService(&_Controller_serviceDesc, srv)
}

func _Controller_CreateFileShare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileShareOpts)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServer).CreateFileShare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshareproto.Controller/CreateFileShare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServer).CreateFileShare(ctx, req.(*CreateFileShareOpts))
	}
	return interceptor(ctx, in, info, handler)
}

var _Controller_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fileshareproto.Controller",
	HandlerType: (*ControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFileShare",
			Handler:    _Controller_CreateFileShare_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model.proto",
}

// FProvisionDockClient is the client API for FProvisionDock service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FProvisionDockClient interface {
	// Create a file share
	CreateFileShare(ctx context.Context, in *CreateFileShareOpts, opts ...grpc.CallOption) (*GenericResponse, error)
}

type fProvisionDockClient struct {
	cc *grpc.ClientConn
}

func NewFProvisionDockClient(cc *grpc.ClientConn) FProvisionDockClient {
	return &fProvisionDockClient{cc}
}

func (c *fProvisionDockClient) CreateFileShare(ctx context.Context, in *CreateFileShareOpts, opts ...grpc.CallOption) (*GenericResponse, error) {
	out := new(GenericResponse)
	err := c.cc.Invoke(ctx, "/fileshareproto.FProvisionDock/CreateFileShare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FProvisionDockServer is the server API for FProvisionDock service.
type FProvisionDockServer interface {
	// Create a file share
	CreateFileShare(context.Context, *CreateFileShareOpts) (*GenericResponse, error)
}

// UnimplementedFProvisionDockServer can be embedded to have forward compatible implementations.
type UnimplementedFProvisionDockServer struct {
}

func (*UnimplementedFProvisionDockServer) CreateFileShare(ctx context.Context, req *CreateFileShareOpts) (*GenericResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFileShare not implemented")
}

func RegisterFProvisionDockServer(s *grpc.Server, srv FProvisionDockServer) {
	s.RegisterService(&_FProvisionDock_serviceDesc, srv)
}

func _FProvisionDock_CreateFileShare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileShareOpts)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FProvisionDockServer).CreateFileShare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileshareproto.FProvisionDock/CreateFileShare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FProvisionDockServer).CreateFileShare(ctx, req.(*CreateFileShareOpts))
	}
	return interceptor(ctx, in, info, handler)
}

var _FProvisionDock_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fileshareproto.FProvisionDock",
	HandlerType: (*FProvisionDockServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFileShare",
			Handler:    _FProvisionDock_CreateFileShare_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model.proto",
}
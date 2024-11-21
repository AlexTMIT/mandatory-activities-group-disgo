// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.1
// source: replication.proto

package replication_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientName string `protobuf:"bytes,1,opt,name=clientName,proto3" json:"clientName,omitempty"`
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	mi := &file_replication_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinRequest.ProtoReflect.Descriptor instead.
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{0}
}

func (x *JoinRequest) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

type JoinReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *JoinReply) Reset() {
	*x = JoinReply{}
	mi := &file_replication_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JoinReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinReply) ProtoMessage() {}

func (x *JoinReply) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinReply.ProtoReflect.Descriptor instead.
func (*JoinReply) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{1}
}

func (x *JoinReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type BidRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount     int32  `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	ClientName string `protobuf:"bytes,2,opt,name=clientName,proto3" json:"clientName,omitempty"`
	Lamport    int32  `protobuf:"varint,3,opt,name=lamport,proto3" json:"lamport,omitempty"`
}

func (x *BidRequest) Reset() {
	*x = BidRequest{}
	mi := &file_replication_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BidRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BidRequest) ProtoMessage() {}

func (x *BidRequest) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BidRequest.ProtoReflect.Descriptor instead.
func (*BidRequest) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{2}
}

func (x *BidRequest) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *BidRequest) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

func (x *BidRequest) GetLamport() int32 {
	if x != nil {
		return x.Lamport
	}
	return 0
}

type BidReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	Lamport  int32  `protobuf:"varint,2,opt,name=lamport,proto3" json:"lamport,omitempty"`
}

func (x *BidReply) Reset() {
	*x = BidReply{}
	mi := &file_replication_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BidReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BidReply) ProtoMessage() {}

func (x *BidReply) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BidReply.ProtoReflect.Descriptor instead.
func (*BidReply) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{3}
}

func (x *BidReply) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

func (x *BidReply) GetLamport() int32 {
	if x != nil {
		return x.Lamport
	}
	return 0
}

type AQueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AQueryRequest) Reset() {
	*x = AQueryRequest{}
	mi := &file_replication_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AQueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AQueryRequest) ProtoMessage() {}

func (x *AQueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AQueryRequest.ProtoReflect.Descriptor instead.
func (*AQueryRequest) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{4}
}

type AQueryReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentAmount int32  `protobuf:"varint,1,opt,name=currentAmount,proto3" json:"currentAmount,omitempty"`
	Result        string `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *AQueryReply) Reset() {
	*x = AQueryReply{}
	mi := &file_replication_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AQueryReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AQueryReply) ProtoMessage() {}

func (x *AQueryReply) ProtoReflect() protoreflect.Message {
	mi := &file_replication_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AQueryReply.ProtoReflect.Descriptor instead.
func (*AQueryReply) Descriptor() ([]byte, []int) {
	return file_replication_proto_rawDescGZIP(), []int{5}
}

func (x *AQueryReply) GetCurrentAmount() int32 {
	if x != nil {
		return x.CurrentAmount
	}
	return 0
}

func (x *AQueryReply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_replication_proto protoreflect.FileDescriptor

var file_replication_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x2d, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x1d, 0x0a, 0x09, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x5e,
	0x0a, 0x0a, 0x42, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x40,
	0x0a, 0x08, 0x42, 0x69, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74,
	0x22, 0x0f, 0x0a, 0x0d, 0x41, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x4b, 0x0a, 0x0b, 0x41, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x24, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0xdf,
	0x01, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x07, 0x42, 0x69, 0x64, 0x64, 0x69, 0x6e, 0x67,
	0x12, 0x17, 0x2e, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x42,
	0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x72, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x69, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x48, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4a, 0x6f, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x2e, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4a,
	0x6f, 0x69, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0c, 0x41, 0x75,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1a, 0x2e, 0x72, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x42, 0x73, 0x5a, 0x71, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x6c, 0x65, 0x78, 0x54, 0x4d, 0x49, 0x54, 0x2f,
	0x6d, 0x61, 0x6e, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x79, 0x2d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x2d, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2d, 0x64, 0x69, 0x73, 0x67, 0x6f,
	0x2f, 0x74, 0x72, 0x65, 0x65, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x6d, 0x61, 0x6e, 0x64, 0x61,
	0x74, 0x6f, 0x72, 0x79, 0x2d, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2d,
	0x35, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_replication_proto_rawDescOnce sync.Once
	file_replication_proto_rawDescData = file_replication_proto_rawDesc
)

func file_replication_proto_rawDescGZIP() []byte {
	file_replication_proto_rawDescOnce.Do(func() {
		file_replication_proto_rawDescData = protoimpl.X.CompressGZIP(file_replication_proto_rawDescData)
	})
	return file_replication_proto_rawDescData
}

var file_replication_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_replication_proto_goTypes = []any{
	(*JoinRequest)(nil),   // 0: replication.JoinRequest
	(*JoinReply)(nil),     // 1: replication.JoinReply
	(*BidRequest)(nil),    // 2: replication.BidRequest
	(*BidReply)(nil),      // 3: replication.BidReply
	(*AQueryRequest)(nil), // 4: replication.AQueryRequest
	(*AQueryReply)(nil),   // 5: replication.AQueryReply
}
var file_replication_proto_depIdxs = []int32{
	2, // 0: replication.ReplicationService.Bidding:input_type -> replication.BidRequest
	0, // 1: replication.ReplicationService.ProcessJoinRequest:input_type -> replication.JoinRequest
	4, // 2: replication.ReplicationService.AuctionQuery:input_type -> replication.AQueryRequest
	3, // 3: replication.ReplicationService.Bidding:output_type -> replication.BidReply
	1, // 4: replication.ReplicationService.ProcessJoinRequest:output_type -> replication.JoinReply
	5, // 5: replication.ReplicationService.AuctionQuery:output_type -> replication.AQueryReply
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_replication_proto_init() }
func file_replication_proto_init() {
	if File_replication_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_replication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_replication_proto_goTypes,
		DependencyIndexes: file_replication_proto_depIdxs,
		MessageInfos:      file_replication_proto_msgTypes,
	}.Build()
	File_replication_proto = out.File
	file_replication_proto_rawDesc = nil
	file_replication_proto_goTypes = nil
	file_replication_proto_depIdxs = nil
}

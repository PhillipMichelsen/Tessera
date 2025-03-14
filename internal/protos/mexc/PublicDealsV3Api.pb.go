// spot@public.deals.v3.api.pb

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.2
// source: PublicDealsV3Api.proto

package protos

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

type PublicDealsV3Api struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deals     []*PublicDealsV3ApiItem `protobuf:"bytes,1,rep,name=deals,proto3" json:"deals,omitempty"`
	EventType string                  `protobuf:"bytes,2,opt,name=eventType,proto3" json:"eventType,omitempty"`
}

func (x *PublicDealsV3Api) Reset() {
	*x = PublicDealsV3Api{}
	mi := &file_PublicDealsV3Api_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PublicDealsV3Api) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicDealsV3Api) ProtoMessage() {}

func (x *PublicDealsV3Api) ProtoReflect() protoreflect.Message {
	mi := &file_PublicDealsV3Api_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicDealsV3Api.ProtoReflect.Descriptor instead.
func (*PublicDealsV3Api) Descriptor() ([]byte, []int) {
	return file_PublicDealsV3Api_proto_rawDescGZIP(), []int{0}
}

func (x *PublicDealsV3Api) GetDeals() []*PublicDealsV3ApiItem {
	if x != nil {
		return x.Deals
	}
	return nil
}

func (x *PublicDealsV3Api) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

type PublicDealsV3ApiItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Price     string `protobuf:"bytes,1,opt,name=price,proto3" json:"price,omitempty"`
	Quantity  string `protobuf:"bytes,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	TradeType int32  `protobuf:"varint,3,opt,name=tradeType,proto3" json:"tradeType,omitempty"`
	Time      int64  `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *PublicDealsV3ApiItem) Reset() {
	*x = PublicDealsV3ApiItem{}
	mi := &file_PublicDealsV3Api_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PublicDealsV3ApiItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublicDealsV3ApiItem) ProtoMessage() {}

func (x *PublicDealsV3ApiItem) ProtoReflect() protoreflect.Message {
	mi := &file_PublicDealsV3Api_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublicDealsV3ApiItem.ProtoReflect.Descriptor instead.
func (*PublicDealsV3ApiItem) Descriptor() ([]byte, []int) {
	return file_PublicDealsV3Api_proto_rawDescGZIP(), []int{1}
}

func (x *PublicDealsV3ApiItem) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *PublicDealsV3ApiItem) GetQuantity() string {
	if x != nil {
		return x.Quantity
	}
	return ""
}

func (x *PublicDealsV3ApiItem) GetTradeType() int32 {
	if x != nil {
		return x.TradeType
	}
	return 0
}

func (x *PublicDealsV3ApiItem) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_PublicDealsV3Api_proto protoreflect.FileDescriptor

var file_PublicDealsV3Api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x44, 0x65, 0x61, 0x6c, 0x73, 0x56, 0x33, 0x41,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x10, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x44, 0x65, 0x61, 0x6c, 0x73, 0x56, 0x33, 0x41, 0x70, 0x69, 0x12, 0x2b, 0x0a, 0x05,
	0x64, 0x65, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x44, 0x65, 0x61, 0x6c, 0x73, 0x56, 0x33, 0x41, 0x70, 0x69, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x05, 0x64, 0x65, 0x61, 0x6c, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0x7a, 0x0a, 0x14, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x44, 0x65, 0x61, 0x6c, 0x73, 0x56, 0x33, 0x41, 0x70, 0x69, 0x49, 0x74, 0x65, 0x6d, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x72, 0x61, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x42, 0x39, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x78, 0x63, 0x2e, 0x70,
	0x75, 0x73, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x42, 0x15, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x44, 0x65, 0x61, 0x6c, 0x73,
	0x56, 0x33, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x01, 0x50, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_PublicDealsV3Api_proto_rawDescOnce sync.Once
	file_PublicDealsV3Api_proto_rawDescData = file_PublicDealsV3Api_proto_rawDesc
)

func file_PublicDealsV3Api_proto_rawDescGZIP() []byte {
	file_PublicDealsV3Api_proto_rawDescOnce.Do(func() {
		file_PublicDealsV3Api_proto_rawDescData = protoimpl.X.CompressGZIP(file_PublicDealsV3Api_proto_rawDescData)
	})
	return file_PublicDealsV3Api_proto_rawDescData
}

var file_PublicDealsV3Api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_PublicDealsV3Api_proto_goTypes = []any{
	(*PublicDealsV3Api)(nil),     // 0: PublicDealsV3Api
	(*PublicDealsV3ApiItem)(nil), // 1: PublicDealsV3ApiItem
}
var file_PublicDealsV3Api_proto_depIdxs = []int32{
	1, // 0: PublicDealsV3Api.deals:type_name -> PublicDealsV3ApiItem
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_PublicDealsV3Api_proto_init() }
func file_PublicDealsV3Api_proto_init() {
	if File_PublicDealsV3Api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_PublicDealsV3Api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_PublicDealsV3Api_proto_goTypes,
		DependencyIndexes: file_PublicDealsV3Api_proto_depIdxs,
		MessageInfos:      file_PublicDealsV3Api_proto_msgTypes,
	}.Build()
	File_PublicDealsV3Api_proto = out.File
	file_PublicDealsV3Api_proto_rawDesc = nil
	file_PublicDealsV3Api_proto_goTypes = nil
	file_PublicDealsV3Api_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: market_data.proto

package generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OHLCV struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Open          float64                `protobuf:"fixed64,1,opt,name=open,proto3" json:"open,omitempty"`
	High          float64                `protobuf:"fixed64,2,opt,name=high,proto3" json:"high,omitempty"`
	Low           float64                `protobuf:"fixed64,3,opt,name=low,proto3" json:"low,omitempty"`
	Close         float64                `protobuf:"fixed64,4,opt,name=close,proto3" json:"close,omitempty"`
	Volume        float64                `protobuf:"fixed64,5,opt,name=volume,proto3" json:"volume,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OHLCV) Reset() {
	*x = OHLCV{}
	mi := &file_market_data_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OHLCV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OHLCV) ProtoMessage() {}

func (x *OHLCV) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OHLCV.ProtoReflect.Descriptor instead.
func (*OHLCV) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{0}
}

func (x *OHLCV) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *OHLCV) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *OHLCV) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *OHLCV) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *OHLCV) GetVolume() float64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *OHLCV) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type Trade struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Price         float64                `protobuf:"fixed64,1,opt,name=price,proto3" json:"price,omitempty"`
	Quantity      float64                `protobuf:"fixed64,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Trade) Reset() {
	*x = Trade{}
	mi := &file_market_data_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trade) ProtoMessage() {}

func (x *Trade) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trade.ProtoReflect.Descriptor instead.
func (*Trade) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{1}
}

func (x *Trade) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Trade) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Trade) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type BookTicker struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BidPrice      float64                `protobuf:"fixed64,2,opt,name=bid_price,json=bidPrice,proto3" json:"bid_price,omitempty"`
	BidQuantity   float64                `protobuf:"fixed64,3,opt,name=bid_quantity,json=bidQuantity,proto3" json:"bid_quantity,omitempty"`
	AskPrice      float64                `protobuf:"fixed64,4,opt,name=ask_price,json=askPrice,proto3" json:"ask_price,omitempty"`
	AskQuantity   float64                `protobuf:"fixed64,5,opt,name=ask_quantity,json=askQuantity,proto3" json:"ask_quantity,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BookTicker) Reset() {
	*x = BookTicker{}
	mi := &file_market_data_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BookTicker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BookTicker) ProtoMessage() {}

func (x *BookTicker) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BookTicker.ProtoReflect.Descriptor instead.
func (*BookTicker) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{2}
}

func (x *BookTicker) GetBidPrice() float64 {
	if x != nil {
		return x.BidPrice
	}
	return 0
}

func (x *BookTicker) GetBidQuantity() float64 {
	if x != nil {
		return x.BidQuantity
	}
	return 0
}

func (x *BookTicker) GetAskPrice() float64 {
	if x != nil {
		return x.AskPrice
	}
	return 0
}

func (x *BookTicker) GetAskQuantity() float64 {
	if x != nil {
		return x.AskQuantity
	}
	return 0
}

func (x *BookTicker) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type OrderBookEntry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Price         float64                `protobuf:"fixed64,1,opt,name=price,proto3" json:"price,omitempty"`
	Quantity      float64                `protobuf:"fixed64,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderBookEntry) Reset() {
	*x = OrderBookEntry{}
	mi := &file_market_data_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderBookEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderBookEntry) ProtoMessage() {}

func (x *OrderBookEntry) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderBookEntry.ProtoReflect.Descriptor instead.
func (*OrderBookEntry) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{3}
}

func (x *OrderBookEntry) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *OrderBookEntry) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type OrderBookUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AskUpdates    []*OrderBookEntry      `protobuf:"bytes,1,rep,name=ask_updates,json=askUpdates,proto3" json:"ask_updates,omitempty"`
	BidUpdates    []*OrderBookEntry      `protobuf:"bytes,2,rep,name=bid_updates,json=bidUpdates,proto3" json:"bid_updates,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderBookUpdate) Reset() {
	*x = OrderBookUpdate{}
	mi := &file_market_data_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderBookUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderBookUpdate) ProtoMessage() {}

func (x *OrderBookUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderBookUpdate.ProtoReflect.Descriptor instead.
func (*OrderBookUpdate) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{4}
}

func (x *OrderBookUpdate) GetAskUpdates() []*OrderBookEntry {
	if x != nil {
		return x.AskUpdates
	}
	return nil
}

func (x *OrderBookUpdate) GetBidUpdates() []*OrderBookEntry {
	if x != nil {
		return x.BidUpdates
	}
	return nil
}

func (x *OrderBookUpdate) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type OrderBookSnapshot struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Asks          []*OrderBookEntry      `protobuf:"bytes,1,rep,name=asks,proto3" json:"asks,omitempty"`
	Bids          []*OrderBookEntry      `protobuf:"bytes,2,rep,name=bids,proto3" json:"bids,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderBookSnapshot) Reset() {
	*x = OrderBookSnapshot{}
	mi := &file_market_data_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderBookSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderBookSnapshot) ProtoMessage() {}

func (x *OrderBookSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderBookSnapshot.ProtoReflect.Descriptor instead.
func (*OrderBookSnapshot) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{5}
}

func (x *OrderBookSnapshot) GetAsks() []*OrderBookEntry {
	if x != nil {
		return x.Asks
	}
	return nil
}

func (x *OrderBookSnapshot) GetBids() []*OrderBookEntry {
	if x != nil {
		return x.Bids
	}
	return nil
}

func (x *OrderBookSnapshot) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type SerializedJSON struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Json          string                 `protobuf:"bytes,1,opt,name=json,proto3" json:"json,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SerializedJSON) Reset() {
	*x = SerializedJSON{}
	mi := &file_market_data_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SerializedJSON) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SerializedJSON) ProtoMessage() {}

func (x *SerializedJSON) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SerializedJSON.ProtoReflect.Descriptor instead.
func (*SerializedJSON) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{6}
}

func (x *SerializedJSON) GetJson() string {
	if x != nil {
		return x.Json
	}
	return ""
}

type MarketData struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Payload:
	//
	//	*MarketData_Ohlcv
	//	*MarketData_Trade
	//	*MarketData_BookTicker
	//	*MarketData_OrderBookUpdate
	//	*MarketData_OrderBookSnapshot
	//	*MarketData_SerializedJson
	Payload       isMarketData_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MarketData) Reset() {
	*x = MarketData{}
	mi := &file_market_data_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MarketData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarketData) ProtoMessage() {}

func (x *MarketData) ProtoReflect() protoreflect.Message {
	mi := &file_market_data_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarketData.ProtoReflect.Descriptor instead.
func (*MarketData) Descriptor() ([]byte, []int) {
	return file_market_data_proto_rawDescGZIP(), []int{7}
}

func (x *MarketData) GetPayload() isMarketData_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *MarketData) GetOhlcv() *OHLCV {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_Ohlcv); ok {
			return x.Ohlcv
		}
	}
	return nil
}

func (x *MarketData) GetTrade() *Trade {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_Trade); ok {
			return x.Trade
		}
	}
	return nil
}

func (x *MarketData) GetBookTicker() *BookTicker {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_BookTicker); ok {
			return x.BookTicker
		}
	}
	return nil
}

func (x *MarketData) GetOrderBookUpdate() *OrderBookUpdate {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_OrderBookUpdate); ok {
			return x.OrderBookUpdate
		}
	}
	return nil
}

func (x *MarketData) GetOrderBookSnapshot() *OrderBookSnapshot {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_OrderBookSnapshot); ok {
			return x.OrderBookSnapshot
		}
	}
	return nil
}

func (x *MarketData) GetSerializedJson() *SerializedJSON {
	if x != nil {
		if x, ok := x.Payload.(*MarketData_SerializedJson); ok {
			return x.SerializedJson
		}
	}
	return nil
}

type isMarketData_Payload interface {
	isMarketData_Payload()
}

type MarketData_Ohlcv struct {
	Ohlcv *OHLCV `protobuf:"bytes,1,opt,name=ohlcv,proto3,oneof"`
}

type MarketData_Trade struct {
	Trade *Trade `protobuf:"bytes,2,opt,name=trade,proto3,oneof"`
}

type MarketData_BookTicker struct {
	BookTicker *BookTicker `protobuf:"bytes,3,opt,name=book_ticker,json=bookTicker,proto3,oneof"`
}

type MarketData_OrderBookUpdate struct {
	OrderBookUpdate *OrderBookUpdate `protobuf:"bytes,4,opt,name=order_book_update,json=orderBookUpdate,proto3,oneof"`
}

type MarketData_OrderBookSnapshot struct {
	OrderBookSnapshot *OrderBookSnapshot `protobuf:"bytes,5,opt,name=order_book_snapshot,json=orderBookSnapshot,proto3,oneof"`
}

type MarketData_SerializedJson struct {
	SerializedJson *SerializedJSON `protobuf:"bytes,6,opt,name=serialized_json,json=serializedJson,proto3,oneof"`
}

func (*MarketData_Ohlcv) isMarketData_Payload() {}

func (*MarketData_Trade) isMarketData_Payload() {}

func (*MarketData_BookTicker) isMarketData_Payload() {}

func (*MarketData_OrderBookUpdate) isMarketData_Payload() {}

func (*MarketData_OrderBookSnapshot) isMarketData_Payload() {}

func (*MarketData_SerializedJson) isMarketData_Payload() {}

var File_market_data_proto protoreflect.FileDescriptor

var file_market_data_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x01, 0x0a,
	0x05, 0x4f, 0x48, 0x4c, 0x43, 0x56, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69,
	0x67, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x10,
	0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x77,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x38,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x73, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x64,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xc6, 0x01,
	0x0a, 0x0a, 0x42, 0x6f, 0x6f, 0x6b, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09,
	0x62, 0x69, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x08, 0x62, 0x69, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x69, 0x64,
	0x5f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0b, 0x62, 0x69, 0x64, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x73, 0x6b, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x08, 0x61, 0x73, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x73, 0x6b,
	0x5f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0b, 0x61, 0x73, 0x6b, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x42, 0x0a, 0x0e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42,
	0x6f, 0x6f, 0x6b, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0xbd, 0x01, 0x0a, 0x0f, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x37,
	0x0a, 0x0b, 0x61, 0x73, 0x6b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x73, 0x6b,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x0b, 0x62, 0x69, 0x64, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x62, 0x69, 0x64, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xa5, 0x01, 0x0a, 0x11, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x12, 0x2a, 0x0a, 0x04, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f,
	0x6b, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x2a, 0x0a, 0x04,
	0x62, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x04, 0x62, 0x69, 0x64, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x22, 0x24, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64,
	0x4a, 0x53, 0x4f, 0x4e, 0x12, 0x12, 0x0a, 0x04, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6a, 0x73, 0x6f, 0x6e, 0x22, 0xf3, 0x02, 0x0a, 0x0a, 0x4d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x25, 0x0a, 0x05, 0x6f, 0x68, 0x6c, 0x63, 0x76,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x4f, 0x48, 0x4c, 0x43, 0x56, 0x48, 0x00, 0x52, 0x05, 0x6f, 0x68, 0x6c, 0x63, 0x76, 0x12, 0x25,
	0x0a, 0x05, 0x74, 0x72, 0x61, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x48, 0x00, 0x52, 0x05,
	0x74, 0x72, 0x61, 0x64, 0x65, 0x12, 0x35, 0x0a, 0x0b, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x69,
	0x63, 0x6b, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x42, 0x6f, 0x6f, 0x6b, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x48, 0x00,
	0x52, 0x0a, 0x62, 0x6f, 0x6f, 0x6b, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x11,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x48, 0x00, 0x52, 0x0f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x4b, 0x0a, 0x13, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x6f, 0x6f,
	0x6b, 0x5f, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x42,
	0x6f, 0x6f, 0x6b, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x48, 0x00, 0x52, 0x11, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x42, 0x6f, 0x6f, 0x6b, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x12, 0x41, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x5f, 0x6a,
	0x73, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4a, 0x53, 0x4f,
	0x4e, 0x48, 0x00, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4a,
	0x73, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x3c,
	0x5a, 0x3a, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x69, 0x63, 0x54, 0x72, 0x61,
	0x64, 0x65, 0x72, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_market_data_proto_rawDescOnce sync.Once
	file_market_data_proto_rawDescData = file_market_data_proto_rawDesc
)

func file_market_data_proto_rawDescGZIP() []byte {
	file_market_data_proto_rawDescOnce.Do(func() {
		file_market_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_market_data_proto_rawDescData)
	})
	return file_market_data_proto_rawDescData
}

var file_market_data_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_market_data_proto_goTypes = []any{
	(*OHLCV)(nil),                 // 0: models.OHLCV
	(*Trade)(nil),                 // 1: models.Trade
	(*BookTicker)(nil),            // 2: models.BookTicker
	(*OrderBookEntry)(nil),        // 3: models.OrderBookEntry
	(*OrderBookUpdate)(nil),       // 4: models.OrderBookUpdate
	(*OrderBookSnapshot)(nil),     // 5: models.OrderBookSnapshot
	(*SerializedJSON)(nil),        // 6: models.SerializedJSON
	(*MarketData)(nil),            // 7: models.MarketData
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_market_data_proto_depIdxs = []int32{
	8,  // 0: models.OHLCV.timestamp:type_name -> google.protobuf.Timestamp
	8,  // 1: models.Trade.timestamp:type_name -> google.protobuf.Timestamp
	8,  // 2: models.BookTicker.timestamp:type_name -> google.protobuf.Timestamp
	3,  // 3: models.OrderBookUpdate.ask_updates:type_name -> models.OrderBookEntry
	3,  // 4: models.OrderBookUpdate.bid_updates:type_name -> models.OrderBookEntry
	8,  // 5: models.OrderBookUpdate.timestamp:type_name -> google.protobuf.Timestamp
	3,  // 6: models.OrderBookSnapshot.asks:type_name -> models.OrderBookEntry
	3,  // 7: models.OrderBookSnapshot.bids:type_name -> models.OrderBookEntry
	8,  // 8: models.OrderBookSnapshot.timestamp:type_name -> google.protobuf.Timestamp
	0,  // 9: models.MarketData.ohlcv:type_name -> models.OHLCV
	1,  // 10: models.MarketData.trade:type_name -> models.Trade
	2,  // 11: models.MarketData.book_ticker:type_name -> models.BookTicker
	4,  // 12: models.MarketData.order_book_update:type_name -> models.OrderBookUpdate
	5,  // 13: models.MarketData.order_book_snapshot:type_name -> models.OrderBookSnapshot
	6,  // 14: models.MarketData.serialized_json:type_name -> models.SerializedJSON
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_market_data_proto_init() }
func file_market_data_proto_init() {
	if File_market_data_proto != nil {
		return
	}
	file_market_data_proto_msgTypes[7].OneofWrappers = []any{
		(*MarketData_Ohlcv)(nil),
		(*MarketData_Trade)(nil),
		(*MarketData_BookTicker)(nil),
		(*MarketData_OrderBookUpdate)(nil),
		(*MarketData_OrderBookSnapshot)(nil),
		(*MarketData_SerializedJson)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_market_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_market_data_proto_goTypes,
		DependencyIndexes: file_market_data_proto_depIdxs,
		MessageInfos:      file_market_data_proto_msgTypes,
	}.Build()
	File_market_data_proto = out.File
	file_market_data_proto_rawDesc = nil
	file_market_data_proto_goTypes = nil
	file_market_data_proto_depIdxs = nil
}

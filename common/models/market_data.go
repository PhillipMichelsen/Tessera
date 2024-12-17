package models

import (
	"AlgorithimcTraderDistributed/common/models/proto/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// MarketData General data type
type MarketData interface{}

type OHLCV struct {
	Open      float64   `json:"O"`
	High      float64   `json:"H"`
	Low       float64   `json:"L"`
	Close     float64   `json:"C"`
	Volume    float64   `json:"V"`
	Timestamp time.Time `json:"T"`
}

type Trade struct {
	Price     float64   `json:"P"`
	Quantity  float64   `json:"Q"`
	Timestamp time.Time `json:"T"`
}

type BookTicker struct {
	BidPrice    float64   `json:"b"`
	BidQuantity float64   `json:"B"`
	AskPrice    float64   `json:"a"`
	AskQuantity float64   `json:"A"`
	Timestamp   time.Time `json:"T"`
}

type OrderBookEntry struct {
	Price    float64 `json:"P"`
	Quantity float64 `json:"Q"`
}

type OrderBookUpdate struct {
	AskUpdates []OrderBookEntry `json:"A"`
	BidUpdates []OrderBookEntry `json:"B"`
	Timestamp  time.Time        `json:"T"`
}

type OrderBookSnapshot struct {
	Asks      []OrderBookEntry `json:"A"`
	Bids      []OrderBookEntry `json:"B"`
	Timestamp time.Time        `json:"T"`
}

type SerializedJSON struct {
	Data string `json:"data"`
}

// ConvertGoToProtoOHLCV converts a Go OHLCV struct to a Protobuf OHLCV message.
func ConvertGoToProtoOHLCV(goOHLCV OHLCV) *generated.OHLCV {
	return &generated.OHLCV{
		Open:      goOHLCV.Open,
		High:      goOHLCV.High,
		Low:       goOHLCV.Low,
		Close:     goOHLCV.Close,
		Volume:    goOHLCV.Volume,
		Timestamp: timestamppb.New(goOHLCV.Timestamp),
	}
}

// ConvertProtoToGoOHLCV converts a Protobuf OHLCV message to a Go OHLCV struct.
func ConvertProtoToGoOHLCV(protoOHLCV *generated.OHLCV) OHLCV {
	return OHLCV{
		Open:      protoOHLCV.Open,
		High:      protoOHLCV.High,
		Low:       protoOHLCV.Low,
		Close:     protoOHLCV.Close,
		Volume:    protoOHLCV.Volume,
		Timestamp: protoOHLCV.Timestamp.AsTime(),
	}
}

// ConvertGoToProtoTrade converts a Go Trade struct to a Protobuf Trade message.
func ConvertGoToProtoTrade(goTrade Trade) *generated.Trade {
	return &generated.Trade{
		Price:     goTrade.Price,
		Quantity:  goTrade.Quantity,
		Timestamp: timestamppb.New(goTrade.Timestamp),
	}
}

// ConvertProtoToGoTrade converts a Protobuf Trade message to a Go Trade struct.
func ConvertProtoToGoTrade(protoTrade *generated.Trade) Trade {
	return Trade{
		Price:     protoTrade.Price,
		Quantity:  protoTrade.Quantity,
		Timestamp: protoTrade.Timestamp.AsTime(),
	}
}

func ConvertGoToProtoBookTicker(goBookTicker BookTicker) *generated.BookTicker {
	return &generated.BookTicker{
		BidPrice:    goBookTicker.BidPrice,
		BidQuantity: goBookTicker.BidQuantity,
		AskPrice:    goBookTicker.AskPrice,
		AskQuantity: goBookTicker.AskQuantity,
		Timestamp:   timestamppb.New(goBookTicker.Timestamp),
	}
}

func ConvertProtoToGoBookTicker(protoBookTicker *generated.BookTicker) BookTicker {
	return BookTicker{
		BidPrice:    protoBookTicker.BidPrice,
		BidQuantity: protoBookTicker.BidQuantity,
		AskPrice:    protoBookTicker.AskPrice,
		AskQuantity: protoBookTicker.AskQuantity,
		Timestamp:   protoBookTicker.Timestamp.AsTime(),
	}
}

// ConvertGoToProtoOrderBookEntry converts a Go OrderBookEntry struct to a Protobuf OrderBookEntry message.
func ConvertGoToProtoOrderBookEntry(goEntry OrderBookEntry) *generated.OrderBookEntry {
	return &generated.OrderBookEntry{
		Price:    goEntry.Price,
		Quantity: goEntry.Quantity,
	}
}

// ConvertProtoToGoOrderBookEntry converts a Protobuf OrderBookEntry message to a Go OrderBookEntry struct.
func ConvertProtoToGoOrderBookEntry(protoEntry *generated.OrderBookEntry) OrderBookEntry {
	return OrderBookEntry{
		Price:    protoEntry.Price,
		Quantity: protoEntry.Quantity,
	}
}

// ConvertGoToProtoOrderBookUpdate converts a Go OrderBookUpdate struct to a Protobuf OrderBookUpdate message.
func ConvertGoToProtoOrderBookUpdate(goUpdate OrderBookUpdate) *generated.OrderBookUpdate {
	protoUpdate := &generated.OrderBookUpdate{
		Timestamp: timestamppb.New(goUpdate.Timestamp),
	}

	for _, entry := range goUpdate.AskUpdates {
		protoUpdate.AskUpdates = append(protoUpdate.AskUpdates, ConvertGoToProtoOrderBookEntry(entry))
	}

	for _, entry := range goUpdate.BidUpdates {
		protoUpdate.BidUpdates = append(protoUpdate.BidUpdates, ConvertGoToProtoOrderBookEntry(entry))
	}

	return protoUpdate
}

// ConvertProtoToGoOrderBookUpdate converts a Protobuf OrderBookUpdate message to a Go OrderBookUpdate struct.
func ConvertProtoToGoOrderBookUpdate(protoUpdate *generated.OrderBookUpdate) OrderBookUpdate {
	goUpdate := OrderBookUpdate{
		Timestamp: protoUpdate.Timestamp.AsTime(),
	}

	for _, entry := range protoUpdate.AskUpdates {
		goUpdate.AskUpdates = append(goUpdate.AskUpdates, ConvertProtoToGoOrderBookEntry(entry))
	}

	for _, entry := range protoUpdate.BidUpdates {
		goUpdate.BidUpdates = append(goUpdate.BidUpdates, ConvertProtoToGoOrderBookEntry(entry))
	}

	return goUpdate
}

// ConvertGoToProtoOrderBookSnapshot converts a Go OrderBookSnapshot struct to a Protobuf OrderBookSnapshot message.
func ConvertGoToProtoOrderBookSnapshot(goSnapshot OrderBookSnapshot) *generated.OrderBookSnapshot {
	protoSnapshot := &generated.OrderBookSnapshot{
		Timestamp: timestamppb.New(goSnapshot.Timestamp),
	}

	for _, entry := range goSnapshot.Asks {
		protoSnapshot.Asks = append(protoSnapshot.Asks, ConvertGoToProtoOrderBookEntry(entry))
	}

	for _, entry := range goSnapshot.Bids {
		protoSnapshot.Bids = append(protoSnapshot.Bids, ConvertGoToProtoOrderBookEntry(entry))
	}

	return protoSnapshot
}

// ConvertProtoToGoOrderBookSnapshot converts a Protobuf OrderBookSnapshot message to a Go OrderBookSnapshot struct.
func ConvertProtoToGoOrderBookSnapshot(protoSnapshot *generated.OrderBookSnapshot) OrderBookSnapshot {
	goSnapshot := OrderBookSnapshot{
		Timestamp: protoSnapshot.Timestamp.AsTime(),
	}

	for _, entry := range protoSnapshot.Asks {
		goSnapshot.Asks = append(goSnapshot.Asks, ConvertProtoToGoOrderBookEntry(entry))
	}

	for _, entry := range protoSnapshot.Bids {
		goSnapshot.Bids = append(goSnapshot.Bids, ConvertProtoToGoOrderBookEntry(entry))
	}

	return goSnapshot
}

func ConvertGoToProtoSerializedJSON(goSerializedJSON SerializedJSON) *generated.SerializedJSON {
	protoSerializedJSON := &generated.SerializedJSON{
		Data: goSerializedJSON.Data,
	}

	return protoSerializedJSON
}

func ConvertProtoToGoSerializedJSON(protoSerializedJson *generated.SerializedJSON) SerializedJSON {
	goSerializedJSON := SerializedJSON{
		Data: protoSerializedJson.Data,
	}

	return goSerializedJSON
}

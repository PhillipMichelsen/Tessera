package models

import (
	"AlgorithimcTraderDistributed/common/models/proto/generated"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// MarketDataPiece represents a piece of market data from a source at a specific point in time.
type MarketDataPiece struct {
	Source            string     `json:"source"`
	Symbol            string     `json:"symbol"`
	BaseType          string     `json:"base_type"`
	Interval          string     `json:"interval,omitempty"`
	Processors        string     `json:"processors"`
	ExternalTimestamp time.Time  `json:"base_timestamp"`
	ReceiveTimestamp  time.Time  `json:"receive_timestamp"`
	SendTimestamp     time.Time  `json:"send_timestamp"`
	Payload           MarketData `json:"payload"`
}

// RoutingPattern and RoutingRule represent routing patterns and their corresponding rules.
type RoutingPattern struct {
	Source     string
	Symbol     string
	BaseType   string
	Interval   string
	Processors string
}

func RoutingPatternToAMQPTable(pattern RoutingPattern) amqp.Table {
	table := amqp.Table{}

	if pattern.Source != "*" {
		table["source"] = pattern.Source
	}
	if pattern.Symbol != "*" {
		table["symbol"] = pattern.Symbol
	}
	if pattern.BaseType != "*" {
		table["base_type"] = pattern.BaseType
	}
	if pattern.Interval != "*" {
		table["interval"] = pattern.Interval
	}
	if pattern.Processors != "*" {
		table["processors"] = pattern.Processors
	}

	table["x-match"] = "all"

	return table
}

func AMQPTableToRoutingPattern(table amqp.Table) RoutingPattern {
	return RoutingPattern{
		Source:     table["source"].(string),
		Symbol:     table["symbol"].(string),
		BaseType:   table["base_type"].(string),
		Interval:   table["interval"].(string),
		Processors: table["processors"].(string),
	}
}

func MarketDataPieceToRoutingPattern(data MarketDataPiece) RoutingPattern {
	return RoutingPattern{
		Source:     data.Source,
		Symbol:     data.Symbol,
		BaseType:   data.BaseType,
		Interval:   data.Interval,
		Processors: data.Processors,
	}
}

// ConvertGoToProtoMarketDataPiece converts a Go MarketDataPiece struct to a Protobuf MarketDataPiece message.
func ConvertGoToProtoMarketDataPiece(goData MarketDataPiece) *generated.MarketDataPiece {
	protoData := &generated.MarketDataPiece{
		Source:            goData.Source,
		Symbol:            goData.Symbol,
		BaseType:          goData.BaseType,
		Interval:          goData.Interval,
		Processors:        goData.Processors,
		ExternalTimestamp: timestamppb.New(goData.ExternalTimestamp),
		ReceiveTimestamp:  timestamppb.New(goData.ReceiveTimestamp),
		SendTimestamp:     timestamppb.New(goData.SendTimestamp),
	}

	// Populate the appropriate payload type in the oneof field
	switch payload := goData.Payload.(type) {
	case OHLCV:
		protoData.Payload = &generated.MarketDataPiece_Ohlcv{Ohlcv: ConvertGoToProtoOHLCV(payload)}
	case Trade:
		protoData.Payload = &generated.MarketDataPiece_Trade{Trade: ConvertGoToProtoTrade(payload)}
	case BookTicker:
		protoData.Payload = &generated.MarketDataPiece_BookTicker{BookTicker: ConvertGoToProtoBookTicker(payload)}
	case OrderBookSnapshot:
		protoData.Payload = &generated.MarketDataPiece_OrderBookSnapshot{OrderBookSnapshot: ConvertGoToProtoOrderBookSnapshot(payload)}
	case SerializedJSON:
		protoData.Payload = &generated.MarketDataPiece_SerializedJson{SerializedJson: ConvertGoToProtoSerializedJSON(payload)}
	}

	return protoData
}

// ConvertProtoToGoMarketDataPiece converts a Protobuf MarketDataPiece message to a Go MarketDataPiece struct.
func ConvertProtoToGoMarketDataPiece(protoData *generated.MarketDataPiece) MarketDataPiece {
	goData := MarketDataPiece{
		Source:            protoData.Source,
		Symbol:            protoData.Symbol,
		BaseType:          protoData.BaseType,
		Interval:          protoData.Interval,
		Processors:        protoData.Processors,
		ExternalTimestamp: protoData.ExternalTimestamp.AsTime(),
		ReceiveTimestamp:  protoData.ReceiveTimestamp.AsTime(),
		SendTimestamp:     protoData.SendTimestamp.AsTime(),
	}

	// Populate the appropriate payload type in the oneof field
	switch payload := protoData.Payload.(type) {
	case *generated.MarketDataPiece_Ohlcv:
		goData.Payload = ConvertProtoToGoOHLCV(payload.Ohlcv)
	case *generated.MarketDataPiece_Trade:
		goData.Payload = ConvertProtoToGoTrade(payload.Trade)
	case *generated.MarketDataPiece_BookTicker:
		goData.Payload = ConvertProtoToGoBookTicker(payload.BookTicker)
	case *generated.MarketDataPiece_OrderBookSnapshot:
		goData.Payload = ConvertProtoToGoOrderBookSnapshot(payload.OrderBookSnapshot)
	case *generated.MarketDataPiece_SerializedJson:
		goData.Payload = ConvertProtoToGoSerializedJSON(payload.SerializedJson)
	}

	return goData
}

func ConvertProtoBytesToGoMarketDataPiece(protoBytes []byte) (MarketDataPiece, error) {
	var protoData generated.MarketDataPiece
	err := proto.Unmarshal(protoBytes, &protoData)
	if err != nil {
		return MarketDataPiece{}, err
	}

	return ConvertProtoToGoMarketDataPiece(&protoData), nil
}

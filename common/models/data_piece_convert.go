package models

import (
	"AlgorithmicTraderDistributed/common/models/proto/generated"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ConvertGoToProtoMarketDataPacket converts a Go MarketDataPacket struct to a Protobuf MarketDataPacket message.
func ConvertGoToProtoMarketDataPacket(goData MarketDataPacket) *generated.MarketDataPacket {
	protoData := &generated.MarketDataPacket{
		CurrentUuid: goData.CurrentUUID.String(),
		TransmissionPath: func() []string {
			paths := make([]string, len(goData.TransmissionPath))
			for i, path := range goData.TransmissionPath {
				paths[i] = path.String()
			}
			return paths
		}(),
		TransmissionTimestamps: func() []*timestamppb.Timestamp {
			timestamps := make([]*timestamppb.Timestamp, len(goData.TransmissionTimestamps))
			for i, ts := range goData.TransmissionTimestamps {
				timestamps[i] = timestamppb.New(ts)
			}
			return timestamps
		}(),
		Payload: ConvertGoToProtoMarketData(goData.Payload),
	}

	return protoData
}

// ConvertProtoToGoMarketDataPacket converts a Protobuf MarketDataPacket message to a Go MarketDataPacket struct.
func ConvertProtoToGoMarketDataPacket(protoData *generated.MarketDataPacket) MarketDataPacket {
	goData := MarketDataPacket{
		CurrentUUID: func() uuid.UUID {
			uuidValue, _ := uuid.Parse(protoData.CurrentUuid)
			return uuidValue
		}(),
		TransmissionPath: func() []uuid.UUID {
			paths := make([]uuid.UUID, len(protoData.TransmissionPath))
			for i, path := range protoData.TransmissionPath {
				paths[i], _ = uuid.Parse(path)
			}
			return paths
		}(),
		TransmissionTimestamps: func() []time.Time {
			timestamps := make([]time.Time, len(protoData.TransmissionTimestamps))
			for i, ts := range protoData.TransmissionTimestamps {
				timestamps[i] = ts.AsTime()
			}
			return timestamps
		}(),
		Payload: ConvertProtoToGoMarketData(protoData.Payload),
	}

	return goData
}

// ConvertProtoBytesToGoMarketDataPacket converts Protobuf-encoded bytes to a Go MarketDataPacket struct.
func ConvertProtoBytesToGoMarketDataPacket(protoBytes []byte) (MarketDataPacket, error) {
	var protoData generated.MarketDataPacket
	err := proto.Unmarshal(protoBytes, &protoData)
	if err != nil {
		return MarketDataPacket{}, err
	}

	return ConvertProtoToGoMarketDataPacket(&protoData), nil
}

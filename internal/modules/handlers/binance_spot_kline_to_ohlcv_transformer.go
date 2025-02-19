package handlers

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type BinanceSpotKlineToOHLCVTransformerConfig struct {
	InputOutputUUIDMappings map[string]string `yaml:"input_output_uuid_mappings"`
}

type BinanceSpotKlineToOHLCVTransformer struct {
	moduleUUID              uuid.UUID
	config                  BinanceSpotKlineToOHLCVTransformerConfig
	inputOutputUUIDMappings map[uuid.UUID]uuid.UUID
	stopSignalChannel       chan struct{}
}

func NewBinanceSpotKlineToOHLCVTransformer(moduleUUID uuid.UUID) *BinanceSpotKlineToOHLCVTransformer {
	return &BinanceSpotKlineToOHLCVTransformer{
		moduleUUID: moduleUUID,
	}
}

func (b *BinanceSpotKlineToOHLCVTransformer) Initialize(rawConfig map[string]interface{}) error {
	// Initialize mapping containers
	b.config.InputOutputUUIDMappings = make(map[string]string)
	b.inputOutputUUIDMappings = make(map[uuid.UUID]uuid.UUID)

	// Get and validate mappings from config
	mappings, ok := rawConfig["input_output_uuid_mappings"].(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("invalid input_output_uuid_mappings: %v", rawConfig["input_output_uuid_mappings"])
	}

	// Process each mapping
	for k, v := range mappings {
		key, okKey := k.(string)
		value, okValue := v.(string)
		if !okKey || !okValue {
			return fmt.Errorf("invalid input_output_uuid_mappings: %v", mappings)
		}

		// Store string mapping
		b.config.InputOutputUUIDMappings[key] = value

		// Parse and store UUID mapping
		keyUUID, err := uuid.Parse(key)
		if err != nil {
			return fmt.Errorf("invalid key UUID: %s, error: %v", key, err)
		}
		valueUUID, err := uuid.Parse(value)
		if err != nil {
			return fmt.Errorf("invalid value UUID: %s, error: %v", value, err)
		}
		b.inputOutputUUIDMappings[keyUUID] = valueUUID
	}

	return nil
}

func (b *BinanceSpotKlineToOHLCVTransformer) Run(instanceAPI api.InstanceAPIInternal, runtimeErrorReceiver func(error)) {
	inputChannel := make(chan interface{})
	instanceAPI.RegisterModuleInputChannel(inputChannel)
	defer instanceAPI.DeregisterModuleInputChannel()

	b.stopSignalChannel = make(chan struct{})

	for {
		select {
		case <-b.stopSignalChannel:
			return
		case message := <-inputChannel:
			packet, ok := message.(models.Packet)
			if !ok {
				runtimeErrorReceiver(fmt.Errorf("invalid message received (Expected Packet type): %v", message))
				continue
			}

			transformedPacket, err := b.handlePacket(packet)
			if err != nil {
				runtimeErrorReceiver(err)
				continue
			}

			instanceAPI.DispatchPacket(transformedPacket)
		}
	}
}

func (b *BinanceSpotKlineToOHLCVTransformer) Stop() error {
	close(b.stopSignalChannel)

	return nil
}

func (b *BinanceSpotKlineToOHLCVTransformer) handlePacket(packet models.Packet) (models.Packet, error) {
	marketData, ok := packet.Payload.(models.MarketData)
	if !ok {
		return models.Packet{}, fmt.Errorf("invalid packet payload type: %T", packet.Payload)
	}

	serializedJSON, ok := marketData.Data.(models.SerializedJSON)
	if !ok {
		return models.Packet{}, fmt.Errorf("invalid market data data type: %T", marketData.Data)
	}

	// Parse kline JSON
	jsonPayload := serializedJSON.JSON
	timestamp := gjson.Get(jsonPayload, "k.t")
	open := gjson.Get(jsonPayload, "k.o")
	high := gjson.Get(jsonPayload, "k.h")
	low := gjson.Get(jsonPayload, "k.l")
	closePrice := gjson.Get(jsonPayload, "k.c")
	volume := gjson.Get(jsonPayload, "k.v")

	if !timestamp.Exists() || !open.Exists() || !high.Exists() ||
		!low.Exists() || !closePrice.Exists() || !volume.Exists() {
		return models.Packet{}, fmt.Errorf("missing required fields in kline JSON: %s", jsonPayload)
	}

	// Create OHLCV
	ohlcv := models.OHLCV{
		Open:      open.Float(),
		High:      high.Float(),
		Low:       low.Float(),
		Close:     closePrice.Float(),
		Volume:    volume.Float(),
		Timestamp: time.UnixMilli(timestamp.Int()),
	}

	// Get destination UUID
	destinationUUID, exists := b.inputOutputUUIDMappings[packet.DestinationModuleUUID]
	if !exists {
		panic(fmt.Sprintf("no mapping for UUID: %s", packet.DestinationModuleUUID))
	}

	return models.Packet{
		SourceModuleUUID:      b.moduleUUID,
		DestinationModuleUUID: destinationUUID,
		Payload: models.MarketData{
			UUID: destinationUUID,
			Data: ohlcv,
		},
	}, nil
}

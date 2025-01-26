package concrete

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"time"
)

type BinanceSpotKlineToOHLCVTransformerConfig struct {
	InputOutputUUIDMappings map[string]string
}

type BinanceSpotKlineToOHLCVTransformer struct {
	*modules.BaseModule
	rawConfig               map[string]interface{}
	config                  BinanceSpotKlineToOHLCVTransformerConfig
	inputOutputUUIDMappings map[uuid.UUID]uuid.UUID
}

// NewBinanceSpotKlineToOHLCVTransformer creates a new instance of the transformer.
func NewBinanceSpotKlineToOHLCVTransformer(
	componentUUID uuid.UUID,
	rawConfig map[string]interface{},
	outputChannel chan interface{},
	updateAlertChannel chan uuid.UUID,
) *BinanceSpotKlineToOHLCVTransformer {
	return &BinanceSpotKlineToOHLCVTransformer{
		BaseModule: modules.NewBaseModule(componentUUID, outputChannel, updateAlertChannel),
		rawConfig:  rawConfig,
	}
}

func (b *BinanceSpotKlineToOHLCVTransformer) Initialize() {
	b.BaseModule.Initialize(func() {
		// Ensure config maps are initialized
		b.config.InputOutputUUIDMappings = make(map[string]string)
		b.inputOutputUUIDMappings = make(map[uuid.UUID]uuid.UUID)

		// Parse and validate `input_output_uuid_mappings` from rawConfig
		mappings, ok := b.rawConfig["input_output_uuid_mappings"].(map[interface{}]interface{})
		if !ok {
			panic("input_output_uuid_mappings is missing or not a valid map")
		}

		// Process and validate the mappings
		for k, v := range mappings {
			key, okKey := k.(string)
			value, okValue := v.(string)
			if !okKey || !okValue {
				panic(fmt.Sprintf("input_output_uuid_mappings key or value is not a string: key=%v, value=%v", k, v))
			}

			// Add to config.InputOutputUUIDMappings
			b.config.InputOutputUUIDMappings[key] = value

			// Parse and validate UUIDs
			keyUUID, err := uuid.Parse(key)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse key UUID: %s, error: %v", key, err))
			}
			valueUUID, err := uuid.Parse(value)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse value UUID: %s, error: %v", value, err))
			}

			// Add to inputOutputUUIDMappings
			b.inputOutputUUIDMappings[keyUUID] = valueUUID
		}
	})
}

func (b *BinanceSpotKlineToOHLCVTransformer) Start() {
	b.BaseModule.Start(b.runWorker)
}

// runWorker processes incoming packets, transforms them, and pushes them to the output channel.
func (b *BinanceSpotKlineToOHLCVTransformer) runWorker(workerChannels modules.WorkerChannels) {
	for {
		select {
		case <-workerChannels.StopSignal:
			return
		case packet := <-workerChannels.InputChannel:
			marketDataPacket, ok := packet.(models.MarketDataPacket)
			if !ok {
				panic(fmt.Sprintf("Received unexpected packet type: %T", packet))
			}
			workerChannels.OutputChannel <- b.handlePacket(marketDataPacket)
		}
	}
}

func (b *BinanceSpotKlineToOHLCVTransformer) handlePacket(packet models.MarketDataPacket) models.MarketDataPacket {
	timeNow := time.Now()
	defer func() {
		// Print the time taken to process the packet
		fmt.Printf("Processed packet with Packet UUID %s at %s\n", packet.CurrentUUID, time.Now().Sub(timeNow))
	}()

	serializedJSONPayload, ok := packet.Payload.(models.SerializedJSON)
	if !ok {
		panic(fmt.Sprintf("'Payload in packet with Packet UUID %s is not of type SerializedJSON.'", packet.CurrentUUID))
	}

	// Extract and validate required fields from the JSON payload
	jsonPayload := serializedJSONPayload.JSON
	timestamp := gjson.Get(jsonPayload, "k.t")
	open := gjson.Get(jsonPayload, "k.o")
	high := gjson.Get(jsonPayload, "k.h")
	low := gjson.Get(jsonPayload, "k.l")
	closePrice := gjson.Get(jsonPayload, "k.c")
	volume := gjson.Get(jsonPayload, "k.v")

	// Check if required fields exist
	if !timestamp.Exists() || !open.Exists() || !high.Exists() || !low.Exists() || !closePrice.Exists() || !volume.Exists() {
		panic(fmt.Sprintf("Invalid JSON: Missing kline values in packet with Packet UUID %s", packet.CurrentUUID))
	}

	// Convert extracted fields
	ohlcv := models.OHLCV{
		Open:      open.Float(),
		High:      high.Float(),
		Low:       low.Float(),
		Close:     closePrice.Float(),
		Volume:    volume.Float(),
		Timestamp: time.UnixMilli(timestamp.Int()),
	}

	newCurrentUUID, exists := b.inputOutputUUIDMappings[packet.CurrentUUID]
	if !exists {
		panic(fmt.Sprintf("'No output mapping found for packet with Packet UUID %s'", packet.CurrentUUID))
	}

	return models.MarketDataPacket{
		CurrentUUID:           newCurrentUUID,
		UUIDHistory:           append(packet.UUIDHistory, packet.CurrentUUID),
		UUIDHistoryTimestamps: append(packet.UUIDHistoryTimestamps, time.Now()),
		Payload:               ohlcv,
	}
}

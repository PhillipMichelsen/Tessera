package workers

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

// BinanceSpotKlineToOHLCVWorkerConfig uses the worker.InputOutputMapping.
type BinanceSpotKlineToOHLCVWorkerConfig struct {
	InputOutputMapping worker.InputOutputMapping
	InputMailboxUUID   uuid.UUID
}

// BinanceSpotKlineToOHLCVWorker implements the worker.Worker interface.
type BinanceSpotKlineToOHLCVWorker struct {
	moduleUUID uuid.UUID
	config     BinanceSpotKlineToOHLCVWorkerConfig

	services worker.Services
}

// NewBinanceSpotKlineToOHLCVWorker creates a new instance.
func NewBinanceSpotKlineToOHLCVWorker(moduleUUID uuid.UUID) *BinanceSpotKlineToOHLCVWorker {
	return &BinanceSpotKlineToOHLCVWorker{
		moduleUUID: moduleUUID,
	}
}

// GetWorkerName returns the worker's name.
func (w *BinanceSpotKlineToOHLCVWorker) GetWorkerName() string {
	return "BinanceSpotKlineToOHLCVWorker"
}

// Run parses the configuration, registers a mailbox using the input mailbox UUID,
// and waits until the context is canceled.
func (w *BinanceSpotKlineToOHLCVWorker) Run(ctx context.Context, config map[string]any, services worker.Services) (worker.ExitCode, error) {
	// Parse the input_output_mapping.
	rawMapping, ok := config["input_output_mapping"].(map[string]interface{})
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("invalid input_output_mapping configuration")
	}
	mapping := make(worker.InputOutputMapping)
	for inputTag, rawOutputMetadata := range rawMapping {
		metadataMap, ok := rawOutputMetadata.(map[string]any)
		if !ok {
			return worker.RuntimeErrorExit, fmt.Errorf("invalid mapping for tag %s", inputTag)
		}
		destMailboxStr, ok := metadataMap["destination_mailbox_uuid"].(string)
		if !ok {
			return worker.RuntimeErrorExit, fmt.Errorf("missing destination_mailbox_uuid for tag %s", inputTag)
		}
		destMailbox, err := uuid.Parse(destMailboxStr)
		if err != nil {
			return worker.RuntimeErrorExit, fmt.Errorf("invalid destination_mailbox_uuid for tag %s: %v", inputTag, err)
		}
		outputTag, ok := metadataMap["tag"].(string)
		if !ok {
			return worker.RuntimeErrorExit, fmt.Errorf("missing tag for mapping with key %s", inputTag)
		}
		mapping[inputTag] = worker.OutputMetadata{
			DestinationMailboxUUID: destMailbox,
			Tag:                    outputTag,
		}
	}

	// Parse the input mailbox UUID.
	inputMailboxStr, ok := config["input_mailbox_uuid"].(string)
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("missing input_mailbox_uuid")
	}
	inputMailbox, err := uuid.Parse(inputMailboxStr)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("invalid input_mailbox_uuid: %v", err)
	}

	// Save the parsed configuration and services.
	w.config = BinanceSpotKlineToOHLCVWorkerConfig{
		InputOutputMapping: mapping,
		InputMailboxUUID:   inputMailbox,
	}
	w.services = services

	// Register a mailbox using the provided input mailbox UUID.
	services.CreateMailbox(w.config.InputMailboxUUID, w.handleMessage)

	// Block until context cancellation.
	<-ctx.Done()

	// Clean up.
	services.RemoveMailbox(w.config.InputMailboxUUID)
	return worker.NormalExit, nil
}

// handleMessage processes an incoming worker.Message.
func (w *BinanceSpotKlineToOHLCVWorker) handleMessage(message worker.Message) {
	// Look up the mapping using the incoming message's tag.
	outputMeta, exists := w.config.InputOutputMapping[message.Tag]
	if !exists {
		fmt.Printf("No mapping found for tag %s\n", message.Tag)
		return
	}

	serializedJSON, ok := message.Payload.(models.SerializedJSON)
	if !ok {
		fmt.Printf("Invalid market data content for message with tag %s\n", message.Tag)
		return
	}

	jsonPayload := serializedJSON.JSON
	timestamp := gjson.Get(jsonPayload, "k.t")
	open := gjson.Get(jsonPayload, "k.o")
	high := gjson.Get(jsonPayload, "k.h")
	low := gjson.Get(jsonPayload, "k.l")
	closePrice := gjson.Get(jsonPayload, "k.c")
	volume := gjson.Get(jsonPayload, "k.v")
	if !timestamp.Exists() || !open.Exists() || !high.Exists() ||
		!low.Exists() || !closePrice.Exists() || !volume.Exists() {
		fmt.Printf("Missing required fields in JSON payload for tag %s: %s\n", message.Tag, jsonPayload)
		return
	}

	ohlcv := models.OHLCV{
		Open:      open.Float(),
		High:      high.Float(),
		Low:       low.Float(),
		Close:     closePrice.Float(),
		Volume:    volume.Float(),
		Timestamp: time.UnixMilli(timestamp.Int()),
	}

	// Build the outgoing message with the new tag.
	newMsg := worker.Message{
		Tag:     outputMeta.Tag,
		Payload: ohlcv,
	}

	// Send the new message to the destination mailbox.
	if err := w.services.SendMessage(outputMeta.DestinationMailboxUUID, newMsg); err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}
}

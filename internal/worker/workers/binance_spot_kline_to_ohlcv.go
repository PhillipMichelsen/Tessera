package workers

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"time"
)

// BinanceSpotKlineToOHLCVConfig represents the YAML configuration for the worker.
type BinanceSpotKlineToOHLCVConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputOutputMapping map[string]struct {
		DestinationMailboxUUID uuid.UUID `yaml:"destination_mailbox_uuid"`
		Tag                    string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
}

// BinanceSpotKlineToOHLCVWorker implements the worker.Worker interface.
type BinanceSpotKlineToOHLCVWorker struct{}

func (w *BinanceSpotKlineToOHLCVWorker) Run(ctx context.Context, config any, services worker.Services) (worker.ExitCode, error) {
	configBytes, ok := config.([]byte)
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("config is not in the expected []byte format")
	}

	var cfg BinanceSpotKlineToOHLCVConfig
	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	inputChannel := make(chan worker.Message)
	services.CreateMailbox(cfg.InputMailboxUUID, worker.BuildChannelMessageReceiverForwarder(inputChannel))

	for {
		select {
		case <-ctx.Done():
			services.RemoveMailbox(cfg.InputMailboxUUID)
			return worker.NormalExit, nil
		case msg := <-inputChannel:
			// Ensure that the incoming message has a tag.
			if msg.Tag == "" {
				return worker.RuntimeErrorExit, fmt.Errorf("message is missing tag")
			}

			// Look up the destination mapping using the message tag.
			mapping, ok := cfg.InputOutputMapping[msg.Tag]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for tag: %s", msg.Tag)
			}

			ohlcv, err := w.parseJSONToOHLCV(msg.Payload.(models.SerializedJSON).JSON)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to parse JSON to OHLCV: %w", err)
			}

			if err := services.SendMessage(mapping.DestinationMailboxUUID, worker.Message{
				Tag:     mapping.Tag,
				Payload: ohlcv,
			}); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

func (w *BinanceSpotKlineToOHLCVWorker) parseJSONToOHLCV(jsonStr string) (models.OHLCV, error) {
	timestamp := gjson.Get(jsonStr, "k.t")
	open := gjson.Get(jsonStr, "k.o")
	high := gjson.Get(jsonStr, "k.h")
	low := gjson.Get(jsonStr, "k.l")
	closePrice := gjson.Get(jsonStr, "k.c")
	volume := gjson.Get(jsonStr, "k.v")
	if !timestamp.Exists() || !open.Exists() || !high.Exists() ||
		!low.Exists() || !closePrice.Exists() || !volume.Exists() {
		return models.OHLCV{}, fmt.Errorf("missing required fields in JSON payload: %s", jsonStr)
	}

	return models.OHLCV{
		Open:      open.Float(),
		High:      high.Float(),
		Low:       low.Float(),
		Close:     closePrice.Float(),
		Volume:    volume.Float(),
		Timestamp: time.UnixMilli(timestamp.Int()).UTC(),
	}, nil
}

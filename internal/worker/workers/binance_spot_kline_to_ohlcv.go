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
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// BinanceSpotKlineToOHLCVWorker implements the worker.Worker interface.
type BinanceSpotKlineToOHLCVWorker struct{}

func (w *BinanceSpotKlineToOHLCVWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	config, err := w.parseRawConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse raw config: %w", err)
	}

	services.CreateMailbox(config.InputMailboxUUID, config.InputMailboxBuffer)
	inputChannel, ok := services.GetMailboxChannel(config.InputMailboxUUID)
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to get input mailbox channel")
	}
	defer services.RemoveMailbox(config.InputMailboxUUID)

	for {
		select {
		case <-ctx.Done():
			return worker.NormalExit, nil
		case rawMessage, ok := <-inputChannel:
			// Check if the main input channel is closed.
			if !ok {
				return worker.PrematureExit, fmt.Errorf("main input channel closed")
			}

			// Cast message to worker.Message.
			message, ok := rawMessage.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message is not of type worker.Message")
			}

			// Find the output destination and new tag for the received message given its tag.
			mappedOutput, ok := config.InputOutputMapping[message.Tag]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for tag: %s", message.Tag)
			}

			ohlcv, err := w.parseJSONToOHLCV(message.Payload.(models.SerializedJSON).JSON)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to parse JSON to OHLCV: %w", err)
			}

			if err := services.SendMessage(mappedOutput.MailboxUUID, worker.Message{
				Tag:     mappedOutput.Tag,
				Payload: ohlcv,
			}, config.BlockingSend); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message blocking: %w", err)
			}
		}
	}
}

func (w *BinanceSpotKlineToOHLCVWorker) parseRawConfig(rawConfig any) (BinanceSpotKlineToOHLCVConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return BinanceSpotKlineToOHLCVConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config BinanceSpotKlineToOHLCVConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return BinanceSpotKlineToOHLCVConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return BinanceSpotKlineToOHLCVConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}

	if len(config.InputOutputMapping) == 0 {
		return BinanceSpotKlineToOHLCVConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return BinanceSpotKlineToOHLCVConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	return config, nil
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

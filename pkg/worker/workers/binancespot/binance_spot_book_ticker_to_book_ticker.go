package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/PhillipMichelsen/Tessera/pkg/models"
	"github.com/PhillipMichelsen/Tessera/pkg/worker"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

// BinanceSpotBookTickerToBookTickerConfig represents the YAML configuration for the worker.
type BinanceSpotBookTickerToBookTickerConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// BinanceSpotBookTickerToBookTickerWorker implements the worker.Worker interface.
type BinanceSpotBookTickerToBookTickerWorker struct{}

func (w *BinanceSpotBookTickerToBookTickerWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	config, err := w.parseRawConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse raw config: %w", err)
	}

	inputChannel, err := services.CreateMailbox(config.InputMailboxUUID, config.InputMailboxBuffer)
	defer services.RemoveMailbox(config.InputMailboxUUID)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to create input mailbox: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return worker.NormalExit, nil
		case rawMessage, ok := <-inputChannel:
			if !ok {
				return worker.PrematureExit, fmt.Errorf("input mailbox channel closed")
			}

			message, ok := rawMessage.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message is not of type worker.Message")
			}

			mappedOutput, ok := config.InputOutputMapping[message.Tag]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for tag: %s", message.Tag)
			}

			bookTicker, err := w.parseJSONToBookTicker(message.Payload.(models.SerializedJSON).JSON)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to parse JSON to BookTicker: %w", err)
			}

			if err := services.SendMessage(mappedOutput.MailboxUUID, worker.Message{
				Tag:     mappedOutput.Tag,
				Payload: bookTicker,
			}, config.BlockingSend); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

func (w *BinanceSpotBookTickerToBookTickerWorker) parseRawConfig(rawConfig any) (BinanceSpotBookTickerToBookTickerConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return BinanceSpotBookTickerToBookTickerConfig{}, fmt.Errorf("config is not in expected []byte format")
	}

	var config BinanceSpotBookTickerToBookTickerConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return BinanceSpotBookTickerToBookTickerConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return BinanceSpotBookTickerToBookTickerConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}
	if len(config.InputOutputMapping) == 0 {
		return BinanceSpotBookTickerToBookTickerConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return BinanceSpotBookTickerToBookTickerConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	return config, nil
}

func (w *BinanceSpotBookTickerToBookTickerWorker) parseJSONToBookTicker(jsonStr string) (models.BookTicker, error) {
	// Extract values using gjson.
	bidPrice := gjson.Get(jsonStr, "b")
	bidQuantity := gjson.Get(jsonStr, "B")
	askPrice := gjson.Get(jsonStr, "a")
	askQuantity := gjson.Get(jsonStr, "A")

	if !bidPrice.Exists() || !bidQuantity.Exists() || !askPrice.Exists() || !askQuantity.Exists() {
		return models.BookTicker{}, fmt.Errorf("missing required fields in JSON payload: %s", jsonStr)
	}

	// As the payload doesn't include a timestamp, assign current UTC time.
	return models.BookTicker{
		BidPrice:    bidPrice.Float(),
		BidQuantity: bidQuantity.Float(),
		AskPrice:    askPrice.Float(),
		AskQuantity: askQuantity.Float(),
		Timestamp:   time.Now().UTC(),
	}, nil
}

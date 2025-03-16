package workers

import (
	"Tessera/internal/models"
	"Tessera/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"time"
)

// BinanceSpotDepthUpdateToOrderBookConfig represents the YAML configuration for the worker.
type BinanceSpotDepthUpdateToOrderBookConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// BinanceSpotDepthUpdateToOrderBookWorker implements the worker.Worker interface.
type BinanceSpotDepthUpdateToOrderBookWorker struct{}

func (w *BinanceSpotDepthUpdateToOrderBookWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
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
				return worker.PrematureExit, fmt.Errorf("main input channel closed")
			}

			message, ok := rawMessage.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message is not of type worker.Message")
			}

			mappedOutput, ok := config.InputOutputMapping[message.Tag]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for tag: %s", message.Tag)
			}

			orderBookUpdate, err := w.parseJSONToOrderBookUpdate(message.Payload.(models.SerializedJSON).JSON)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to parse JSON to OrderBookUpdate: %w", err)
			}

			if err := services.SendMessage(mappedOutput.MailboxUUID, worker.Message{
				Tag:     mappedOutput.Tag,
				Payload: orderBookUpdate,
			}, config.BlockingSend); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

func (w *BinanceSpotDepthUpdateToOrderBookWorker) parseRawConfig(rawConfig any) (BinanceSpotDepthUpdateToOrderBookConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return BinanceSpotDepthUpdateToOrderBookConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config BinanceSpotDepthUpdateToOrderBookConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return BinanceSpotDepthUpdateToOrderBookConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return BinanceSpotDepthUpdateToOrderBookConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}

	if len(config.InputOutputMapping) == 0 {
		return BinanceSpotDepthUpdateToOrderBookConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return BinanceSpotDepthUpdateToOrderBookConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	return config, nil
}

func (w *BinanceSpotDepthUpdateToOrderBookWorker) parseJSONToOrderBookUpdate(jsonStr string) (models.OrderBook, error) {
	// Extract the event time (in milliseconds) from the "E" field.
	eventTimeResult := gjson.Get(jsonStr, "E")
	if !eventTimeResult.Exists() {
		return models.OrderBook{}, fmt.Errorf("missing event time in JSON payload: %s", jsonStr)
	}
	timestamp := time.UnixMilli(eventTimeResult.Int()).UTC()

	// Extract bids and asks arrays using keys "b" and "a".
	bidsResult := gjson.Get(jsonStr, "b")
	asksResult := gjson.Get(jsonStr, "a")

	if !bidsResult.Exists() || !asksResult.Exists() {
		return models.OrderBook{}, fmt.Errorf("missing required fields in JSON payload: %s", jsonStr)
	}

	var bidUpdates []models.OrderBookEntry
	var askUpdates []models.OrderBookEntry

	// Parse bids: each bid is an array [price, quantity].
	bidsResult.ForEach(func(_, bid gjson.Result) bool {
		arr := bid.Array()
		if len(arr) < 2 {
			return true
		}
		bidUpdates = append(bidUpdates, models.OrderBookEntry{
			Price:    arr[0].Float(),
			Quantity: arr[1].Float(),
		})
		return true
	})

	// Parse asks: each ask is an array [price, quantity].
	asksResult.ForEach(func(_, ask gjson.Result) bool {
		arr := ask.Array()
		if len(arr) < 2 {
			return true
		}
		askUpdates = append(askUpdates, models.OrderBookEntry{
			Price:    arr[0].Float(),
			Quantity: arr[1].Float(),
		})
		return true
	})

	return models.OrderBook{
		Bids:      bidUpdates,
		Asks:      askUpdates,
		Timestamp: timestamp,
	}, nil
}

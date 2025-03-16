package workers

import (
	"Tessera/internal/models"
	protos "Tessera/internal/protos/mexc"
	"Tessera/internal/worker"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// MEXCSpotBookTickerToBookTickerConfig represents the YAML configuration for the worker.
type MEXCSpotBookTickerToBookTickerConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// MEXCSpotBookTickerToBookTickerWorker implements the worker.Worker interface.
type MEXCSpotBookTickerToBookTickerWorker struct{}

// Run listens for incoming messages, converts the payload from protobuf to an internal BookTicker, and sends it onward.
func (w *MEXCSpotBookTickerToBookTickerWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
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

			// cast message.Payload to protos.PushDataV3ApiWrapper
			// call parseMEXCProtobufPushBodyToBookTicker
			pushData := message.Payload.(*protos.PushDataV3ApiWrapper)
			bookTicker, err := w.parseMEXCProtobufPushBodyToBookTicker(pushData)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to parse protobuf to BookTicker: %w", err)
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

// parseRawConfig unmarshals the YAML configuration.
func (w *MEXCSpotBookTickerToBookTickerWorker) parseRawConfig(rawConfig any) (MEXCSpotBookTickerToBookTickerConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return MEXCSpotBookTickerToBookTickerConfig{}, fmt.Errorf("config is not in expected []byte format")
	}

	var config MEXCSpotBookTickerToBookTickerConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return MEXCSpotBookTickerToBookTickerConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return MEXCSpotBookTickerToBookTickerConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}
	if len(config.InputOutputMapping) == 0 {
		return MEXCSpotBookTickerToBookTickerConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return MEXCSpotBookTickerToBookTickerConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	return config, nil
}

// parseMEXCProtobufPushBodyToBookTicker unmarshals the payload into a protobuf message,
// asserts its type, and maps it to an internal BookTicker.
func (w *MEXCSpotBookTickerToBookTickerWorker) parseMEXCProtobufPushBodyToBookTicker(pushData *protos.PushDataV3ApiWrapper) (models.BookTicker, error) {
	protoBookTicker := pushData.GetPublicAggreBookTicker()
	if protoBookTicker == nil {
		return models.BookTicker{}, fmt.Errorf("failed to get PublicAggreBookTicker")
	}

	bidPrice, err := strconv.ParseFloat(protoBookTicker.BidPrice, 64)
	if err != nil {
		return models.BookTicker{}, fmt.Errorf("failed to parse bidPrice %q: %w", protoBookTicker.BidPrice, err)
	}

	bidQuantity, err := strconv.ParseFloat(protoBookTicker.BidQuantity, 64)
	if err != nil {
		return models.BookTicker{}, fmt.Errorf("failed to parse bidQuantity %q: %w", protoBookTicker.BidQuantity, err)
	}

	askPrice, err := strconv.ParseFloat(protoBookTicker.AskPrice, 64)
	if err != nil {
		return models.BookTicker{}, fmt.Errorf("failed to parse askPrice %q: %w", protoBookTicker.AskPrice, err)
	}

	askQuantity, err := strconv.ParseFloat(protoBookTicker.AskQuantity, 64)
	if err != nil {
		return models.BookTicker{}, fmt.Errorf("failed to parse askQuantity %q: %w", protoBookTicker.AskQuantity, err)
	}

	bookTicker := models.BookTicker{
		BidPrice:    bidPrice,
		BidQuantity: bidQuantity,
		AskPrice:    askPrice,
		AskQuantity: askQuantity,
		Timestamp:   time.Now().UTC(),
	}

	return bookTicker, nil
}

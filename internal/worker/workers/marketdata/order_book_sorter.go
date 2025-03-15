package workers

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"sort"
)

// OrderBookSorterConfig represents the YAML configuration for the sorter worker.
type OrderBookSorterConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// OrderBookSorterWorker implements the worker.Worker interface.
type OrderBookSorterWorker struct{}

func (w *OrderBookSorterWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
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

			snapshot, ok := message.Payload.(models.OrderBook)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message payload is not of type models.OrderBook")
			}

			// Sort the order book: asks ascending, bids descending.
			sortedSnapshot := sortSnapshot(snapshot)

			err = services.SendMessage(mappedOutput.MailboxUUID, worker.Message{
				Tag:     mappedOutput.Tag,
				Payload: sortedSnapshot,
			}, config.BlockingSend)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send sorted snapshot: %w", err)
			}
		}
	}
}

func (w *OrderBookSorterWorker) parseRawConfig(rawConfig any) (OrderBookSorterConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return OrderBookSorterConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config OrderBookSorterConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return OrderBookSorterConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return OrderBookSorterConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}

	if len(config.InputOutputMapping) == 0 {
		return OrderBookSorterConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return OrderBookSorterConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	return config, nil
}

// sortSnapshot returns an order book with asks sorted in ascending order and bids sorted in descending order.
func sortSnapshot(snapshot models.OrderBook) models.OrderBook {
	// Copy and sort asks in ascending order.
	sortedAsks := make([]models.OrderBookEntry, len(snapshot.Asks))
	copy(sortedAsks, snapshot.Asks)
	sort.Slice(sortedAsks, func(i, j int) bool {
		return sortedAsks[i].Price < sortedAsks[j].Price
	})

	// Copy and sort bids in descending order.
	sortedBids := make([]models.OrderBookEntry, len(snapshot.Bids))
	copy(sortedBids, snapshot.Bids)
	sort.Slice(sortedBids, func(i, j int) bool {
		return sortedBids[i].Price > sortedBids[j].Price
	})

	return models.OrderBook{
		Asks:      sortedAsks,
		Bids:      sortedBids,
		Timestamp: snapshot.Timestamp,
	}
}

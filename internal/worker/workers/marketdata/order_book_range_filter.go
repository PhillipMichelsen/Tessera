package workers

import (
	"Tessera/internal/models"
	"Tessera/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// OrderBookRangeFilterConfig represents the YAML configuration for the worker.
type OrderBookRangeFilterConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
	InputOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"input_output_mapping"`
	BlockingSend    bool    `yaml:"blocking_send"`
	RangePercentage float64 `yaml:"range_percentage"` // n% range filter
}

// OrderBookRangeFilterWorker implements the worker.Worker interface.
type OrderBookRangeFilterWorker struct{}

func (w *OrderBookRangeFilterWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
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
			// Ensure the main input channel is open.
			if !ok {
				return worker.PrematureExit, fmt.Errorf("main input channel closed")
			}

			// Cast the incoming message.
			message, ok := rawMessage.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message is not of type worker.Message")
			}

			// Find the output destination based on the input message tag.
			mappedOutput, ok := config.InputOutputMapping[message.Tag]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for tag: %s", message.Tag)
			}

			// Cast the payload to OrderBookSnapshot.
			snapshot, ok := message.Payload.(models.OrderBook)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("message payload is not of type models.OrderBookSnapshot")
			}

			// Filter the snapshot within the specified n% range about the mid-price.
			filteredSnapshot, err := filterSnapshot(snapshot, config.RangePercentage)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to filter snapshot: %w", err)
			}

			// Send the filtered snapshot.
			err = services.SendMessage(mappedOutput.MailboxUUID, worker.Message{
				Tag:     mappedOutput.Tag,
				Payload: filteredSnapshot,
			}, config.BlockingSend)
			if err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send filtered snapshot: %w", err)
			}
		}
	}
}

func (w *OrderBookRangeFilterWorker) parseRawConfig(rawConfig any) (OrderBookRangeFilterConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return OrderBookRangeFilterConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config OrderBookRangeFilterConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return OrderBookRangeFilterConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return OrderBookRangeFilterConfig{}, fmt.Errorf("input_mailbox_uuid is required in configuration")
	}

	if len(config.InputOutputMapping) == 0 {
		return OrderBookRangeFilterConfig{}, fmt.Errorf("input_output_mapping is required in configuration")
	}

	for tag, mapping := range config.InputOutputMapping {
		if mapping.MailboxUUID == uuid.Nil {
			return OrderBookRangeFilterConfig{}, fmt.Errorf("mailbox_uuid is required for tag: %s", tag)
		}
	}

	// Ensure a valid range percentage.
	if config.RangePercentage <= 0 {
		return OrderBookRangeFilterConfig{}, fmt.Errorf("range_percentage must be greater than 0")
	}

	return config, nil
}

// filterSnapshot computes the mid-price from the best ask and bid, and returns a snapshot containing only the
// order book entries within ±rangePercentage% of that mid-price.
// The snapshot is assumed to be sorted in ascending order for asks and descending order for bids.
func filterSnapshot(snapshot models.OrderBook, rangePercentage float64) (models.OrderBook, error) {
	var bestAsk, bestBid, midPrice float64

	// Directly use the first element for best ask and bid if available.
	if len(snapshot.Asks) > 0 {
		bestAsk = snapshot.Asks[0].Price
	}
	if len(snapshot.Bids) > 0 {
		bestBid = snapshot.Bids[0].Price
	}

	// Calculate midPrice based on available sides.
	if len(snapshot.Asks) > 0 && len(snapshot.Bids) > 0 {
		midPrice = (bestAsk + bestBid) / 2.0
	} else if len(snapshot.Asks) > 0 {
		midPrice = bestAsk
	} else if len(snapshot.Bids) > 0 {
		midPrice = bestBid
	} else {
		return models.OrderBook{}, fmt.Errorf("order book contains neither asks nor bids")
	}

	lowerBound := midPrice * (1 - rangePercentage/100)
	upperBound := midPrice * (1 + rangePercentage/100)

	// Filter asks (assumed sorted in ascending order).
	filteredAsks := make([]models.OrderBookEntry, 0, len(snapshot.Asks))
	for _, ask := range snapshot.Asks {
		// Once an ask exceeds the upper bound, further asks are out-of-bound.
		if ask.Price > upperBound {
			break
		}
		if ask.Price >= lowerBound {
			filteredAsks = append(filteredAsks, ask)
		}
	}

	// Filter bids (assumed sorted in descending order).
	filteredBids := make([]models.OrderBookEntry, 0, len(snapshot.Bids))
	for _, bid := range snapshot.Bids {
		// Once a bid falls below the lower bound, further bids are out-of-bound.
		if bid.Price < lowerBound {
			break
		}
		if bid.Price <= upperBound {
			filteredBids = append(filteredBids, bid)
		}
	}

	return models.OrderBook{
		Asks:      filteredAsks,
		Bids:      filteredBids,
		Timestamp: snapshot.Timestamp,
	}, nil
}

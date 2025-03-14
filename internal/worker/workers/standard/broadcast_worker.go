package workers

import (
	"context"
	"fmt"

	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// BroadcastWorkerConfig defines the YAML configuration for mapping input tags to destination mailbox UUIDs.
type BroadcastWorkerConfig struct {
	InputMailboxUUID uuid.UUID              `yaml:"input_mailbox_uuid"`
	TagDestinations  map[string][]uuid.UUID `yaml:"tag_destinations"`
	BlockingSend     bool                   `yaml:"blocking_send"`
}

// BroadcastWorker takes a message from an input mailbox and broadcasts it to multiple destination mailboxes based on the message's tag.
type BroadcastWorker struct{}

func (w *BroadcastWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	config, err := w.parseConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse config: %w", err)
	}

	// Retrieve the input mailbox channel.
	msgCh, exists := services.GetMailboxChannel(config.InputMailboxUUID)
	if !exists {
		return worker.RuntimeErrorExit, fmt.Errorf("input mailbox channel not found: %s", config.InputMailboxUUID)
	}

	// Process messages from the input mailbox.
	for {
		select {
		case <-ctx.Done():
			return worker.NormalExit, nil
		case msg, ok := <-msgCh:
			if !ok {
				return worker.PrematureExit, fmt.Errorf("input mailbox channel closed")
			}

			// Ensure the message is of the expected type.
			m, ok := msg.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("unexpected message type: %T", msg)
			}

			// Lookup the destinations for the message's tag.
			destinations, exists := config.TagDestinations[m.Tag]
			if !exists {
				return worker.RuntimeErrorExit, fmt.Errorf("no destinations found for tag: %s", m.Tag)
			}

			// Broadcast the message with the same tag to all destination mailboxes.
			for _, dest := range destinations {
				if err := services.SendMessage(dest, m, config.BlockingSend); err != nil {
					return worker.RuntimeErrorExit, fmt.Errorf("failed to send message to %s: %w", dest, err)
				}
			}
		}
	}
}

// parseConfig unmarshals and validates the worker configuration.
func (w *BroadcastWorker) parseConfig(rawConfig any) (BroadcastWorkerConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return BroadcastWorkerConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config BroadcastWorkerConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return BroadcastWorkerConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return BroadcastWorkerConfig{}, fmt.Errorf("input_mailbox_uuid is required")
	}
	if len(config.TagDestinations) == 0 {
		return BroadcastWorkerConfig{}, fmt.Errorf("at least one tag mapping must be provided in configuration")
	}

	return config, nil
}

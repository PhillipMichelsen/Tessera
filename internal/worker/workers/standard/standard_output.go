package workers

import (
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// StandardOutputConfig represents the YAML configuration for the StandardOutputWorker.
type StandardOutputConfig struct {
	InputMailboxUUID   uuid.UUID `yaml:"input_mailbox_uuid"`
	InputMailboxBuffer int       `yaml:"input_mailbox_buffer"`
}

// StandardOutputWorker implements the worker.Worker interface.
// It simply prints any received message to standard output.
type StandardOutputWorker struct{}

// Run initializes the mailbox using the mailbox_uuid from configuration,
// then continuously prints any received message to stdout.
func (w *StandardOutputWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	config, err := w.parseRawConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse raw config: %w", err)
	}

	// Create mailbox and forward messages to inputChannel.
	services.CreateMailbox(config.InputMailboxUUID, config.InputMailboxBuffer)
	inputChannel, ok := services.GetMailboxChannel(config.InputMailboxUUID)
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to get input mailbox channel")
	}

	for {
		select {
		case <-ctx.Done():
			services.RemoveMailbox(config.InputMailboxUUID)
			return worker.NormalExit, nil
		case msg := <-inputChannel:
			// _ = msg
			fmt.Printf("%+v\n", msg)
		}
	}
}

func (w *StandardOutputWorker) parseRawConfig(rawConfig any) (StandardOutputConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return StandardOutputConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config StandardOutputConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return StandardOutputConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.InputMailboxUUID == uuid.Nil {
		return StandardOutputConfig{}, fmt.Errorf("mailbox_uuid is nil")
	}

	return config, nil
}

package workers

import (
	"context"
	"fmt"
	"time"

	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// StandardOutputConfig represents the YAML configuration for the StandardOutputWorker.
type StandardOutputConfig struct {
	MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
}

// StandardOutputWorker implements the worker.Worker interface.
// It simply prints any received message to standard output.
type StandardOutputWorker struct{}

// Run initializes the mailbox using the mailbox_uuid from configuration,
// then continuously prints any received message to stdout.
func (w *StandardOutputWorker) Run(ctx context.Context, config any, services worker.Services) (worker.ExitCode, error) {
	configBytes, ok := config.([]byte)
	if !ok {
		return worker.RuntimeErrorExit, fmt.Errorf("config is not in the expected []byte format")
	}

	var cfg StandardOutputConfig
	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if cfg.MailboxUUID == uuid.Nil {
		return worker.RuntimeErrorExit, fmt.Errorf("mailbox_uuid is required in configuration")
	}

	// Create mailbox and forward messages to inputChannel.
	inputChannel := make(chan worker.Message)
	services.CreateMailbox(cfg.MailboxUUID, worker.BuildChannelMessageReceiverForwarder(inputChannel))

	for {
		select {
		case <-ctx.Done():
			services.RemoveMailbox(cfg.MailboxUUID)
			return worker.NormalExit, nil
		case msg := <-inputChannel:
			// Print the entire message struct to stdout.
			fmt.Printf("[%s] Received message: %+v\n", time.Now().Format("15:04:05"), msg)
		}
	}
}

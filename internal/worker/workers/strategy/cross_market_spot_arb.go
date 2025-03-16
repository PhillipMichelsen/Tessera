package workers

import (
	"Tessera/internal/models"
	"Tessera/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// CrossMarketSpotArbitrageStrategyConfig represents the YAML configuration for the worker.
type CrossMarketSpotArbitrageStrategyConfig struct {
	Market1BookTickerMailboxUUID uuid.UUID `yaml:"market_1_book_ticker_mailbox_uuid"`
	Market2BookTickerMailboxUUID uuid.UUID `yaml:"market_2_book_ticker_mailbox_uuid"`
	MailboxBuffers               int       `yaml:"mailbox_buffers"`
	Output                       struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"output"`
	BlockingSend bool `yaml:"blocking_send"`
}

// CrossMarketSpotArbitrageStrategyWorker implements the worker.Worker interface.
type CrossMarketSpotArbitrageStrategyWorker struct {
	market1LastBookTicker models.BookTicker
	market2LastBookTicker models.BookTicker
}

func (w *CrossMarketSpotArbitrageStrategyWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	config, err := w.parseRawConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse raw config: %w", err)
	}

	// Create mailboxes for each market.
	market1Channel, err := services.CreateMailbox(config.Market1BookTickerMailboxUUID, config.MailboxBuffers)
	defer services.RemoveMailbox(config.Market1BookTickerMailboxUUID)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to create market 1 mailbox: %w", err)
	}

	market2Channel, err := services.CreateMailbox(config.Market2BookTickerMailboxUUID, config.MailboxBuffers)
	defer services.RemoveMailbox(config.Market2BookTickerMailboxUUID)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to create market 2 mailbox: %w", err)
	}

	// Main loop: listen for updates from both channels.
	for {
		select {
		case msg := <-market1Channel:
			message, ok := msg.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("invalid message type on market1 channel: %T", msg)
			}

			bookTicker, ok := message.Payload.(models.BookTicker)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("invalid payload type on market1 channel: %T", message.Payload)
			}

			w.market1LastBookTicker = bookTicker
			w.compareAndSend(services, config)
		case msg := <-market2Channel:
			message, ok := msg.(worker.Message)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("invalid message type on market2 channel: %T", msg)
			}

			bookTicker, ok := message.Payload.(models.BookTicker)
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("invalid payload type on market2 channel: %T", message.Payload)
			}

			w.market2LastBookTicker = bookTicker
			w.compareAndSend(services, config)
		case <-ctx.Done():
			return worker.NormalExit, nil
		}
	}
}

// compareAndSend computes the arbitrage differences and sends a message if an opportunity exists.
func (w *CrossMarketSpotArbitrageStrategyWorker) compareAndSend(services worker.Services, config CrossMarketSpotArbitrageStrategyConfig) {
	// Check if both tickers have been updated.
	if w.market1LastBookTicker.AskPrice == 0 || w.market2LastBookTicker.AskPrice == 0 {
		return
	}

	// Opportunity 1: Buy at Market1 (ask) and sell at Market2 (bid).
	percentDiff1 := ((w.market2LastBookTicker.BidPrice - w.market1LastBookTicker.AskPrice) / w.market1LastBookTicker.AskPrice) * 100

	// Opportunity 2: Buy at Market2 (ask) and sell at Market1 (bid).
	percentDiff2 := ((w.market1LastBookTicker.BidPrice - w.market2LastBookTicker.AskPrice) / w.market2LastBookTicker.AskPrice) * 100

	var message string
	// Determine the best opportunity based on the percentage difference.
	if percentDiff1 > percentDiff2 && percentDiff1 > 0 {
		message = fmt.Sprintf("Buy on Market1 at %.2f, sell on Market2 at %.2f: Profit = %.3f%%",
			w.market1LastBookTicker.AskPrice, w.market2LastBookTicker.BidPrice, percentDiff1)
	} else if percentDiff2 > 0 {
		message = fmt.Sprintf("Buy on Market2 at %.2f, sell on Market1 at %.2f: Profit = %.3f%%",
			w.market2LastBookTicker.AskPrice, w.market1LastBookTicker.BidPrice, percentDiff2)
	} else {
		return
	}

	// Send the message to the output mailbox.
	err := services.SendMessage(config.Output.MailboxUUID, worker.Message{
		Tag:     config.Output.Tag,
		Payload: message,
	}, config.BlockingSend)
	if err != nil {
		fmt.Printf("failed to send arbitrage message: %v\n", err)
	}
}

func (w *CrossMarketSpotArbitrageStrategyWorker) parseRawConfig(rawConfig any) (CrossMarketSpotArbitrageStrategyConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config CrossMarketSpotArbitrageStrategyConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.Market1BookTickerMailboxUUID == uuid.Nil {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("market_1_book_ticker_mailbox_uuid is required in configuration")
	}

	if config.Market2BookTickerMailboxUUID == uuid.Nil {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("market_2_book_ticker_mailbox_uuid is required in configuration")
	}

	if config.Output.MailboxUUID == uuid.Nil {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("output_mailbox_uuid is required in configuration")
	}

	if config.Output.Tag == "" {
		return CrossMarketSpotArbitrageStrategyConfig{}, fmt.Errorf("output_tag is required in configuration")
	}

	return config, nil
}

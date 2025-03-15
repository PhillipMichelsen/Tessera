package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

type BinanceSpotWebsocketStreamMessage struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// BinanceSpotWebsocketWorkerConfig defines the YAML configuration.
type BinanceSpotWebsocketWorkerConfig struct {
	BaseURL              string `yaml:"base_url"`
	StreamsOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"streams_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// BinanceSpotWebsocketWorker implements the worker.Worker interface.
type BinanceSpotWebsocketWorker struct{}

// Run reads the YAML config, connects to the Binance websocket, subscribes to the streams,
// and routes each received message to the configured destination mailboxes.
func (w *BinanceSpotWebsocketWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
	cfg, err := w.parseRawConfig(rawConfig)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to parse raw config: %w", err)
	}

	// Build a slice of stream names from the config.
	streamNames := make([]string, 0, len(cfg.StreamsOutputMapping))
	for streamName := range cfg.StreamsOutputMapping {
		streamNames = append(streamNames, streamName)
	}

	// Build the websocket URL. (Assumes config.BaseURL is provided without the "wss://" prefix.)
	u := url.URL{Scheme: "wss", Host: cfg.BaseURL, Path: "/stream"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to connect to Binance Spot WebSocket: %w", err)
	}
	// Ensure connection is closed on exit.
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("failed to close websocket connection: %v\n", err)
		}
	}(conn)

	// Subscribe to streams.
	subscribeMsg := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": streamNames,
		"id":     1,
	}
	if err = conn.WriteJSON(subscribeMsg); err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to send subscription message: %w", err)
	}

	// Create channels for messages and errors.
	msgCh := make(chan []byte)
	errCh := make(chan error)

	// Launch a goroutine that reads from the websocket.
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			msgCh <- message
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return worker.NormalExit, nil
		case err := <-errCh:
			return worker.RuntimeErrorExit, fmt.Errorf("failed to read message: %w", err)
		case message := <-msgCh:
			var msg BinanceSpotWebsocketStreamMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to unmarshal message: %w", err)
			}

			// In case of a subscription response, continue.
			if len(msg.Data) == 0 {
				fmt.Printf("Received subscription response: %s\n", message)
				continue
			}

			serializedJSON := models.SerializedJSON{
				JSON: string(msg.Data),
			}

			output, ok := cfg.StreamsOutputMapping[msg.Stream]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for stream: %s", msg.Stream)
			}

			if err := services.SendMessage(output.MailboxUUID, worker.Message{
				Tag:     output.Tag,
				Payload: serializedJSON,
			}, cfg.BlockingSend); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

func (w *BinanceSpotWebsocketWorker) parseRawConfig(rawConfig any) (BinanceSpotWebsocketWorkerConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return BinanceSpotWebsocketWorkerConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config BinanceSpotWebsocketWorkerConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return BinanceSpotWebsocketWorkerConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.BaseURL == "" {
		return BinanceSpotWebsocketWorkerConfig{}, fmt.Errorf("base_url is required in configuration")
	}
	if len(config.StreamsOutputMapping) == 0 {
		return BinanceSpotWebsocketWorkerConfig{}, fmt.Errorf("at least one stream must be provided in configuration")
	}

	return config, nil
}

package workers

import (
	protos "AlgorithmicTraderDistributed/internal/protos/mexc"
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"net/url"
)

// MEXCSpotWebsocketWorkerConfig defines the YAML configuration.
type MEXCSpotWebsocketWorkerConfig struct {
	BaseURL              string `yaml:"base_url"`
	StreamsOutputMapping map[string]struct {
		MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
		Tag         string    `yaml:"tag"`
	} `yaml:"streams_output_mapping"`
	BlockingSend bool `yaml:"blocking_send"`
}

// MEXCSpotWebsocketWorker implements the worker.Worker interface.
type MEXCSpotWebsocketWorker struct{}

func (w *MEXCSpotWebsocketWorker) Run(ctx context.Context, rawConfig any, services worker.Services) (worker.ExitCode, error) {
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
	u := url.URL{Scheme: "wss", Host: cfg.BaseURL, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return worker.RuntimeErrorExit, fmt.Errorf("failed to connect to MEXC Spot websocket: %w", err)
	}
	// Ensure connection is closed on exit.
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("failed to close websocket connection: %v\n", err)
		}
	}(conn)

	subscribeMsg := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": streamNames,
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
			// Catch the subscription response, which for some reason is sent as a JSON object.
			if message[0] == '{' {
				fmt.Printf("Received subscription response: %+v\n", string(message))
				continue
			}

			// Otherwise, assume it's a protobuf PushDataV3ApiWrapper message.
			var msg protos.PushDataV3ApiWrapper
			if err := proto.Unmarshal(message, &msg); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to unmarshal protobuf message: %w", err)
			}

			output, ok := cfg.StreamsOutputMapping[msg.Channel]
			if !ok {
				return worker.RuntimeErrorExit, fmt.Errorf("destination mapping not found for channel: %s", msg.Channel)
			}

			if err := services.SendMessage(output.MailboxUUID, worker.Message{
				Tag:     output.Tag,
				Payload: &msg,
			}, cfg.BlockingSend); err != nil {
				return worker.RuntimeErrorExit, fmt.Errorf("failed to send message: %w", err)
			}
		}
	}

}

// parseRawConfig converts the raw YAML configuration into MEXCSpotWebsocketWorkerConfig.
func (w *MEXCSpotWebsocketWorker) parseRawConfig(rawConfig any) (MEXCSpotWebsocketWorkerConfig, error) {
	configBytes, ok := rawConfig.([]byte)
	if !ok {
		return MEXCSpotWebsocketWorkerConfig{}, fmt.Errorf("config is not in the expected []byte format")
	}

	var config MEXCSpotWebsocketWorkerConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return MEXCSpotWebsocketWorkerConfig{}, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if config.BaseURL == "" {
		return MEXCSpotWebsocketWorkerConfig{}, fmt.Errorf("base_url is required in configuration")
	}
	if len(config.StreamsOutputMapping) == 0 {
		return MEXCSpotWebsocketWorkerConfig{}, fmt.Errorf("at least one stream must be provided in configuration")
	}

	return config, nil
}

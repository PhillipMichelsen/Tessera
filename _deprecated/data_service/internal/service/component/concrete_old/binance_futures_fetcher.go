package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type BinanceFuturesWebsocketFetcher struct {
	name          string
	uuid          uuid.UUID
	running       bool
	config        map[string]string
	routeFunction func(data models.MarketDataPiece)

	conn *websocket.Conn
	mu   sync.Mutex
}

func NewBinanceFuturesWebsocketFetcher(routeFunction func(data models.MarketDataPiece), config map[string]string) *BinanceFuturesWebsocketFetcher {
	return &BinanceFuturesWebsocketFetcher{
		name:          "BinanceFuturesWebsocketFetcher",
		uuid:          uuid.New(),
		running:       false,
		config:        config,
		routeFunction: routeFunction,
	}
}

func (r *BinanceFuturesWebsocketFetcher) Handle(data models.MarketDataPiece) {
	return
}

func (r *BinanceFuturesWebsocketFetcher) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("receiver is already running")
	}

	streams, ok := r.config["streams"]
	var streamsArray []string
	if !ok {
		return errors.New("streams not found in config")
	}
	if streams != "" {
		streamsArray = strings.Split(streams, ",")
	}

	baseURL, ok := r.config["base_url"]
	if !ok {
		return errors.New("base_url not found in config")
	}

	// Create WebSocket URL without stream-specific query parameters
	u := url.URL{Scheme: "wss", Host: baseURL, Path: "/stream"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Binance Futures WebSocket: %w", err)
	}
	r.conn = conn

	// Create subscription message to subscribe to streams
	subscribeMsg := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": streamsArray,
		"id":     1,
	}
	if err := r.conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to send subscription message: %w", err)
	}
	log.Printf("Subscribed to streams on Binance Futures Websocket: %v", streamsArray)

	// Start listening for incoming messages
	r.running = true
	go r.listen()

	return nil
}

func (r *BinanceFuturesWebsocketFetcher) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return errors.New("receiver is not running")
	}

	r.running = false

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close WebSocket connection: %w", err)
	}

	log.Println("BinanceFuturesWebsocketFetcher stopped")
	return nil
}

func (r *BinanceFuturesWebsocketFetcher) GetName() string {
	return r.name
}

func (r *BinanceFuturesWebsocketFetcher) GetUUID() uuid.UUID {
	return r.uuid
}

// listen reads incoming WebSocket messages and handles them
func (r *BinanceFuturesWebsocketFetcher) listen() {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(r.conn)

	for r.running {
		_, message, err := r.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Handle each message based on stream type
		r.HandleMessage(message)
	}
}

// HandleMessage processes each WebSocket message dynamically
func (r *BinanceFuturesWebsocketFetcher) HandleMessage(message []byte) {
	var msg struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"` // Use RawMessage to defer JSON decoding
	}

	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	streamInfo := r.parseStreamInformation(msg.Stream)
	if streamInfo == nil {
		log.Printf("Unrecognized stream type or subscription confirmation: %v, %v", msg.Stream, string(msg.Data))
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		log.Printf("Error unmarshalling data field: %v", err)
		return
	}

	// Populate `MarketDataPiece` structure
	marketDataPiece := models.MarketDataPiece{
		Source:           "BinanceFutures",
		Symbol:           streamInfo["symbol"],
		BaseType:         streamInfo["baseType"],
		Interval:         streamInfo["interval"],
		ReceiveTimestamp: time.Now().UTC(),
		Payload:          models.MarketData(data),
	}

	// Pass the structured data to the routing function
	r.routeFunction(marketDataPiece)
}

// parseStreamInformation extracts metadata from stream names for Futures
func (r *BinanceFuturesWebsocketFetcher) parseStreamInformation(stream string) map[string]string {
	stream = strings.TrimSpace(stream)
	parts := strings.Split(stream, "@")
	if len(parts) < 2 {
		return nil
	}

	result := make(map[string]string)
	result["symbol"] = parts[0]
	streamType := parts[1]

	switch {
	case strings.HasPrefix(streamType, "kline_"):
		result["baseType"] = "kline"
		result["interval"] = strings.TrimPrefix(streamType, "kline_")
	case streamType == "aggTrade":
		result["baseType"] = "trade"
		result["interval"] = "100ms"
	case streamType == "depth":
		result["baseType"] = "depth"
		result["interval"] = "250ms"
	case streamType == "markPrice":
		result["baseType"] = "markPrice"
		result["interval"] = "1s"
	case streamType == "bookTicker":
		result["baseType"] = "bookTicker"
		result["interval"] = "Realtime"
	case streamType == "forceOrder":
		result["baseType"] = "forceOrder"
		result["interval"] = "1s"
	default:
		result["baseType"] = "UnrecognizedStreamType"
	}

	return result
}

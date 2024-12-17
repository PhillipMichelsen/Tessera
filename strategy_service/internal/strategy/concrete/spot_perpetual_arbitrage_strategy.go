package strategy

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/strategy_service/internal/helper"
	"AlgorithimcTraderDistributed/strategy_service/internal/strategy"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

type SpotPerpetualArbitrageStrategy struct {
	// Core
	StrategyUUID  string
	sessionUUID   string
	positions     []strategy.Position
	orders        []strategy.Order
	isInitialized bool
	isTrading     bool
	isConsuming   bool

	// Configuration
	spotAsset      strategy.Asset
	perpetualAsset strategy.Asset

	// State
	currentPerpetualBookTicker models.BookTicker
	currentSpotBookTicker      models.BookTicker
	largestContago             float64
	smallestContago            float64
	largestBackwardation       float64
	smallestBackwardation      float64

	// RabbitMQ
	rabbitMQChannel *amqp.Channel
	rabbitMQQueue   amqp.Queue
}

func NewSpotPerpetualArbitrageStrategy(sessionUuid string, rabbitMQChannel *amqp.Channel) *SpotPerpetualArbitrageStrategy {
	return &SpotPerpetualArbitrageStrategy{
		StrategyUUID:    uuid.New().String(),
		sessionUUID:     sessionUuid,
		positions:       []strategy.Position{},
		orders:          []strategy.Order{},
		rabbitMQChannel: rabbitMQChannel,
		isInitialized:   false,
		isTrading:       false,
		isConsuming:     false,
	}
}

// Initialize configures the strategy
func (s *SpotPerpetualArbitrageStrategy) Initialize(config map[string]string) error {
	// Parse config and initialize components
	s.spotAsset = strategy.AssetFromString(config["spotAsset"])
	s.perpetualAsset = strategy.AssetFromString(config["perpetualAsset"])

	// Setup RabbitMQ
	var err error
	s.rabbitMQQueue, err = helper.CreateRabbitMQQueueForStrategy(
		s.rabbitMQChannel,
		"SpotPerpetualArbitrageStrategy",
		s.StrategyUUID,
		s.sessionUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ queue: %v", err)
	}

	err = helper.BindQueueToMarketDataExchange(
		s.rabbitMQChannel,
		s.rabbitMQQueue,
		models.RoutingPatternToAMQPTable(
			models.RoutingPattern{
				Source:     s.spotAsset.Source,
				Symbol:     s.spotAsset.Symbol,
				BaseType:   "bookTicker",
				Interval:   "*",
				Processors: "BinanceSpotBookTickerTransformerProcessor.",
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	err = helper.BindQueueToMarketDataExchange(
		s.rabbitMQChannel,
		s.rabbitMQQueue,
		models.RoutingPatternToAMQPTable(
			models.RoutingPattern{
				Source:     s.perpetualAsset.Source,
				Symbol:     s.perpetualAsset.Symbol,
				BaseType:   "bookTicker",
				Interval:   "*",
				Processors: "BinanceFuturesBookTickerTransformerProcessor.",
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	s.isInitialized = true
	return nil
}

// Start begins the data consumption process and optionally auto-lives the strategy
func (s *SpotPerpetualArbitrageStrategy) Start(autoLive bool) error {
	if !s.isInitialized {
		return fmt.Errorf("strategy not initialized")
	}

	s.isConsuming = true
	go s.consumeMessages()

	if autoLive {
		go func() {
			for s.isConsuming {
				if s.IsReady() {
					if err := s.BeginTrading(); err != nil {
						fmt.Printf("failed to bring strategy live: %v\n", err)
					}
					return
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	log.Printf("Strategy %s has been started.", s.StrategyUUID)
	go s.printLargestDifference()
	return nil
}

func (s *SpotPerpetualArbitrageStrategy) printLargestDifference() {
	for {
		time.Sleep(60 * time.Second)
		fmt.Printf("Largest Contago: %f | Smallest Contago: %f | Largest Backwardation: %f | Smallest Backwardation: %f\n", s.largestContago, s.smallestContago, s.largestBackwardation, s.smallestBackwardation)

		s.largestContago = math.Inf(-1)
		s.smallestContago = math.Inf(1)
		s.largestBackwardation = math.Inf(-1)
		s.smallestBackwardation = math.Inf(1)
	}
}

func (s *SpotPerpetualArbitrageStrategy) IsReady() bool {
	return true
}

func (s *SpotPerpetualArbitrageStrategy) BeginTrading() error {
	if !s.isInitialized {
		return fmt.Errorf("strategy not initialized")
	}

	if !s.IsReady() {
		return fmt.Errorf("strategy not ready")
	}

	s.isTrading = true
	log.Printf("Strategy %s is now live.", s.StrategyUUID)
	return nil
}

// Stop gracefully stops the strategy
func (s *SpotPerpetualArbitrageStrategy) Stop() error {
	s.isConsuming = false
	s.isTrading = false

	if err := s.rabbitMQChannel.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
	}
	return nil
}

// Private Methods

// consumeMessages listens for messages and processes them
func (s *SpotPerpetualArbitrageStrategy) consumeMessages() {
	msgs, err := s.rabbitMQChannel.Consume(
		s.rabbitMQQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("failed to register consumer: %v\n", err)
		return
	}

	for s.isConsuming {
		select {
		case msg := <-msgs:
			s.process(msg.Body)
		}
	}
}

// process handles message data and optionally performs trading logic
func (s *SpotPerpetualArbitrageStrategy) process(messageBody []byte) {
	marketDataPiece, err := models.ConvertProtoBytesToGoMarketDataPiece(messageBody)
	if err != nil {
		fmt.Printf("failed to convert message to MarketDataPiece: %v\n", err)
		return
	}
	bookTicker, ok := marketDataPiece.Payload.(models.BookTicker)
	if !ok {
		fmt.Println("message payload is not OHLCV")
		return
	}

	if strategy.AssetFromMarketDataPiece(marketDataPiece) == s.spotAsset {
		s.currentSpotBookTicker = bookTicker
	}
	if strategy.AssetFromMarketDataPiece(marketDataPiece) == s.perpetualAsset {
		s.currentPerpetualBookTicker = bookTicker
	}

	if s.isTrading && s.IsReady() {
		s.processDecision()
	}
}

// performCriterionChecks evaluates moving averages and sends orders
func (s *SpotPerpetualArbitrageStrategy) processDecision() {
	if s.currentSpotBookTicker.BidPrice == 0 || s.currentPerpetualBookTicker.BidPrice == 0 {
		return
	}

	shortSpotLongPerpProfit := ((s.currentPerpetualBookTicker.AskPrice - s.currentSpotBookTicker.BidPrice) / s.currentSpotBookTicker.BidPrice) * 100
	longSpotShortPerpProfit := ((s.currentSpotBookTicker.AskPrice - s.currentPerpetualBookTicker.BidPrice) / s.currentPerpetualBookTicker.BidPrice) * 100

	if shortSpotLongPerpProfit > 0 { // This indicates contango
		if shortSpotLongPerpProfit > s.largestContago {
			s.largestContago = shortSpotLongPerpProfit
		}
		if shortSpotLongPerpProfit < s.smallestContago {
			s.smallestContago = shortSpotLongPerpProfit
		}
	}

	// Check and update largest/smallest backwardation
	if longSpotShortPerpProfit > 0 { // This indicates backwardation
		if longSpotShortPerpProfit > s.largestBackwardation {
			s.largestBackwardation = longSpotShortPerpProfit
		}
		if longSpotShortPerpProfit < s.smallestBackwardation {
			s.smallestBackwardation = longSpotShortPerpProfit
		}
	}
}

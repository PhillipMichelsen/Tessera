package strategy

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/strategy_service/internal/component"
	"AlgorithimcTraderDistributed/strategy_service/internal/helper"
	"AlgorithimcTraderDistributed/strategy_service/internal/strategy"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

type OrderBookImbalanceStrategy struct {
	// Core
	strategyUUID string
	sessionUUID  string
	positions    []strategy.Position
	orders       []strategy.Order
	status       strategy.Status
	autoLive     bool

	// Configuration
	asset strategy.Asset

	// Components
	orderBookImbalance *component.OrderBookSnapshotImbalance
	zScore             *component.RollingFloat64ZScoreComponent

	// RabbitMQ
	rabbitMQChannel *amqp.Channel
	rabbitMQQueue   amqp.Queue
}

func NewOrderBookImbalanceStrategy(strategyUUID string, sessionUUID string, rabbitMQChannel *amqp.Channel) *OrderBookImbalanceStrategy {
	return &OrderBookImbalanceStrategy{
		strategyUUID:    strategyUUID,
		sessionUUID:     sessionUUID,
		positions:       []strategy.Position{},
		orders:          []strategy.Order{},
		rabbitMQChannel: rabbitMQChannel,
		status:          strategy.Uninitialized,
	}
}

// Initialize configures the strategy
func (s *OrderBookImbalanceStrategy) Initialize(config map[string]string) error {
	// Parse config and initialize components
	s.asset = strategy.AssetFromString(config["asset"])

	zScorePeriod, err := strconv.Atoi(config["zScorePeriod"])
	if err != nil {
		return fmt.Errorf("invalid zScorePeriod: %v", err)
	}
	s.zScore = component.NewRollingFloat64ZScoreComponent(zScorePeriod)
	s.orderBookImbalance = component.NewOrderBookSnapshotImbalance()

	// Setup RabbitMQ
	s.rabbitMQQueue, err = helper.CreateRabbitMQQueueForStrategy(
		s.rabbitMQChannel,
		"OrderBookImbalanceStrategy",
		s.strategyUUID,
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
				Source:     s.asset.Source,
				Symbol:     s.asset.Symbol,
				BaseType:   "depth",
				Interval:   "*",
				Processors: "BinanceSpotOrderBookUpdateTransformerProcessor.BinanceSpotOrderBookAggregatorProcessor.OrderBookSnapshotRangeFilterProcessor_2.00%.",
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	s.status = strategy.Initialized
	return nil
}

// Start begins the data consumption process and optionally begins trading automatically when ready
func (s *OrderBookImbalanceStrategy) Start() error {
	if s.status < strategy.Initialized {
		return fmt.Errorf("strategy not initialized")
	}

	s.status = strategy.Staging
	go s.consumeMessages()

	log.Printf("Strategy %s has been started.", s.strategyUUID)
	return nil
}

func (s *OrderBookImbalanceStrategy) SetLive(live bool) error {
	if live {
		if s.status != strategy.Staged {
			return fmt.Errorf("strategy not staged")
		}
		s.status = strategy.Live
	} else {
		s.status = strategy.Staged
	}

	return nil
}

func (s *OrderBookImbalanceStrategy) SetAutoLive(autoLive bool) {
	s.autoLive = autoLive
}

// Stop gracefully stops the strategy
func (s *OrderBookImbalanceStrategy) Stop() error {
	s.status = strategy.Uninitialized

	if err := s.rabbitMQChannel.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
	}
	return nil
}

// Private Methods

// consumeMessages listens for messages and processes them
func (s *OrderBookImbalanceStrategy) consumeMessages() {
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

	for s.status >= strategy.Staging {
		select {
		case msg := <-msgs:
			s.process(msg.Body)
		}
	}
}

// process handles message data and optionally performs trading logic
func (s *OrderBookImbalanceStrategy) process(messageBody []byte) {
	timeNow := time.Now()
	marketDataPiece, err := models.ConvertProtoBytesToGoMarketDataPiece(messageBody)
	if err != nil {
		fmt.Printf("failed to convert message to MarketDataPiece: %v\n", err)
		return
	}
	orderbookSnapshot, ok := marketDataPiece.Payload.(models.OrderBookSnapshot)
	if !ok {
		fmt.Println("message payload is not OrderBookSnapshot")
		return
	}

	imbalance := s.orderBookImbalance.CalculateOrderBookSnapshotImbalance(orderbookSnapshot)
	s.zScore.AddNewData(imbalance)

	if s.status == strategy.Staging {
		s.checkIfStaged()
	}
	if s.status == strategy.Staged {
		if s.autoLive {
			s.status = strategy.Live
			fmt.Println("Strategy is live")
		}
	}
	if s.status == strategy.Live {
		log.Printf("Imbalance: %f, Z-Score: %f", imbalance, s.zScore.GetZScore())
		log.Printf("Data Received by data service at: %v, recieved and fully processed in: %v. Process function isolated: %v", marketDataPiece.ReceiveTimestamp, time.Since(marketDataPiece.SendTimestamp), time.Since(timeNow))
	}
}

func (s *OrderBookImbalanceStrategy) checkIfStaged() {
	if s.zScore.IsReady() {
		s.status = strategy.Staged
		fmt.Println("Strategy is staged")
		return
	}
}

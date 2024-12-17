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
)

type MovingAverageCrossoverStrategy struct {
	// Core
	strategyUUID string
	sessionUUID  string
	positions    []strategy.Position
	orders       []strategy.Order
	status       strategy.Status
	autoLive     bool

	// Configuration
	asset      strategy.Asset
	fastPeriod int
	slowPeriod int

	// Components
	fastRollingMA *component.MovingAverageOHLCVRollingComponent
	slowRollingMA *component.MovingAverageOHLCVRollingComponent

	// RabbitMQ
	rabbitMQChannel *amqp.Channel
	rabbitMQQueue   amqp.Queue
}

func NewMovingAverageCrossoverStrategy(strategyUUID string, sessionUUID string, rabbitMQChannel *amqp.Channel) *MovingAverageCrossoverStrategy {
	return &MovingAverageCrossoverStrategy{
		strategyUUID:    strategyUUID,
		sessionUUID:     sessionUUID,
		positions:       []strategy.Position{},
		orders:          []strategy.Order{},
		rabbitMQChannel: rabbitMQChannel,
		status:          strategy.Uninitialized,
	}
}

// Initialize configures the strategy
func (s *MovingAverageCrossoverStrategy) Initialize(config map[string]string) error {
	// Parse config and initialize components
	fastPeriod, err := strconv.Atoi(config["fastMAPeriod"])
	if err != nil {
		return fmt.Errorf("invalid fastMAPeriod: %v", err)
	}
	s.fastPeriod = fastPeriod

	slowPeriod, err := strconv.Atoi(config["slowMAPeriod"])
	if err != nil {
		return fmt.Errorf("invalid slowMAPeriod: %v", err)
	}
	s.slowPeriod = slowPeriod

	s.asset = strategy.AssetFromString(config["asset"])

	s.fastRollingMA = component.NewMovingAverageOHLCVRollingComponent(s.fastPeriod)
	s.slowRollingMA = component.NewMovingAverageOHLCVRollingComponent(s.slowPeriod)

	// Setup RabbitMQ
	s.rabbitMQQueue, err = helper.CreateRabbitMQQueueForStrategy(
		s.rabbitMQChannel,
		"MovingAverageCrossoverStrategy",
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
				BaseType:   "kline",
				Interval:   "1s",
				Processors: "*",
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	s.status = strategy.Initialized
	return nil
}

// Start begins the data consumption process and optionally begins trading automatically when ready
func (s *MovingAverageCrossoverStrategy) Start() error {
	if s.status < strategy.Initialized {
		return fmt.Errorf("strategy not initialized")
	}

	s.status = strategy.Staging
	go s.consumeMessages()

	log.Printf("Strategy %s has been started.", s.strategyUUID)
	return nil
}

func (s *MovingAverageCrossoverStrategy) SetLive(live bool) error {
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

func (s *MovingAverageCrossoverStrategy) SetAutoLive(autoLive bool) {
	s.autoLive = autoLive
}

// Stop gracefully stops the strategy
func (s *MovingAverageCrossoverStrategy) Stop() error {
	s.status = strategy.Uninitialized

	if err := s.rabbitMQChannel.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
	}
	return nil
}

// Private Methods

// consumeMessages listens for messages and processes them
func (s *MovingAverageCrossoverStrategy) consumeMessages() {
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
func (s *MovingAverageCrossoverStrategy) process(messageBody []byte) {
	marketDataPiece, err := models.ConvertProtoBytesToGoMarketDataPiece(messageBody)
	if err != nil {
		fmt.Printf("failed to convert message to MarketDataPiece: %v\n", err)
		return
	}
	ohlcv, ok := marketDataPiece.Payload.(models.OHLCV)
	if !ok {
		fmt.Println("message payload is not OHLCV")
		return
	}
	s.processComponentUpdate(ohlcv)

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
		s.processDecision()
	}
}

func (s *MovingAverageCrossoverStrategy) checkIfStaged() {
	if s.fastRollingMA.IsReady() && s.slowRollingMA.IsReady() {
		s.status = strategy.Staged
		fmt.Println("Strategy is staged")
		return
	}
}

func (s *MovingAverageCrossoverStrategy) processComponentUpdate(ohlcv models.OHLCV) {
	s.fastRollingMA.AddNewData(ohlcv)
	s.slowRollingMA.AddNewData(ohlcv)
}

// performCriterionChecks evaluates moving averages and sends orders
func (s *MovingAverageCrossoverStrategy) processDecision() {
	fastMA := s.fastRollingMA.GetMovingAverage()
	slowMA := s.slowRollingMA.GetMovingAverage()

	if fastMA > slowMA {
		fmt.Println("Buy signal triggered, Fast MA:", fastMA, "Slow MA:", slowMA)
	} else if fastMA < slowMA {
		fmt.Println("Sell signal triggered, Fast MA:", fastMA, "Slow MA:", slowMA)
	}
}

package manager

import (
	"AlgorithimcTraderDistributed/strategy_service/internal/strategy"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Manager struct {
	strategies         map[string]strategy.Strategy
	rabbitMQConnection *amqp.Connection
}

func NewManager(sessionUUID string, rabbitMQURL string) (*Manager, error) {
	rabbitMQConnection, err := amqp.DialConfig(rabbitMQURL, amqp.Config{
		Properties: amqp.Table{
			"connection_name": fmt.Sprintf("StrategyService_%s", sessionUUID),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return &Manager{
		strategies:         make(map[string]strategy.Strategy),
		rabbitMQConnection: rabbitMQConnection,
	}, nil
}

func (m *Manager) AddStrategy(strategyUUID string, strategy strategy.Strategy) {
	m.strategies[strategyUUID] = strategy
}

func (m *Manager) GetStrategy(strategyUUID string) strategy.Strategy {
	return m.strategies[strategyUUID]
}

func (m *Manager) GetRabbitMQChannel() (*amqp.Channel, error) {
	channel, err := m.rabbitMQConnection.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (m *Manager) StartAllStrategiesAndAutoLive() {
	for _, strat := range m.strategies {
		err := strat.Start()
		if err != nil {
			return
		}
		strat.SetAutoLive(true)
	}
}

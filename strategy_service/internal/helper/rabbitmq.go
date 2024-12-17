package helper

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateRabbitMQQueueForStrategy(channel *amqp.Channel, strategyName, strategyUUID, sessionUUID string) (amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		fmt.Sprintf("%s_%s_%s", strategyName, strategyUUID, sessionUUID),
		false,
		true,
		false,
		true,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue: %v", err)
	}

	return queue, nil
}

func BindQueueToMarketDataExchange(channel *amqp.Channel, queue amqp.Queue, headerParameters amqp.Table) error {
	err := channel.QueueBind(
		queue.Name,
		"",
		"market_data",
		false,
		headerParameters,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}

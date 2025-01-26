package sender

import (
	"AlgorithimcTraderDistributed/common/models"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

type RabbitMQSender struct {
	SessionUUID string
	RabbitMQURL string
	Exchange    string
	conn        *amqp.Connection
	channel     *amqp.Channel
}

// NewRabbitMQSender creates a new instance of RabbitMQSender (no connection yet).
func NewRabbitMQSender(sessionUUID, rabbitmqurl, exchange string) *RabbitMQSender {
	return &RabbitMQSender{
		SessionUUID: sessionUUID,
		RabbitMQURL: rabbitmqurl,
		Exchange:    exchange,
	}
}

// Start initializes the connection, channel, and exchange for RabbitMQ.
func (r *RabbitMQSender) Start() error {
	// Establish connection to RabbitMQ with a custom name (via the ConnectionName option).
	conn, err := amqp.DialConfig(r.RabbitMQURL, amqp.Config{
		Properties: amqp.Table{
			"connection_name": fmt.Sprintf("DataService_%s", r.SessionUUID),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare the exchange with type "headers" instead of "topic"
	err = channel.ExchangeDeclare(
		r.Exchange,
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Set the connection and channel
	r.conn = conn
	r.channel = channel

	log.Printf("RabbitMQ Sender started with exchange %s", r.Exchange)
	return nil
}

// Send sends a serialized MarketDataPiece to RabbitMQ using headers for routing.
func (r *RabbitMQSender) Send(marketDataPiece models.MarketDataPiece) {
	marketDataPiece.SendTimestamp = time.Now()

	marketDataPieceAsProto := models.ConvertGoToProtoMarketDataPiece(marketDataPiece)
	bytes, err := proto.Marshal(marketDataPieceAsProto)
	if err != nil {
		log.Printf("failed to marshal market data piece: %v", err)
		return
	}

	// Publish the message with headers for routing
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        bytes,
			Headers:     models.RoutingPatternToAMQPTable(models.MarketDataPieceToRoutingPattern(marketDataPiece)),
		},
	)
	if err != nil {
		log.Printf("failed to publish message: %v", err)
		return
	}
}

// Stop closes the RabbitMQ connection and channel.
func (r *RabbitMQSender) Stop() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("error closing channel: %v", err)
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("error closing connection: %v", err)
		}
	}

	log.Println("RabbitMQ connection and channel closed.")
	return nil
}

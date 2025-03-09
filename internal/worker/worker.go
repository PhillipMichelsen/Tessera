package worker

import (
	"context"
	"github.com/google/uuid"
)

// ExitCode represents the exit status of a worker.
type ExitCode int

const (
	NormalExit ExitCode = iota
	PrematureExit
	RuntimeErrorExit
	PanicExit
)

// Services defines the services (interface) that a worker can use to interact with the system.
type Services interface {
	SendMessage(message OutboundMessage) error
	StartReceivingMessages(receiverFunc func(message InboundMessage))
	StopReceivingMessages()
}

// InboundMessage represents a message received by a worker.
type InboundMessage struct {
	SourceWorkerUUID uuid.UUID
	MessageTag       string
	Payload          interface{}
}

// OutboundMessage represents a message to be sent by a worker.
type OutboundMessage struct {
	DestinationWorkerUUID uuid.UUID
	MessageTag            string
	Payload               interface{}
}

type Worker interface {
	Run(ctx context.Context, config map[string]interface{}, services Services) (ExitCode, error)
	GetWorkerName() string
}

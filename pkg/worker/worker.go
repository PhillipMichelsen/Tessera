package worker

import (
	"context"
	"github.com/google/uuid"
	"time"
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
	SendMessage(destinationWorkerUUID uuid.UUID, payload interface{}) error
	StartReceivingMessages(receiverFunc func(message InboundMessage))
	StopReceivingMessages()
}

// InboundMessage represents a message received by a worker.
type InboundMessage struct {
	SourceWorkerUUID uuid.UUID
	SentTimestamp    time.Time
	Payload          interface{}
}

// Worker is the interface that concrete workers implement.
type Worker interface {
	Run(ctx context.Context, config map[string]interface{}, services Services) (ExitCode, error)
	GetWorkerName() string
}

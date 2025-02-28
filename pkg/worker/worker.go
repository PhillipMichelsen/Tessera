package worker

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// ExitCode represents exit status.
type ExitCode int

const (
	NormalExit ExitCode = iota
	PrematureExit
	RuntimeErrorExit
	PanicExit
)

type InstanceServices interface {
	SendMessage(destinationWorkerUUID uuid.UUID, payload interface{}) error
	StartReceivingMessages(receiverFunc func(message InboundMessage))
	StopReceivingMessages()
}

type InboundMessage struct {
	SourceWorkerUUID uuid.UUID
	SentTimestamp    time.Time
	Payload          interface{}
}

// Worker is the interface that concrete workers implement.
type Worker interface {
	Run(ctx context.Context, config map[string]interface{}, instanceServices InstanceServices) (ExitCode, error)
	GetWorkerName() string
}

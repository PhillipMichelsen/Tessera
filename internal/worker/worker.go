package worker

import (
	"context"
	"github.com/google/uuid"
)

// ExitCode represents the exit status of a worker. It is used to communicate the reason for the worker's termination to the node.
type ExitCode int

const (
	NormalExit ExitCode = iota
	PrematureExit
	RuntimeErrorExit
	PanicExit
)

// Services defines the services (interface) that a worker can use to interact with the system.
type Services interface {
	SendMessage(destinationMailboxUUID uuid.UUID, message Message, block bool) error
	CreateMailbox(mailboxUUID uuid.UUID, bufferSize int) (<-chan any, error)
	RemoveMailbox(mailboxUUID uuid.UUID)
}

// Message represents a message that can be sent or received by a worker. Identifications of source and purpose are done via tags.
type Message struct {
	Tag     string
	Payload interface{}
}

// Worker interface. Every worker that wants to be deployed by the node must implement this interface.
type Worker interface {
	Run(ctx context.Context, rawConfig any, services Services) (ExitCode, error)
}

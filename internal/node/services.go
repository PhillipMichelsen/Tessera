package node

import (
	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/google/uuid"
)

// Ensure Node implements the worker.Services interface.
var _ worker.Services = &WorkerServices{}

// WorkerServices provides methods for a worker to interact with the system.
type WorkerServices struct {
	instance   *Node
	workerUUID uuid.UUID
}

// NewWorkerServices initializes a new WorkerServices.
func NewWorkerServices(instance *Node, workerUUID uuid.UUID) *WorkerServices {
	return &WorkerServices{
		instance:   instance,
		workerUUID: workerUUID,
	}
}

// TODO: To be overhauled when the router is implemented. Currently, it sends messages directly with the dispatcher.

// SendMessage sends a message to another worker.
func (ws *WorkerServices) SendMessage(message worker.OutboundMessage) error {
	err := ws.instance.dispatcher.SendMessage(message.DestinationWorkerUUID, worker.InboundMessage{
		SourceWorkerUUID: ws.workerUUID,
		MessageTag:       message.MessageTag,
		Payload:          message.Payload,
	})
	if err != nil {
		return err
	}
	return nil
}

// StartReceivingMessages registers a mailbox to start receiving messages.
// It adapts the internal MailboxMessage into the worker's InboundMessage type.
func (ws *WorkerServices) StartReceivingMessages(receiverFunc func(message worker.InboundMessage)) {
	ws.instance.dispatcher.CreateMailbox(
		ws.workerUUID,
		func(message interface{}) {
			receiverFunc(message.(worker.InboundMessage))
		},
		1000,
	)
}

// StopReceivingMessages removes the mailbox for the worker.
func (ws *WorkerServices) StopReceivingMessages() {
	ws.instance.dispatcher.RemoveMailbox(ws.workerUUID)
}

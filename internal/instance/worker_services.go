package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
	"github.com/google/uuid"
)

// Ensure that WorkerServices implements the worker.Services interface.
var _ worker.Services = &WorkerServices{}

// WorkerServices provides methods for a worker to interact with the system.
type WorkerServices struct {
	instance   *Instance
	workerUUID uuid.UUID
}

// NewWorkerServices initializes a new WorkerServices.
func NewWorkerServices(instance *Instance, workerUUID uuid.UUID) *WorkerServices {
	return &WorkerServices{
		instance:   instance,
		workerUUID: workerUUID,
	}
}

// SendMessage sends a message to another worker.
func (ws *WorkerServices) SendMessage(destinationWorkerUUID uuid.UUID, payload interface{}) error {
	return ws.instance.dispatcher.SendMessage(ws.workerUUID, destinationWorkerUUID, payload)
}

// StartReceivingMessages registers a mailbox to start receiving messages.
// It adapts the internal MailboxMessage into the worker's InboundMessage type.
func (ws *WorkerServices) StartReceivingMessages(receiverFunc func(message worker.InboundMessage)) {
	ws.instance.dispatcher.CreateMailbox(ws.workerUUID, func(message MailboxMessage) {
		receiverFunc(worker.InboundMessage{
			SourceWorkerUUID: message.SourceWorkerUUID,
			SentTimestamp:    message.SentTimestamp,
			Payload:          message.Payload,
		})
	})
}

// StopReceivingMessages removes the mailbox for the worker.
func (ws *WorkerServices) StopReceivingMessages() {
	ws.instance.dispatcher.RemoveMailbox(ws.workerUUID)
}

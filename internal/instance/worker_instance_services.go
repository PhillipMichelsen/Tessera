package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
	"github.com/google/uuid"
)

type WorkerInstanceServices struct {
	Instance   *Instance
	WorkerUUID uuid.UUID
}

func (wis *WorkerInstanceServices) SendMessage(recipientWorkerUUID uuid.UUID, payload interface{}) error {
	return wis.Instance.dispatcher.SendMessage(wis.WorkerUUID, recipientWorkerUUID, payload)
}

func (wis *WorkerInstanceServices) StartReceivingMessages(receiverFunc func(message worker.InboundMessage)) {
	wis.Instance.dispatcher.CreateMailbox(wis.WorkerUUID, func(message MailboxMessage) {
		receiverFunc(worker.InboundMessage{
			SourceWorkerUUID: message.SourceWorkerUUID,
			SentTimestamp:    message.SentTimestamp,
			Payload:          message.Payload,
		})
	})
}

func (wis *WorkerInstanceServices) StopReceivingMessages() {
	wis.Instance.dispatcher.RemoveMailbox(wis.WorkerUUID)
}

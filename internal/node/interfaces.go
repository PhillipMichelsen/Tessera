package node

import (
	"AlgorithmicTraderDistributed/internal/node/communication"
	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/google/uuid"
)

type WorkerManager interface {
	RegisterWorker(workerUUID uuid.UUID, worker worker.Worker)
	DeregisterWorker(workerUUID uuid.UUID) error
	StartWorker(workerUUID uuid.UUID, config map[string]interface{}, services worker.Services) error
	StopWorker(workerUUID uuid.UUID) error
	IsWorkerActive(workerUUID uuid.UUID) (bool, error)
}

type Dispatcher interface {
	SendMessage(sourceWorkerUUID uuid.UUID, destinationWorkerUUID uuid.UUID, payload interface{}) error
	CreateMailbox(workerUUID uuid.UUID, receiverFunc func(message communication.IntraNodeMessage))
	RemoveMailbox(workerUUID uuid.UUID)
}

type WorkerFactory interface {
	InstantiateWorker(workerType string, workerUUID uuid.UUID) (worker.Worker, error)
}

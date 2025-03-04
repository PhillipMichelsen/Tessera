package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
	"fmt"
	"github.com/google/uuid"
)

type Instance struct {
	workerManager *WorkerManager
	dispatcher    *Dispatcher

	workerFactory *worker.Factory
}

// NewInstance initializes a new instance.
func NewInstance(workerFactory *worker.Factory) *Instance {
	return &Instance{
		workerManager: NewWorkerManager(),
		dispatcher:    NewDispatcher(),
		workerFactory: workerFactory,
	}
}

// Start starts the instance. A no-op for now.
func (i *Instance) Start() {}

// Stop stops the instance. A no-op for now.
func (i *Instance) Stop() {}

// CreateWorker creates and registers a new worker of the given type.
func (i *Instance) CreateWorker(workerType string, workerUUID uuid.UUID) error {
	instantiatedWorker, err := i.workerFactory.InstantiateWorker(workerType, workerUUID)

	if err != nil {
		fmt.Println(err)
	}

	i.workerManager.RegisterWorker(workerUUID, instantiatedWorker)
	return nil
}

// RemoveWorker removes a worker from the instance.
func (i *Instance) RemoveWorker(workerUUID uuid.UUID) error {
	return i.workerManager.DeregisterWorker(workerUUID)
}

// StartWorker starts a worker with the given configuration.
func (i *Instance) StartWorker(workerUUID uuid.UUID, config map[string]interface{}) error {
	workerServices := NewWorkerServices(i, workerUUID)
	return i.workerManager.StartWorker(workerUUID, config, workerServices)
}

// StopWorker stops a worker.
func (i *Instance) StopWorker(workerUUID uuid.UUID) error {
	return i.workerManager.StopWorker(workerUUID)
}

// IsWorkerActive checks if a worker is active.
func (i *Instance) IsWorkerActive(workerUUID uuid.UUID) (bool, error) {
	return i.workerManager.IsWorkerActive(workerUUID)
}

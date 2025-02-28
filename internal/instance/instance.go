package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
	"fmt"
	"github.com/google/uuid"
)

// Instance represents the main high-level object for the application.
type Instance struct {
	dispatcher    *Dispatcher
	workerManager *WorkerManager

	workerFactory func(workerType string) (worker.Worker, error)
}

// NewInstance initializes a new instance.
func NewInstance(workerFactory func(string) (worker.Worker, error)) *Instance {
	return &Instance{
		dispatcher:    NewDispatcher(),
		workerManager: NewWorkerManager(),
		workerFactory: workerFactory,
	}
}

// Start begins processing for the dispatcher (mailboxes).
func (i *Instance) Start() {
	i.dispatcher.Start()
}

// Stop terminates the dispatcher and stops all active workers.
func (i *Instance) Stop() {
	i.dispatcher.Stop()
	workers := i.workerManager.GetWorkers()
	for workerID := range workers {
		err := i.workerManager.StopWorker(workerID)
		if err != nil {
			return
		}
	}
}

func (i *Instance) SpawnWorker(workerType string, workerUUID uuid.UUID) error {
	instantiatedWorker, err := i.workerFactory(workerType)

	if err != nil {
		return fmt.Errorf("could not create worker %s: %w", workerType, err)
	}

	i.workerManager.AddWorker(workerUUID, instantiatedWorker)

	return nil
}

func (i *Instance) RemoveWorker(workerUUID uuid.UUID) {
	i.workerManager.RemoveWorker(workerUUID)
}

func (i *Instance) StartWorker(workerUUID uuid.UUID, config map[string]interface{}) error {
	err := i.workerManager.StartWorker(workerUUID, config, i.buildWorkerInstanceServices(workerUUID))
	if err != nil {
		return err
	}

	return nil
}

func (i *Instance) buildWorkerInstanceServices(workerUUID uuid.UUID) worker.InstanceServices {
	return &WorkerInstanceServices{
		Instance:   i,
		WorkerUUID: workerUUID,
	}
}

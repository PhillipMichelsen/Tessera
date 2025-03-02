package instance

import "github.com/google/uuid"

type ControllerInterface interface {
	Shutdown() // Shutdown gracefully stops the instance.
	Halt()     // Halt stops the instance immediately.

	CreateWorker(workerType string, workerUUID uuid.UUID) error            // CreateWorker creates a new worker.
	RemoveWorker(workerUUID uuid.UUID) error                               // RemoveWorker removes a worker.
	StartWorker(workerUUID uuid.UUID, config map[string]interface{}) error // StartWorker starts a worker.
	StopWorker(workerUUID uuid.UUID) error                                 // StopWorker stops a worker.
	GetWorkerStatus(workerUUID uuid.UUID) (string, error)                  // GetWorkerStatus returns the status of a worker.
}

type Instance struct {
	workerManager *WorkerManager
	dispatcher    *Dispatcher

	workerRegistry *WorkerRegistry
}

// NewInstance initializes a new instance.
func NewInstance() *Instance {
	return &Instance{
		workerManager:  NewWorkerManager(),
		dispatcher:     NewDispatcher(),
		workerRegistry: NewWorkerRegistry(),
	}
}

func (i *Instance) Start() {}

func (i *Instance) Stop() {}

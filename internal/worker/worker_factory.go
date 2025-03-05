package worker

import (
	"fmt"
	"github.com/google/uuid"
)

type Factory struct {
	workerCreationFunctions map[string]func(uuid uuid.UUID) Worker
}

func NewFactory() *Factory {
	return &Factory{
		workerCreationFunctions: make(map[string]func(uuid uuid.UUID) Worker),
	}
}

func (wf *Factory) RegisterWorkerCreationFunction(workerType string, creationFunc func(uuid uuid.UUID) Worker) {
	wf.workerCreationFunctions[workerType] = creationFunc
}

func (wf *Factory) InstantiateWorker(workerType string, workerUUID uuid.UUID) (Worker, error) {
	creationFunc, exists := wf.workerCreationFunctions[workerType]
	if !exists {
		return nil, fmt.Errorf("worker type %s not registered", workerType)
	}
	return creationFunc(workerUUID), nil
}

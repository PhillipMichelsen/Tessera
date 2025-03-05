package node

import (
	"AlgorithmicTraderDistributed/internal/node/internal"
	"AlgorithmicTraderDistributed/internal/worker"
	"fmt"
	"github.com/google/uuid"
)

// Node represents a node in the distributed system. Typically, a singleton, but not prohibited.
// It is responsible for managing workers and providing Node-level services to them as defined in the worker.Services interface.
// Such services are implemented in the WorkerServices struct.
type Node struct {
	workerManager *internal.WorkerManager
	dispatcher    *internal.Dispatcher

	workerFactory *worker.Factory
}

// NewNode initializes a new node.
func NewNode(workerFactory *worker.Factory) *Node {
	return &Node{
		workerManager: internal.NewWorkerManager(),
		dispatcher:    internal.NewDispatcher(),
		workerFactory: workerFactory,
	}
}

// Start starts the node. A no-op for now.
func (i *Node) Start() {}

// Stop stops the node. A no-op for now.
func (i *Node) Stop() {}

// CreateWorker creates and registers a new worker of the given type.
func (i *Node) CreateWorker(workerType string, workerUUID uuid.UUID) error {
	instantiatedWorker, err := i.workerFactory.InstantiateWorker(workerType, workerUUID)

	if err != nil {
		return fmt.Errorf("failed to instantiate worker: %v", err)
	}

	i.workerManager.RegisterWorker(workerUUID, instantiatedWorker)
	return nil
}

// RemoveWorker removes a worker from the node.
func (i *Node) RemoveWorker(workerUUID uuid.UUID) error {
	return i.workerManager.DeregisterWorker(workerUUID)
}

// StartWorker starts a worker with the given configuration.
func (i *Node) StartWorker(workerUUID uuid.UUID, config map[string]interface{}) error {
	workerServices := NewWorkerServices(i, workerUUID)
	return i.workerManager.StartWorker(workerUUID, config, workerServices)
}

// StopWorker stops a worker.
func (i *Node) StopWorker(workerUUID uuid.UUID) error {
	return i.workerManager.StopWorker(workerUUID)
}

// IsWorkerActive checks if a worker is active.
func (i *Node) IsWorkerActive(workerUUID uuid.UUID) (bool, error) {
	return i.workerManager.IsWorkerActive(workerUUID)
}

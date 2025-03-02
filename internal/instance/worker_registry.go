package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
)

type WorkerRegistry struct {
	workers map[string]worker.Worker
}

func NewWorkerRegistry() *WorkerRegistry {
	return &WorkerRegistry{
		workers: make(map[string]worker.Worker),
	}
}

func (wr *WorkerRegistry) RegisterWorker(worker worker.Worker) {
	wr.workers[worker.GetWorkerName()] = worker
}

func (wr *WorkerRegistry) GetWorker(workerName string) worker.Worker {
	return wr.workers[workerName]
}

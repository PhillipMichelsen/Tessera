package instance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"AlgorithmicTraderDistributed/pkg/worker"

	"github.com/google/uuid"
)

// WorkerManager manages workers and their lifecycle.
type WorkerManager struct {
	workers map[uuid.UUID]*WorkerContainer
	mu      sync.Mutex
}

// WorkerStatus represents the current state of a worker
type WorkerStatus struct {
	isActive  bool
	exitCode  worker.ExitCode
	error     error
	lastStart time.Time
	lastExit  time.Time
}

// WorkerContainer needs to be updated to include status
type WorkerContainer struct {
	uuid       uuid.UUID
	worker     worker.Worker
	status     WorkerStatus
	cancelFunc context.CancelFunc
	done       chan struct{}
}

func NewWorkerManager() *WorkerManager {
	return &WorkerManager{
		workers: make(map[uuid.UUID]*WorkerContainer),
	}
}

// AddWorker registers a new worker with a given uuid.
func (wm *WorkerManager) AddWorker(workerUUID uuid.UUID, worker worker.Worker) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// We will skip the workerUUID collision checks, such 2^128 collisions are best seen as an act of god.

	wc := &WorkerContainer{
		uuid:   workerUUID,
		worker: worker,
		status: WorkerStatus{
			isActive: false,
		},
		cancelFunc: nil, // Cancel func created along with context in StartWorker
	}
	wm.workers[workerUUID] = wc
}

// RemoveWorker removes a worker from the manager.
func (wm *WorkerManager) RemoveWorker(workerUUID uuid.UUID) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	delete(wm.workers, workerUUID)
}

// StartWorker starts a registered worker using its uuid and configuration.
func (wm *WorkerManager) StartWorker(workerUUID uuid.UUID, config map[string]interface{}) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wc, exists := wm.workers[workerUUID]
	if !exists || wc.status.isActive {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wc.cancelFunc = cancel
	wc.done = make(chan struct{})
	wc.status.isActive = true
	wc.status.lastStart = time.Now()
	wc.status.error = nil
	wc.status.exitCode = worker.NormalExit

	go func() {
		defer func() {
			if r := recover(); r != nil {
				wm.handleWorkerExit(workerUUID, worker.PanicExit, fmt.Errorf("%v", r))
			}
		}()

		exitCode, err := wc.worker.Run(ctx, config)
		wm.handleWorkerExit(workerUUID, exitCode, err)
	}()
}

// StopWorker stops a running worker by its uuid.
func (wm *WorkerManager) StopWorker(workerUUID uuid.UUID) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wc, exists := wm.workers[workerUUID]
	if !exists || !wc.status.isActive {
		return // worker not registered or not active
	}

	if wc.cancelFunc == nil {
		panic("Cancel function not created for worker, should never happen under normal circumstances")
		return // Cancel function not created, should never happen
	}

	wc.cancelFunc()
}

func (wm *WorkerManager) GetWorkers() map[uuid.UUID]*WorkerContainer {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.workers
}

func (wm *WorkerManager) handleWorkerExit(workerUUID uuid.UUID, exitCode worker.ExitCode, err error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wc, exists := wm.workers[workerUUID]
	if !exists {
		return
	}

	wc.status.isActive = false
	wc.status.exitCode = exitCode
	wc.status.error = err
	wc.status.lastExit = time.Now()
	wc.cancelFunc = nil
	close(wc.done)

	// TODO: Add logging here
}

func (wm *WorkerManager) blockUntilWorkerExited(workerUUID uuid.UUID) {
	wm.mu.Lock()
	wc, exists := wm.workers[workerUUID]
	wm.mu.Unlock()

	if !exists {
		return
	}

	<-wc.done
}

package node

import (
	"context"
	"fmt"
	"github.com/PhillipMichelsen/Tessera/internal/worker"
	"github.com/google/uuid"
	"sync"
	"time"
)

// WorkerStatus tracks the state of a worker.
type WorkerStatus struct {
	isActive  bool
	exitCode  worker.ExitCode
	error     error
	lastStart time.Time
	lastExit  time.Time
}

// WorkerContainer wraps a worker along with its status and control channels.
type WorkerContainer struct {
	uuid       uuid.UUID
	worker     worker.Worker
	workerType string
	status     WorkerStatus
	services   WorkerServices
	cancelFunc context.CancelFunc
	done       chan struct{}
}

type WorkerFactory interface {
	InstantiateWorker(workerType string) (worker.Worker, error)
}

// Node represents the node which holds and manages workers.
type Node struct {
	dispatcher    *Dispatcher
	workerFactory WorkerFactory
	workers       map[uuid.UUID]*WorkerContainer
	mu            sync.Mutex
}

// NewNode initializes a new Node instance.
func NewNode(workerFactory WorkerFactory) *Node {
	return &Node{
		dispatcher:    NewDispatcher(),
		workerFactory: workerFactory,
		workers:       make(map[uuid.UUID]*WorkerContainer),
	}
}

func (n *Node) ProcessTask(task Task) error {
	for _, instruction := range task.Instructions {
		switch instruction.Type {
		case "create_worker":
			args, ok := instruction.Args.(CreateWorkerInstructionArgs)
			if !ok {
				return fmt.Errorf("failed to decode create_worker args")
			}

			if err := n.createWorker(args.WorkerType, args.WorkerUUID); err != nil {
				return fmt.Errorf("error creating worker: %v", err)
			}

		case "start_worker":
			args, ok := instruction.Args.(StartWorkerInstructionArgs)
			if !ok {
				return fmt.Errorf("failed to decode start_worker args")
			}

			if err := n.startWorker(args.WorkerUUID, args.WorkerRawConfig); err != nil {
				return fmt.Errorf("error starting worker: %v", err)
			}

		case "remove_worker":
			args, ok := instruction.Args.(RemoveWorkerInstructionArgs)
			if !ok {
				return fmt.Errorf("failed to decode remove_worker args")
			}

			if err := n.removeWorker(args.WorkerUUID); err != nil {
				return fmt.Errorf("error removing worker: %v", err)
			}

		case "stop_worker":
			args, ok := instruction.Args.(StopWorkerInstructionArgs)
			if !ok {
				return fmt.Errorf("failed to decode stop_worker args")
			}

			if err := n.stopWorker(args.WorkerUUID); err != nil {
				return fmt.Errorf("error stopping worker: %v", err)
			}

		default:
			return fmt.Errorf("unknown instruction: %s", instruction.Type)
		}
	}

	return nil
}

func (n *Node) ParseTask(yamlBytes []byte) (Task, error) {
	return parseTaskFromYaml(yamlBytes)
}

// createWorker instantiates and registers a new worker.
func (n *Node) createWorker(workerType string, workerUUID uuid.UUID) error {
	instantiatedWorker, err := n.workerFactory.InstantiateWorker(workerType)
	if err != nil {
		return fmt.Errorf("failed to instantiate worker: %v", err)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	// Skipping collision checks, 2^128 collision is too unlikely.
	wc := &WorkerContainer{
		uuid:       workerUUID,
		worker:     instantiatedWorker,
		workerType: workerType,
		status:     WorkerStatus{isActive: false},
	}
	n.workers[workerUUID] = wc

	return nil
}

// removeWorker de-registers a worker. It returns an error if the worker is active.
func (n *Node) removeWorker(workerUUID uuid.UUID) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	wc, exists := n.workers[workerUUID]
	if !exists {
		return fmt.Errorf("worker %s not registered", workerUUID)
	}
	if wc.status.isActive {
		return fmt.Errorf("worker %s is active", workerUUID)
	}

	delete(n.workers, workerUUID)
	return nil
}

// startWorker starts a worker using its configuration and node-provided services.
func (n *Node) startWorker(workerUUID uuid.UUID, rawConfig any) error {
	n.mu.Lock()
	wc, exists := n.workers[workerUUID]
	if !exists || wc.status.isActive {
		n.mu.Unlock()
		return fmt.Errorf("worker %s not registered or already active", workerUUID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wc.cancelFunc = cancel
	wc.done = make(chan struct{})
	wc.status.isActive = true
	wc.status.lastStart = time.Now()
	wc.status.error = nil
	wc.status.exitCode = worker.NormalExit
	wc.services = NewWorkerServices(n)
	n.mu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				n.handleWorkerExit(workerUUID, worker.PanicExit, fmt.Errorf("%v", r))
			}
		}()

		exitCode, err := wc.worker.Run(ctx, rawConfig, wc.services)
		n.handleWorkerExit(workerUUID, exitCode, err)
	}()

	return nil
}

// stopWorker stops a running worker. Blocks until the worker exits.
func (n *Node) stopWorker(workerUUID uuid.UUID) error {
	n.mu.Lock()

	wc, exists := n.workers[workerUUID]
	if !exists || !wc.status.isActive {
		n.mu.Unlock()
		return fmt.Errorf("worker %s not registered or not active", workerUUID)
	}

	if wc.cancelFunc == nil {
		n.mu.Unlock()
		return fmt.Errorf("worker %s has no cancel function", workerUUID)
	}

	wc.cancelFunc()

	n.mu.Unlock()
	<-wc.done

	return nil
}

// handleWorkerExit updates the status of a worker once it exits.
func (n *Node) handleWorkerExit(workerUUID uuid.UUID, exitCode worker.ExitCode, err error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	wc, exists := n.workers[workerUUID]
	if !exists {
		return
	}

	wc.services.cleanupMailboxes()

	wc.status.isActive = false
	wc.status.exitCode = exitCode
	wc.status.error = err
	wc.status.lastExit = time.Now()
	wc.cancelFunc = nil
	wc.services = WorkerServices{}
	close(wc.done)

	fmt.Printf("Worker %s exited with code %d and error: %v\n", workerUUID, exitCode, err)

	// TODO: Implement propagation of worker exit codes to orchestrator.
	// If worker exit code is not NormalExit, propagate to orchestrator.
	// If worker exit code is NormalExit, do not propagate to orchestrator.
}

package node

import (
	"AlgorithmicTraderDistributed/internal/worker"
	"context"
	"fmt"
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
	cancelFunc context.CancelFunc
	done       chan struct{}
}

type WorkerFactory interface {
	InstantiateWorker(workerType string, workerUUID uuid.UUID) (worker.Worker, error)
}

// Assert that Node implements worker.Services.
var _ worker.Services = (*Node)(nil)

// Node represents the node which holds and manages workers.
type Node struct {
	dispatcher *Dispatcher

	workerFactory WorkerFactory

	workers map[uuid.UUID]*WorkerContainer

	mu sync.Mutex
}

// NewNode initializes a new Node instance.
func NewNode(workerFactory *worker.Factory) *Node {
	return &Node{
		dispatcher:    NewDispatcher(),
		workerFactory: workerFactory,
		workers:       make(map[uuid.UUID]*WorkerContainer),
	}
}

// Start is a no-op for now. Will be used to start sub-systems within the node.
func (n *Node) Start() {}

// Stop is a no-op for now. Will be used to stop sub-systems within the node.
func (n *Node) Stop() {}

// CreateWorker instantiates and registers a new worker.
func (n *Node) CreateWorker(workerType string, workerUUID uuid.UUID) error {
	instantiatedWorker, err := n.workerFactory.InstantiateWorker(workerType, workerUUID)
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

// RemoveWorker de-registers a worker. It returns an error if the worker is active.
func (n *Node) RemoveWorker(workerUUID uuid.UUID) error {
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

// StartWorker starts a worker using its configuration and node-provided services.
func (n *Node) StartWorker(workerUUID uuid.UUID, rawConfig any) error {
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
	n.mu.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				n.handleWorkerExit(workerUUID, worker.PanicExit, fmt.Errorf("%v", r))
			}
		}()

		exitCode, err := wc.worker.Run(ctx, rawConfig, n)
		n.handleWorkerExit(workerUUID, exitCode, err)
	}()

	return nil
}

// StopWorker stops a running worker. Blocks until the worker exits.
func (n *Node) StopWorker(workerUUID uuid.UUID) error {
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

// IsWorkerActive checks if a worker is active.
func (n *Node) IsWorkerActive(workerUUID uuid.UUID) (bool, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	wc, exists := n.workers[workerUUID]
	if !exists {
		return false, fmt.Errorf("worker %s not registered", workerUUID)
	}

	return wc.status.isActive, nil
}

func (n *Node) CreateMailbox(mailboxUUID uuid.UUID, receiverFunc func(message worker.Message)) {
	wrapperReceiverFunc := func(message any) {
		receiverFunc(message.(worker.Message))
	}

	n.dispatcher.CreateMailbox(mailboxUUID, wrapperReceiverFunc, 1000)
}

func (n *Node) RemoveMailbox(mailboxUUID uuid.UUID) {
	n.dispatcher.RemoveMailbox(mailboxUUID)
}

func (n *Node) SendMessage(destinationMailboxUUID uuid.UUID, message worker.Message) error {
	// Intra-node message case, can be directly pushed to mailbox.
	if n.dispatcher.CheckMailboxExists(destinationMailboxUUID) {
		return n.dispatcher.PushMessage(destinationMailboxUUID, message)
	}

	// Inter-node message case, needs to be routed through the bridge.
	// First, resolve the destination node. Done with the cluster discovery service.
	// Then, bridge the message to the destination node with the bridge.
	// TODO: Implement the above.

	return fmt.Errorf("unimplemented non intra-node message routing")
}

// handleWorkerExit updates the status of a worker once it exits.
func (n *Node) handleWorkerExit(workerUUID uuid.UUID, exitCode worker.ExitCode, err error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	wc, exists := n.workers[workerUUID]
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
	fmt.Printf("Worker %s exited with code %d and error: %v\n", workerUUID, exitCode, err)
}

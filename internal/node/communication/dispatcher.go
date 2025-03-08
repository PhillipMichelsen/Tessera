package communication

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Dispatcher manages mailboxes and their processing.
// For each mailbox it creates, it spawns a goroutine that
// continuously dequeues messages and passes them to the receiver function.
type Dispatcher struct {
	mu        sync.RWMutex
	mailboxes map[uuid.UUID]chan IntraNodeMessage
	receivers map[uuid.UUID]func(message IntraNodeMessage)
	wg        sync.WaitGroup
}

// NewDispatcher initializes the dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]chan IntraNodeMessage),
		receivers: make(map[uuid.UUID]func(message IntraNodeMessage)),
	}
}

// processMailbox continuously dequeues messages from a mailbox and passes them to its receiver.
func (d *Dispatcher) processMailbox(mailbox chan IntraNodeMessage, receiverFunc func(message IntraNodeMessage)) {
	defer d.wg.Done()

	// Process messages until the channel is closed.
	for msg := range mailbox {
		receiverFunc(msg)
	}
	// Optionally log that mailbox processing for mailboxUUID has ended.
}

// CreateMailbox registers a worker's mailbox with its message handler.
// It creates a new mailbox and spawns a processing goroutine.
func (d *Dispatcher) CreateMailbox(workerUUID uuid.UUID, receiverFunc func(message IntraNodeMessage)) {
	// Create a new mailbox channel.
	mailbox := make(chan IntraNodeMessage)

	d.mu.Lock()
	d.mailboxes[workerUUID] = mailbox
	d.receivers[workerUUID] = receiverFunc
	d.wg.Add(1)
	d.mu.Unlock()

	// Start processing the mailbox in its own goroutine.
	go d.processMailbox(mailbox, receiverFunc)
}

// RemoveMailbox unregisters a worker's mailbox.
// It closes the mailbox so that its processing goroutine can exit.
func (d *Dispatcher) RemoveMailbox(workerID uuid.UUID) {
	d.mu.Lock()
	mailbox, exists := d.mailboxes[workerID]
	if exists {
		delete(d.mailboxes, workerID)
		delete(d.receivers, workerID)
		close(mailbox)
	}
	d.mu.Unlock()
}

// SendMessage queues a message for delivery to the destination worker.
func (d *Dispatcher) SendMessage(sourceWorkerUUID, destinationWorkerUUID uuid.UUID, payload interface{}) error {
	// Safely retrieve the mailbox channel.
	d.mu.RLock()
	mailbox, exists := d.mailboxes[destinationWorkerUUID]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("worker %v does not have a mailbox", destinationWorkerUUID)
	}

	// Build the message.
	msg := IntraNodeMessage{
		SourceWorkerUUID: sourceWorkerUUID,
		SentTimestamp:    time.Now(),
		Payload:          payload,
	}

	// Use select to avoid blocking indefinitely.
	select {
	case mailbox <- msg:
		return nil
	default:
		return fmt.Errorf("mailbox for worker %v is not accepting messages", destinationWorkerUUID)
	}
}

// Wait blocks until all mailbox processing goroutines have exited.
// This is useful for graceful shutdown.
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

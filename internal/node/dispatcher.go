package node

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

// Dispatcher manages mailboxes and their processing.
// For each mailbox it creates, it spawns a goroutine that
// continuously dequeues messages and passes them to the receiver function.
type Dispatcher struct {
	mu        sync.RWMutex
	mailboxes map[uuid.UUID]chan any
	receivers map[uuid.UUID]func(message any)
	wg        sync.WaitGroup
}

// NewDispatcher initializes the dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]chan any),
		receivers: make(map[uuid.UUID]func(message any)),
	}
}

// CreateMailbox registers a worker's mailbox with its message handler.
// It creates a new mailbox and spawns a processing goroutine.
func (d *Dispatcher) CreateMailbox(mailboxUUID uuid.UUID, receiverFunc func(message any), bufferSize int) {
	mailbox := make(chan any, bufferSize)

	d.mu.Lock()
	d.mailboxes[mailboxUUID] = mailbox
	d.receivers[mailboxUUID] = receiverFunc
	d.wg.Add(1)
	d.mu.Unlock()

	// Start processing the mailbox in its own goroutine.
	go d.processMailbox(mailbox, receiverFunc)
}

// RemoveMailbox unregisters a worker's mailbox.
// It closes the mailbox so that its processing goroutine can exit.
func (d *Dispatcher) RemoveMailbox(mailboxUUID uuid.UUID) {
	d.mu.Lock()
	mailbox, exists := d.mailboxes[mailboxUUID]
	if exists {
		delete(d.mailboxes, mailboxUUID)
		delete(d.receivers, mailboxUUID)
		close(mailbox)
	}
	d.mu.Unlock()
}

// PushMessage queues a message for delivery to the destination worker.
func (d *Dispatcher) PushMessage(destinationMailboxUUID uuid.UUID, message any) error {
	// Safely retrieve the mailbox channel.
	d.mu.RLock()
	mailbox, exists := d.mailboxes[destinationMailboxUUID]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("mailbox %v does not exist", destinationMailboxUUID)
	}

	select {
	case mailbox <- message:
		return nil
	default:
		return fmt.Errorf("mailbox %v is full", destinationMailboxUUID)
	}
}

// CheckMailboxExists checks if a mailbox exists.
func (d *Dispatcher) CheckMailboxExists(mailboxUUID uuid.UUID) bool {
	d.mu.RLock()
	_, exists := d.mailboxes[mailboxUUID]
	d.mu.RUnlock()

	return exists
}

// GetMailboxLength returns the number of messages in a worker's mailbox.
func (d *Dispatcher) GetMailboxLength(mailboxUUID uuid.UUID) int {
	d.mu.RLock()
	mailbox, exists := d.mailboxes[mailboxUUID]
	d.mu.RUnlock()

	if !exists {
		return 0
	}

	return len(mailbox)
}

// Wait blocks until all mailbox processing goroutines have exited. Used in graceful shutdown.
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

// processMailbox continuously dequeues messages from a mailbox and passes them to its receiver.
func (d *Dispatcher) processMailbox(mailbox chan any, receiverFunc func(message any)) {
	defer d.wg.Done()

	// Process messages until the channel is closed.
	for msg := range mailbox {
		receiverFunc(msg)
	}
}

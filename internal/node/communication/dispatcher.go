package communication

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
	mailboxes map[uuid.UUID]chan interface{}
	receivers map[uuid.UUID]func(message interface{})
	wg        sync.WaitGroup
}

// NewDispatcher initializes the dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]chan interface{}),
		receivers: make(map[uuid.UUID]func(message interface{})),
	}
}

// CreateMailbox registers a worker's mailbox with its message handler.
// It creates a new mailbox and spawns a processing goroutine.
func (d *Dispatcher) CreateMailbox(mailboxUUID uuid.UUID, receiverFunc func(message interface{}), bufferSize int) {
	mailbox := make(chan interface{}, bufferSize)

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

// SendMessage queues a message for delivery to the destination worker.
func (d *Dispatcher) SendMessage(destinationMailboxUUID uuid.UUID, message interface{}) error {
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

// GetMailboxLength returns the number of messages in a worker's mailbox.
func (d *Dispatcher) GetMailboxLength(workerID uuid.UUID) int {
	d.mu.RLock()
	mailbox, exists := d.mailboxes[workerID]
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
func (d *Dispatcher) processMailbox(mailbox chan interface{}, receiverFunc func(message interface{})) {
	defer d.wg.Done()

	// Process messages until the channel is closed.
	for msg := range mailbox {
		receiverFunc(msg)
	}
}

package instance

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

// Dispatcher manages mailboxes for workers.
type Dispatcher struct {
	mu        sync.RWMutex
	mailboxes map[uuid.UUID]*Mailbox
}

// NewDispatcher initializes the dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]*Mailbox),
	}
}

// Start starts all the mailboxes.
func (d *Dispatcher) Start() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, mailbox := range d.mailboxes {
		mailbox.Start()
	}
}

// Stop stops all mailboxes gracefully.
func (d *Dispatcher) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, mailbox := range d.mailboxes {
		mailbox.Stop()
	}
}

// CreateMailbox registers a worker's mailbox with its message handler.
func (d *Dispatcher) CreateMailbox(workerID uuid.UUID, receiverFunc func(MailboxMessage)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	mb := NewMailbox(receiverFunc)
	d.mailboxes[workerID] = mb
}

// RemoveMailbox unregisters a worker's mailbox.
func (d *Dispatcher) RemoveMailbox(workerID uuid.UUID) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.mailboxes, workerID)
}

// SendMessage queues a message for delivery.
func (d *Dispatcher) SendMessage(sourceWorkerUUID, destinationWorkerUUID uuid.UUID, payload interface{}) error {
	d.mu.RLock()
	mb, exists := d.mailboxes[destinationWorkerUUID]
	d.mu.RUnlock()
	if !exists {
		return fmt.Errorf("worker %s does not have a mailbox", destinationWorkerUUID)
	}

	msg := MailboxMessage{
		SourceWorkerUUID: sourceWorkerUUID,
		SentTimestamp:    time.Now(),
		Payload:          payload,
	}

	mb.PushMessage(msg)
	return nil
}

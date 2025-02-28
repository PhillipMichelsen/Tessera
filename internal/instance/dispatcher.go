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

// SendMessage queues a message for delivery.
func (d *Dispatcher) SendMessage(senderWorkerUUID, receiverWorkerUUID uuid.UUID, payload interface{}) error {
	d.mu.RLock()
	mb, exists := d.mailboxes[receiverWorkerUUID]
	d.mu.RUnlock()
	if !exists {
		return fmt.Errorf("worker %s does not have a mailbox", receiverWorkerUUID)
	}

	msg := MailboxMessage{
		SenderWorkerUUID: senderWorkerUUID,
		SendTimestamp:    time.Now(),
		Payload:          payload,
	}

	mb.PushMessage(msg)
	return nil
}

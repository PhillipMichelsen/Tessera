package instance

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
	mailboxes map[uuid.UUID]*Mailbox
	receivers map[uuid.UUID]func(MailboxMessage)
	wg        sync.WaitGroup
}

// NewDispatcher initializes the dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]*Mailbox),
		receivers: make(map[uuid.UUID]func(MailboxMessage)),
	}
}

// processMailbox continuously dequeues messages from a mailbox and passes them to its receiver.
func (d *Dispatcher) processMailbox(mb *Mailbox, receiverFunc func(MailboxMessage)) {
	defer d.wg.Done()
	for {
		msg, ok := mb.Dequeue()
		if !ok {
			// Mailbox is closed and empty.
			return
		}
		receiverFunc(msg)
	}
}

// CreateMailbox registers a worker's mailbox with its message handler.
// It creates a new mailbox and spawns a processing goroutine.
func (d *Dispatcher) CreateMailbox(workerID uuid.UUID, receiverFunc func(MailboxMessage)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	mb := NewMailbox()
	d.mailboxes[workerID] = mb
	d.receivers[workerID] = receiverFunc
	d.wg.Add(1)
	go d.processMailbox(mb, receiverFunc)
}

// RemoveMailbox unregisters a worker's mailbox.
// It closes the mailbox so that its processing goroutine can exit.
func (d *Dispatcher) RemoveMailbox(workerID uuid.UUID) {
	d.mu.Lock()
	mb, exists := d.mailboxes[workerID]
	if exists {
		mb.Close()
		delete(d.mailboxes, workerID)
		delete(d.receivers, workerID)
	}
	d.mu.Unlock()
}

// SendMessage queues a message for delivery to the destination worker.
func (d *Dispatcher) SendMessage(sourceWorkerUUID, destinationWorkerUUID uuid.UUID, payload interface{}) error {
	d.mu.RLock()
	mb, exists := d.mailboxes[destinationWorkerUUID]
	d.mu.RUnlock()
	if !exists {
		return fmt.Errorf("worker %v does not have a mailbox", destinationWorkerUUID)
	}

	msg := MailboxMessage{
		SourceWorkerUUID: sourceWorkerUUID,
		SentTimestamp:    time.Now(),
		Payload:          payload,
	}

	return mb.PushMessage(msg)
}

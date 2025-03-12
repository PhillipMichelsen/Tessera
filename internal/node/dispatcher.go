package node

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
	"time"
)

// Dispatcher manages mailboxes and their processing.
// For each mailbox it creates, it spawns a goroutine that
// continuously dequeues messages and passes them to the receiver function.
type Dispatcher struct {
	mu         sync.RWMutex
	mailboxes  map[uuid.UUID]chan any
	receivers  map[uuid.UUID]func(message any)
	wg         sync.WaitGroup
	pushCounts map[uuid.UUID]*int64
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
func (d *Dispatcher) CreateMailbox(mailboxUUID uuid.UUID, bufferSize int) {
	mailbox := make(chan any, bufferSize)

	d.mu.Lock()
	d.mailboxes[mailboxUUID] = mailbox
	if d.pushCounts == nil {
		d.pushCounts = make(map[uuid.UUID]*int64)
	}

	d.pushCounts[mailboxUUID] = new(int64)
	d.wg.Add(1)
	d.mu.Unlock()
}

func (d *Dispatcher) GetMailboxChannel(mailboxUUID uuid.UUID) (<-chan any, bool) {
	d.mu.RLock()
	mailbox, exists := d.mailboxes[mailboxUUID]
	d.mu.RUnlock()

	return mailbox, exists
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
	d.mu.RLock()
	mailbox, exists := d.mailboxes[destinationMailboxUUID]
	counter, countExists := d.pushCounts[destinationMailboxUUID]
	d.mu.RUnlock()

	if !exists || !countExists {
		return fmt.Errorf("mailbox %v does not exist", destinationMailboxUUID)
	}

	select {
	case mailbox <- message:
		atomic.AddInt64(counter, 1)
		return nil
	default:
		return fmt.Errorf("mailbox %v is full", destinationMailboxUUID)
	}
}

func (d *Dispatcher) PushMessageBlocking(destinationMailboxUUID uuid.UUID, message any) error {
	d.mu.RLock()
	mailbox, exists := d.mailboxes[destinationMailboxUUID]
	counter, countExists := d.pushCounts[destinationMailboxUUID]
	d.mu.RUnlock()

	if !exists || !countExists {
		return fmt.Errorf("mailbox %v does not exist", destinationMailboxUUID)
	}

	mailbox <- message
	atomic.AddInt64(counter, 1)
	return nil
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

func (d *Dispatcher) LogPushRates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		d.mu.RLock()
		for _, counter := range d.pushCounts {
			// atomically retrieve and reset the counter.
			count := atomic.SwapInt64(counter, 0)
			fmt.Printf("%d, ", count)
		}
		fmt.Println()
		d.mu.RUnlock()
	}
}

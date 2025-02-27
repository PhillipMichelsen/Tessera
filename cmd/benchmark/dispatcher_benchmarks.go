package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Message holds a message's details.
type Message struct {
	Sender    uuid.UUID
	Receiver  uuid.UUID
	Payload   interface{}
	Timestamp time.Time
}

// ----------------------------------------------------------------
// Mailbox: a slice-based variable-sized buffer with its own lock and condition.
// ----------------------------------------------------------------
type Mailbox struct {
	messages []Message  // The message buffer (slice)
	mu       sync.Mutex // Protects the slice
	cond     *sync.Cond // Used to signal arrival of new messages
	// receiverFunc is the function called when a message arrives.
	receiverFunc func(Message)
	// stop channel signals that this mailbox should stop processing.
	stop chan struct{}
}

// NewMailbox creates a new mailbox with the given receiver function.
func NewMailbox(receiverFunc func(Message)) *Mailbox {
	m := &Mailbox{
		messages:     make([]Message, 0),
		receiverFunc: receiverFunc,
		stop:         make(chan struct{}),
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// AppendMessage adds a message to the end of the mailbox.
func (m *Mailbox) AppendMessage(msg Message) {
	m.mu.Lock()
	m.messages = append(m.messages, msg)
	// Signal the condition variable to wake up the waiting goroutine.
	m.cond.Signal()
	m.mu.Unlock()
}

// Start starts the mailbox's dispatcher goroutine.
// Instead of batching messages, it processes one message at a time as soon as it arrives.
func (m *Mailbox) Start() {
	go func() {
		for {
			m.mu.Lock()
			// Wait until there is a message or stop is signaled.
			for len(m.messages) == 0 {
				// Instead of busy waiting or sleeping, we block.
				m.cond.Wait()
				// Check for stop signal.
				select {
				case <-m.stop:
					m.mu.Unlock()
					return
				default:
				}
			}
			// Pop the first message.
			msg := m.messages[0]
			m.messages = m.messages[1:]
			m.mu.Unlock()

			// Immediately call the receiver function.
			m.receiverFunc(msg)
		}
	}()
}

// Stop signals the mailbox to stop processing messages.
func (m *Mailbox) Stop() {
	close(m.stop)
	m.cond.Broadcast()
}

// ----------------------------------------------------------------
// Dispatcher: manages mailboxes and message routing.
// ----------------------------------------------------------------
type Dispatcher struct {
	mu        sync.RWMutex           // Protects the mailboxes map.
	mailboxes map[uuid.UUID]*Mailbox // Map of worker ID to mailbox.
}

// NewDispatcher creates a new Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		mailboxes: make(map[uuid.UUID]*Mailbox),
	}
}

// Start is provided for symmetry (if we want to initialize resources globally).
func (d *Dispatcher) Start() {
	// Nothing extra to do here since each mailbox starts its own goroutine.
}

// Stop stops all mailboxes.
func (d *Dispatcher) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, mb := range d.mailboxes {
		mb.Stop()
	}
}

// CreateMailbox creates a mailbox for the worker with the provided receiver function.
// It calls Start() on the mailbox so that it begins dispatching messages.
func (d *Dispatcher) CreateMailbox(workerID uuid.UUID, receiverFunc func(Message)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	mb := NewMailbox(receiverFunc)
	d.mailboxes[workerID] = mb
	mb.Start()
}

// SendMessage adds a message to the end of the receiver's mailbox.
// It automatically stamps the message with the sender's UUID and the current time.
func (d *Dispatcher) SendMessage(sender, receiver uuid.UUID, payload interface{}) bool {
	d.mu.RLock()
	mb, exists := d.mailboxes[receiver]
	d.mu.RUnlock()
	if !exists {
		return false
	}
	msg := Message{
		Sender:    sender,
		Receiver:  receiver,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	mb.AppendMessage(msg)
	return true
}

// ----------------------------------------------------------------
// Main: Benchmarking the Dispatcher with multiple workers concurrently.
// ----------------------------------------------------------------
func main() {
	numWorkers := 100
	numMessagesPerWorker := 100000

	dispatcher := NewDispatcher()
	dispatcher.Start()

	// Create worker IDs.
	workerIDs := make([]uuid.UUID, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerIDs[i] = uuid.New()
	}

	// Use a WaitGroup to wait for all messages to be processed.
	var wg sync.WaitGroup
	totalMessages := numWorkers * numMessagesPerWorker
	wg.Add(totalMessages)

	// Register a mailbox for each worker. The receiver function just decrements the WaitGroup.
	for i, id := range workerIDs {
		_ = fmt.Sprintf("Worker-%d", i+1)
		receiverFunc := func(msg Message) {
			// For demonstration, print occasionally.
			//fmt.Printf("[%s] Received from %s: %v\n", workerName, msg.Sender, msg.Payload)
			wg.Done()
		}
		dispatcher.CreateMailbox(id, receiverFunc)
	}

	// Benchmark: each worker sends numMessagesPerWorker messages to the next worker (cyclic).
	start := time.Now()

	// Launch a goroutine for each worker to send messages concurrently.
	for i, senderID := range workerIDs {
		receiverID := workerIDs[(i+1)%numWorkers]
		go func(sender, receiver uuid.UUID) {
			for j := 0; j < numMessagesPerWorker; j++ {
				dispatcher.SendMessage(sender, receiver, j)
			}
		}(senderID, receiverID)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Dispatcher processed %d messages among %d workers in %v\n", totalMessages, numWorkers, elapsed)

	dispatcher.Stop()
}

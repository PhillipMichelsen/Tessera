package instance

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

// MailboxMessage represents a worker-to-worker message.
type MailboxMessage struct {
	SenderWorkerUUID uuid.UUID
	SendTimestamp    time.Time
	Payload          interface{}
}

// Mailbox holds messages for a worker.
type Mailbox struct {
	messages     []MailboxMessage
	mu           sync.Mutex
	cond         *sync.Cond
	receiverFunc func(MailboxMessage)
	stop         chan struct{}
	startOnce    sync.Once
}

// NewMailbox initializes a mailbox with a receiver function.
func NewMailbox(receiverFunc func(MailboxMessage)) *Mailbox {
	m := &Mailbox{
		messages:     make([]MailboxMessage, 0),
		receiverFunc: receiverFunc,
		stop:         make(chan struct{}),
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// Start launches the mailbox processing goroutine.
func (m *Mailbox) Start() {
	m.startOnce.Do(func() {
		go m.process()
	})
}

// Stop signals the mailbox to stop processing messages.
func (m *Mailbox) Stop() {
	close(m.stop)
	m.cond.Broadcast()
}

// process is the mailbox's message loop.
func (m *Mailbox) process() {
	for {
		m.mu.Lock()
		// Wait until there's a message or a stop signal.
		for len(m.messages) == 0 {
			select {
			case <-m.stop:
				m.mu.Unlock()
				return
			default:
			}
			m.cond.Wait()
		}
		// Dequeue the first message.
		msg := m.messages[0]
		m.messages = m.messages[1:]
		// Signal potential pushers that a slot is free.
		m.cond.Signal()
		m.mu.Unlock()

		// Process the message outside the lock.
		m.receiverFunc(msg)
	}
}

// PushMessage adds a message to the mailbox and signals the worker.
func (m *Mailbox) PushMessage(msg MailboxMessage) {
	m.mu.Lock()

	select {
	case <-m.stop:
		m.mu.Unlock()
		return
	default:
	}
	m.messages = append(m.messages, msg)
	m.cond.Signal()
	m.mu.Unlock()
}

// GetMessages returns the current slice of messages (for inspection).
func (m *Mailbox) GetMessages() []MailboxMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.messages
}

// GetMessageCount returns the number of messages in the mailbox.
func (m *Mailbox) GetMessageCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.messages)
}

// ClearMailbox removes all messages from the mailbox.
func (m *Mailbox) ClearMailbox() {
	m.mu.Lock()
	m.messages = make([]MailboxMessage, 0)
	m.mu.Unlock()
}

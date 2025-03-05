package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MailboxMessage represents a worker-to-worker message.
type MailboxMessage struct {
	SourceWorkerUUID uuid.UUID
	SentTimestamp    time.Time
	Payload          interface{}
}

// Mailbox is a thread-safe, dynamically sized queue for messages.
type Mailbox struct {
	mu       sync.Mutex
	cond     *sync.Cond
	messages []MailboxMessage
	closed   bool
}

// NewMailbox creates a new Mailbox.
func NewMailbox() *Mailbox {
	m := &Mailbox{
		messages: make([]MailboxMessage, 0),
	}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// PushMessage adds a message to the mailbox.
// Returns an error if the mailbox has been closed.
func (m *Mailbox) PushMessage(msg MailboxMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return fmt.Errorf("mailbox is closed")
	}
	m.messages = append(m.messages, msg)
	m.cond.Signal()
	return nil
}

// Dequeue removes and returns the first message in the mailbox.
// It blocks until a message is available or the mailbox is closed.
// Returns (message, true) if a message was dequeued, or (zero, false)
// if the mailbox is closed and empty.
func (m *Mailbox) Dequeue() (MailboxMessage, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for len(m.messages) == 0 && !m.closed {
		m.cond.Wait()
	}
	if len(m.messages) > 0 {
		// Retrieve and remove the first message.
		msg := m.messages[0]
		// Zero out the first element for garbage collection.
		m.messages[0] = MailboxMessage{}
		m.messages = m.messages[1:]
		return msg, true
	}
	// Mailbox closed and empty.
	var zero MailboxMessage
	return zero, false
}

// GetMessages returns a copy of all messages currently in the mailbox.
func (m *Mailbox) GetMessages() []MailboxMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	cpy := make([]MailboxMessage, len(m.messages))
	copy(cpy, m.messages)
	return cpy
}

// GetMessageCount returns the current number of messages in the mailbox.
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

// Close marks the mailbox as closed and signals all waiting goroutines.
func (m *Mailbox) Close() {
	m.mu.Lock()
	m.closed = true
	m.mu.Unlock()
	m.cond.Broadcast()
}

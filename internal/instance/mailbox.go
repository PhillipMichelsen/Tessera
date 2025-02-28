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
}

// NewMailbox initializes a mailbox with a receiver function.
func NewMailbox(receiverFunc func(MailboxMessage)) *Mailbox {
	m := &Mailbox{
		messages:     make([]MailboxMessage, 0),
		receiverFunc: receiverFunc,
		stop:         make(chan struct{}),
	}
	m.cond = sync.NewCond(&m.mu)

	go func() {
		for {
			m.mu.Lock()
			for len(m.messages) == 0 {
				m.cond.Wait()
				select {
				case <-m.stop:
					m.mu.Unlock()
					return
				default:
				}
			}

			msg := m.messages[0]
			m.messages = m.messages[1:]
			m.mu.Unlock()

			// Invoke the worker's receiver function, sending the message to the worker.
			m.receiverFunc(msg)
		}
	}()

	return m
}

// Start does nothing really...
func (m *Mailbox) Start() {}

// Stop signals the mailbox to stop processing messages.
func (m *Mailbox) Stop() {
	close(m.stop)
	m.cond.Broadcast()
}

// PushMessage adds a message to the mailbox and signals the worker.
func (m *Mailbox) PushMessage(msg MailboxMessage) {
	m.mu.Lock()
	m.messages = append(m.messages, msg)
	m.cond.Signal() // Wake up the receiver goroutine
	m.mu.Unlock()
}

func (m *Mailbox) GetMessages() []MailboxMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.messages
}

func (m *Mailbox) GetMessageCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.messages)
}

func (m *Mailbox) ClearMailbox() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = make([]MailboxMessage, 0)
}

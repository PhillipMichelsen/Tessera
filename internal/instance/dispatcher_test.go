package instance

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDispatcherSendMessage(t *testing.T) {
	dispatcher := NewDispatcher()
	// Use a channel to capture the delivered message.
	msgCh := make(chan MailboxMessage, 1)
	workerID := uuid.New()

	dispatcher.CreateMailbox(workerID, func(msg MailboxMessage) {
		msgCh <- msg
	})
	dispatcher.Start()

	senderID := uuid.New()
	payload := "hello, world"
	err := dispatcher.SendMessage(senderID, workerID, payload)
	if err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}

	// Wait until the message is processed without an explicit sleep.
	select {
	case receivedMsg := <-msgCh:
		if receivedMsg.Payload != payload {
			t.Errorf("Expected payload %v, got %v", payload, receivedMsg.Payload)
		}
		if receivedMsg.SenderWorkerUUID != senderID {
			t.Errorf("Expected sender %v, got %v", senderID, receivedMsg.SenderWorkerUUID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message to be processed")
	}

	dispatcher.Stop()
}

func TestDispatcherSendMessageNonExistent(t *testing.T) {
	dispatcher := NewDispatcher()
	// Attempt to send a message to a worker that has no mailbox.
	err := dispatcher.SendMessage(uuid.New(), uuid.New(), "test")
	if err == nil {
		t.Fatal("Expected error when sending message to non-existent mailbox, got nil")
	}
}

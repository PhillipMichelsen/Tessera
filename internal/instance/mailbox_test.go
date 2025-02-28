package instance

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMailboxProcessing(t *testing.T) {
	// Channel to capture the processed message.
	msgCh := make(chan MailboxMessage, 1)
	mailbox := NewMailbox(func(msg MailboxMessage) {
		msgCh <- msg
	})
	mailbox.Start()

	testMsg := MailboxMessage{
		SenderWorkerUUID: uuid.New(),
		Payload:          "test payload",
		SendTimestamp:    time.Now(),
	}
	mailbox.PushMessage(testMsg)

	// Use a select with timeout to avoid time.Sleep.
	select {
	case received := <-msgCh:
		if received.Payload != testMsg.Payload {
			t.Errorf("Expected payload %v, got %v", testMsg.Payload, received.Payload)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for mailbox to process message")
	}
	mailbox.Stop()
}

func TestMailboxGetAndClear(t *testing.T) {
	// Do not start processing so that messages remain in the queue.
	mailbox := NewMailbox(func(msg MailboxMessage) {
		// no-op receiver
	})

	testMsg1 := MailboxMessage{SenderWorkerUUID: uuid.New(), Payload: "msg1"}
	testMsg2 := MailboxMessage{SenderWorkerUUID: uuid.New(), Payload: "msg2"}

	mailbox.PushMessage(testMsg1)
	mailbox.PushMessage(testMsg2)

	count := mailbox.GetMessageCount()
	if count != 2 {
		t.Errorf("Expected message count 2, got %d", count)
	}
	msgs := mailbox.GetMessages()
	if len(msgs) != 2 {
		t.Errorf("Expected GetMessages to return 2 messages, got %d", len(msgs))
	}

	// Clear the mailbox and verify.
	mailbox.ClearMailbox()
	count = mailbox.GetMessageCount()
	if count != 0 {
		t.Errorf("Expected message count 0 after clearing, got %d", count)
	}
}

func TestMailboxStopPreventsProcessingNoWait(t *testing.T) {
	processedCh := make(chan struct{}, 1)
	mailbox := NewMailbox(func(msg MailboxMessage) {
		processedCh <- struct{}{}
	})
	mailbox.Start()
	mailbox.Stop()

	testMsg := MailboxMessage{SenderWorkerUUID: uuid.New(), Payload: "should not process"}
	mailbox.PushMessage(testMsg)

	// Verify that no message was added to the mailbox.
	if count := mailbox.GetMessageCount(); count != 0 {
		t.Fatalf("Expected 0 messages after Stop, got %d", count)
	}

	// Non-blocking select ensures no processing occurred.
	select {
	case <-processedCh:
		t.Fatal("Message processed after Stop")
	default:
		// Success: nothing received.
	}
}

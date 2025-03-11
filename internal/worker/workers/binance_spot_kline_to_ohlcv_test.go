package workers_test

import (
	"AlgorithmicTraderDistributed/internal/worker/workers"
	"context"
	"sync"
	"testing"
	"time"

	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/worker"

	"github.com/google/uuid"
)

// SentMessageRecord is used to record calls to SendMessage.
type SentMessageRecord struct {
	Destination uuid.UUID
	Message     worker.Message
}

// MockServices implements worker.Services for testing purposes.
type MockServices struct {
	mu           sync.Mutex
	Mailboxes    map[uuid.UUID]func(message worker.Message)
	SentMessages []SentMessageRecord
}

func (ms *MockServices) CreateMailbox(mailboxUUID uuid.UUID, receiverFunc func(message worker.Message)) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if ms.Mailboxes == nil {
		ms.Mailboxes = make(map[uuid.UUID]func(message worker.Message))
	}
	ms.Mailboxes[mailboxUUID] = receiverFunc
}

func (ms *MockServices) RemoveMailbox(mailboxUUID uuid.UUID) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.Mailboxes, mailboxUUID)
}

func (ms *MockServices) SendMessage(destinationMailboxUUID uuid.UUID, message worker.Message) error {
	ms.mu.Lock()
	ms.SentMessages = append(ms.SentMessages, SentMessageRecord{
		Destination: destinationMailboxUUID,
		Message:     message,
	})
	receiver, exists := ms.Mailboxes[destinationMailboxUUID]
	ms.mu.Unlock()
	if exists {
		// Simulate asynchronous delivery.
		go receiver(message)
	}
	return nil
}

func TestBinanceSpotKlineToOHLCVWorker(t *testing.T) {
	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Instantiate the mock services.
	mockServices := &MockServices{}

	// Create UUIDs for input mailbox and destination mailbox.
	inputMailboxUUID := uuid.New()
	destMailboxUUID := uuid.New()

	// Build the configuration.
	config := map[string]any{
		"input_output_mapping": map[string]any{
			"binance": map[string]any{
				"destination_mailbox_uuid": destMailboxUUID.String(),
				"tag":                      "ohlcv",
			},
		},
		"input_mailbox_uuid": inputMailboxUUID.String(),
	}

	// Create the worker.
	moduleUUID := uuid.New()
	workerInstance := workers.NewBinanceSpotKlineToOHLCVWorker(moduleUUID)

	// Run the worker in a separate goroutine.
	exitChan := make(chan worker.ExitCode)
	go func() {
		exitCode, err := workerInstance.Run(ctx, config, mockServices)
		if err != nil {
			t.Errorf("Worker exited with error: %v", err)
		}
		exitChan <- exitCode
	}()

	// Allow a brief moment for the worker to register its mailbox.
	time.Sleep(10 * time.Millisecond)

	// Prepare a sample Binance kline JSON payload.
	// k.t: timestamp, k.o: open, k.h: high, k.l: low, k.c: close, k.v: volume.
	jsonStr := `{"k": {"t": 1618317040000, "o":60000, "h":61000, "l":59000, "c":60500, "v":1200}}`
	serialized := models.SerializedJSON{JSON: jsonStr}

	// Create an incoming message with tag "binance".
	testMessage := worker.Message{
		Tag:     "binance",
		Payload: serialized,
	}

	// Simulate sending the message to the input mailbox.
	mockServices.mu.Lock()
	inputMailboxFunc, exists := mockServices.Mailboxes[inputMailboxUUID]
	mockServices.mu.Unlock()
	if !exists {
		t.Fatalf("Input mailbox with UUID %s not found", inputMailboxUUID)
	}
	inputMailboxFunc(testMessage)

	// Allow time for processing.
	time.Sleep(10 * time.Millisecond)

	// Verify that a message was sent.
	mockServices.mu.Lock()
	if len(mockServices.SentMessages) != 1 {
		mockServices.mu.Unlock()
		t.Fatalf("Expected 1 sent message, got %d", len(mockServices.SentMessages))
	}
	sentRecord := mockServices.SentMessages[0]
	mockServices.mu.Unlock()

	// Check that the message was sent to the correct destination.
	if sentRecord.Destination != destMailboxUUID {
		t.Errorf("Expected destination mailbox UUID %s, got %s", destMailboxUUID, sentRecord.Destination)
	}

	// Verify that the outgoing message has the new tag.
	if sentRecord.Message.Tag != "ohlcv" {
		t.Errorf("Expected message tag 'ohlcv', got '%s'", sentRecord.Message.Tag)
	}

	// Verify the payload is an OHLCV struct with the expected values.
	ohlcv, ok := sentRecord.Message.Payload.(models.OHLCV)
	if !ok {
		t.Fatalf("Expected payload type models.OHLCV, got %T", sentRecord.Message.Payload)
	}
	if ohlcv.Open != 60000 {
		t.Errorf("Expected Open 60000, got %f", ohlcv.Open)
	}
	if ohlcv.High != 61000 {
		t.Errorf("Expected High 61000, got %f", ohlcv.High)
	}
	if ohlcv.Low != 59000 {
		t.Errorf("Expected Low 59000, got %f", ohlcv.Low)
	}
	if ohlcv.Close != 60500 {
		t.Errorf("Expected Close 60500, got %f", ohlcv.Close)
	}
	if ohlcv.Volume != 1200 {
		t.Errorf("Expected Volume 1200, got %f", ohlcv.Volume)
	}
	expectedTimestamp := time.UnixMilli(1618317040000)
	if !ohlcv.Timestamp.Equal(expectedTimestamp) {
		t.Errorf("Expected Timestamp %v, got %v", expectedTimestamp, ohlcv.Timestamp)
	}

	// Cancel the context and wait for the worker to exit.
	cancel()
	select {
	case exitCode := <-exitChan:
		if exitCode != worker.NormalExit {
			t.Errorf("Expected exit code NormalExit, got %d", exitCode)
		}
	case <-time.After(time.Second):
		t.Errorf("Worker did not exit in time")
	}
}

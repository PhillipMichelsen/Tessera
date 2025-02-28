package tests

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"AlgorithmicTraderDistributed/internal/instance"

	"github.com/google/uuid"
)

const (
	numWorkers          = 100
	messagePayloadRange = 1000
)

func BenchmarkDispatcher(b *testing.B) {
	// Use b.N as the number of messages to send.
	numMessages := b.N

	// Initialize the dispatcher.
	dispatcher := instance.NewDispatcher()

	// Worker registration.
	workerIDs := make([]uuid.UUID, numWorkers)
	var wg sync.WaitGroup
	wg.Add(numMessages) // Expect numMessages to be processed.

	// Register workers & their mailboxes.
	for i := 0; i < numWorkers; i++ {
		workerIDs[i] = uuid.New()
		// Each mailbox simply calls wg.Done() upon receiving a message.
		dispatcher.CreateMailbox(workerIDs[i], func(msg instance.MailboxMessage) {
			wg.Done()
		})
	}

	dispatcher.Start()

	// Generate messages outside the timed section.
	messages := make([]int, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = rand.Intn(messagePayloadRange)
	}
	// Create a destination slice sized to numMessages, not numWorkers.
	destinations := make([]uuid.UUID, numMessages)
	for i := 0; i < numMessages; i++ {
		destinations[i] = workerIDs[rand.Intn(len(workerIDs))]
	}

	sendUUID := uuid.New()

	b.ResetTimer()
	start := time.Now()

	// Dispatch messages.
	for i, msg := range messages {
		if err := dispatcher.SendMessage(sendUUID, destinations[i], msg); err != nil {
			b.Fatal(err)
		}
	}

	// Wait until all messages have been processed.
	wg.Wait()
	elapsed := time.Since(start)
	b.StopTimer()

	fmt.Printf("Processed %d messages among %d workers in %v\n", numMessages, numWorkers, elapsed)

	dispatcher.Stop()
}

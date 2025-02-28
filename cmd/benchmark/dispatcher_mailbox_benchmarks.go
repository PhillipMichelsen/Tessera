package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"AlgorithmicTraderDistributed/internal/instance"

	"github.com/google/uuid"
)

const (
	numWorkers          = 100        // Number of workers
	numMessages         = 10_000_000 // Total messages to send
	messagePayloadRange = 1000       // Random payload range
)

func main() {
	// Initialize dispatcher
	dispatcher := instance.NewDispatcher()

	// Worker registration
	workerIDs := make([]uuid.UUID, numWorkers)
	var wg sync.WaitGroup
	wg.Add(numMessages) // We expect numMessages to be processed

	// Register workers & their mailboxes
	for i := 0; i < numWorkers; i++ {
		workerIDs[i] = uuid.New()
		dispatcher.CreateMailbox(workerIDs[i], func(msg instance.MailboxMessage) {
			wg.Done()
		})
	}

	// Generate messages before starting timer
	messages := make([]instance.MailboxMessage, numMessages)
	for i := 0; i < numMessages; i++ {
		sender := workerIDs[rand.Intn(numWorkers)]
		_ = workerIDs[rand.Intn(numWorkers)]
		messages[i] = instance.MailboxMessage{
			SenderWorkerUUID: sender,
			SendTimestamp:    time.Now(),
			Payload:          rand.Intn(messagePayloadRange),
		}
	}

	// Benchmark: Start sending messages
	fmt.Printf("Dispatching %d messages...\n", len(messages))
	dispatcher.Start()
	start := time.Now()

	for _, msg := range messages {
		err := dispatcher.SendMessage(msg.SenderWorkerUUID, msg.SenderWorkerUUID, msg.Payload)
		if err != nil {
			return
		}
	}

	// Wait for all messages to be processed
	wg.Wait()
	elapsed := time.Since(start)

	// Print benchmark results
	fmt.Printf("Processed %d messages among %d workers in %v\n", numMessages, numWorkers, elapsed)

	// Shutdown dispatcher
	dispatcher.Stop()
}

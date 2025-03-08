package communication_test

import (
	"AlgorithmicTraderDistributed/internal/node/communication"
	"github.com/google/uuid"
	"testing"
)

// dummyReceiver is a no-op function that simply consumes the message.
func dummyReceiver(msg communication.IntraNodeMessage) {
	// minimal processing; optionally simulate work with time.Sleep, etc.
}

func BenchmarkDispatcherSend(b *testing.B) {
	// Create a new dispatcher and a mailbox for a single worker.
	d := communication.NewDispatcher()
	workerID := uuid.New()
	d.CreateMailbox(workerID, dummyReceiver)

	// Reset timer to exclude setup time.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SendMessage(workerID, workerID, "benchmark payload", true)
	}
	b.StopTimer()

	d.RemoveMailbox(workerID)
	d.Wait()
}

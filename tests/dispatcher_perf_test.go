package tests

import (
	"AlgorithmicTraderDistributed/internal/node"
	"github.com/google/uuid"
	"testing"
)

// dummyReceiver is a no-op function that simply consumes the message (by returning nothing).
func dummyReceiver(_ interface{}) {
	return
}

func BenchmarkDispatcherSend(b *testing.B) {
	// Create a new dispatcher and a mailbox for a single worker.
	d := node.NewDispatcher()
	workerID := uuid.New()
	d.CreateMailbox(workerID, dummyReceiver, 1000)

	// Reset timer to exclude setup time.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.PushMessage(workerID, "test")
	}
	b.StopTimer()

	d.RemoveMailbox(workerID)
	d.Wait()
}

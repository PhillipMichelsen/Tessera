package communication_test

import (
	"AlgorithmicTraderDistributed/internal/node/communication"
	"github.com/google/uuid"
	"testing"
)

// dummyReceiver is a no-op function that simply consumes the message.
func dummyReceiver(_ interface{}) {
	return
}

func BenchmarkDispatcherSend(b *testing.B) {
	// Create a new dispatcher and a mailbox for a single worker.
	d := communication.NewDispatcher()
	workerID := uuid.New()
	d.CreateMailbox(workerID, dummyReceiver, 1000)

	communication.NewDispatcher()

	// Reset timer to exclude setup time.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SendMessage(workerID, "test")
	}
	b.StopTimer()

	d.RemoveMailbox(workerID)
	d.Wait()
}

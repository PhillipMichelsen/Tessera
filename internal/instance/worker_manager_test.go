package instance

import (
	"AlgorithmicTraderDistributed/pkg/worker"
	"context"
	"testing"

	"github.com/google/uuid"
)

type dummyWorker struct{}

func (d *dummyWorker) Run(ctx context.Context, config map[string]interface{}) (worker.ExitCode, error) {
	exitMode, _ := config["exitMode"].(string) // Get the exit mode from config

	switch exitMode {
	case "panic":
		panic("simulated panic")

	case "selfExit":
		return worker.PrematureExit, nil

	case "waitForStop":
		<-ctx.Done() // Waits for StopWorker() call
		return worker.NormalExit, nil

	default:
		// Default: wait for StopWorker()
		<-ctx.Done()
		return worker.NormalExit, nil
	}
}

func (d *dummyWorker) GetWorkerName() string {
	return "dummyWorker"
}

// TestAddWorker verifies that a worker is added correctly.
func TestAddWorker(t *testing.T) {
	wm := NewWorkerManager()
	workerUUID := uuid.New()
	dw := &dummyWorker{}

	wm.AddWorker(workerUUID, dw)

	wm.mu.Lock()
	_, exists := wm.workers[workerUUID]
	wm.mu.Unlock()

	if !exists {
		t.Fatalf("Worker with UUID %v was not added", workerUUID)
	}
}

// TestSelfExitWorker ensures a worker that exits on its own is marked as "Premature Exit"
func TestSelfExitWorker(t *testing.T) {
	wm := NewWorkerManager()
	workerUUID := uuid.New()
	dw := &dummyWorker{}
	wm.AddWorker(workerUUID, dw)

	wm.StartWorker(workerUUID, map[string]interface{}{
		"exitMode": "selfExit",
	})

	wm.blockUntilWorkerExited(workerUUID)

	wm.mu.Lock()
	defer wm.mu.Unlock()
	wc := wm.workers[workerUUID]
	if wc.status.isActive {
		t.Errorf("Worker should be inactive after self-exit")
	}
	if wc.status.exitCode != worker.PrematureExit {
		t.Errorf("Expected premature exit, got %v", wc.status.exitCode)
	}
}

// TestStopWorker verifies that calling StopWorker results in a "Normal Exit"
func TestStopWorker(t *testing.T) {
	wm := NewWorkerManager()
	workerUUID := uuid.New()
	dw := &dummyWorker{}
	wm.AddWorker(workerUUID, dw)

	wm.StartWorker(workerUUID, map[string]interface{}{
		"exitMode": "waitForStop",
	})

	wm.StopWorker(workerUUID) // Ask it to stop

	wm.blockUntilWorkerExited(workerUUID)

	wm.mu.Lock()
	defer wm.mu.Unlock()
	wc := wm.workers[workerUUID]
	if wc.status.isActive {
		t.Errorf("Worker should be inactive after StopWorker")
	}
	if wc.status.exitCode != worker.NormalExit {
		t.Errorf("Expected normal exit, got %v", wc.status.exitCode)
	}
}

// TestWorkerPanic simulates a worker panic and ensures it is handled correctly.
func TestWorkerPanic(t *testing.T) {
	wm := NewWorkerManager()
	workerUUID := uuid.New()
	dw := &dummyWorker{}
	wm.AddWorker(workerUUID, dw)

	wm.StartWorker(workerUUID, map[string]interface{}{
		"exitMode": "panic",
	})

	wm.blockUntilWorkerExited(workerUUID)

	wm.mu.Lock()
	defer wm.mu.Unlock()
	wc := wm.workers[workerUUID]
	if wc.status.exitCode != worker.PanicExit {
		t.Errorf("Expected panic exit, got %v", wc.status.exitCode)
	}
}

// TestRemoveWorker verifies that a worker is removed correctly.
func TestRemoveWorker(t *testing.T) {
	wm := NewWorkerManager()
	workerUUID := uuid.New()
	dw := &dummyWorker{}
	wm.AddWorker(workerUUID, dw)

	wm.RemoveWorker(workerUUID)

	wm.mu.Lock()
	_, exists := wm.workers[workerUUID]
	wm.mu.Unlock()

	if exists {
		t.Fatalf("Worker with UUID %v was not removed", workerUUID)
	}
}

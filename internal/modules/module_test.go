package modules

import (
	"AlgorithmicTraderDistributed/internal/constants"
	"github.com/google/uuid"
	"sync"
	"testing"
	"time"
)

// TestModuleStatusThreadSafety concurrently reads and writes the module status.
// TestModuleStatusThreadSafety concurrently reads and writes the module status.
func TestModuleStatusThreadSafety(t *testing.T) {
	// Use the MockCore so Start/Stop don’t trigger a panic.
	module := NewModule(uuid.New(), DefaultCoreFactory("MockCore"))

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrently set the status and read it.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				module.setStatus(constants.StartedModuleStatus)
			} else {
				_ = module.GetStatus()
			}
		}(i)
	}

	wg.Wait()

	finalStatus := module.GetStatus()
	if finalStatus != constants.StartedModuleStatus {
		t.Errorf("Unexpected final module status: got %v, expected %v", finalStatus, constants.StartedModuleStatus)
	}
}

// TestConcurrentStartStop simulates multiple concurrent calls to Initialize, Start(), and Stop().
func TestConcurrentStartStop(t *testing.T) {
	// Use the MockCore to avoid immediate panics.
	module := NewModule(uuid.New(), DefaultCoreFactory("MockCore"))

	// Initialize the module before starting it.
	module.Initialize(map[string]interface{}{"dummy": "config"}, nil)

	var wg sync.WaitGroup
	numGoroutines := 50

	// Concurrently call Start() and Stop() on the module.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				module.Start()
			} else {
				module.Stop()
			}
		}(i)
	}

	wg.Wait()

	// Allow some time for all goroutines to settle.
	time.Sleep(500 * time.Millisecond)

	// The module should end in a valid state: either started or stopped.
	finalStatus := module.GetStatus()
	if finalStatus != constants.StartedModuleStatus && finalStatus != constants.StoppedModuleStatus {
		t.Errorf("Module ended in an inconsistent state: %v", finalStatus)
	}
}

// TestCorePanicHandling verifies that a module using a core that panics transitions to the stopped state.
func TestCorePanicHandling(t *testing.T) {
	// Use the MockPanicCore, which panics immediately.
	module := NewModule(uuid.New(), DefaultCoreFactory("MockPanicCore"))

	// Initialize the module before starting it.
	module.Initialize(map[string]interface{}{"dummy": "config"}, nil)

	// Start the module. The core should panic and trigger the panic handling.
	module.Start()

	// Wait a moment for the panic to be recovered and handled.
	time.Sleep(500 * time.Millisecond)

	finalStatus := module.GetStatus()
	if finalStatus != constants.StoppedModuleStatus {
		t.Errorf("Expected module to be stopped after panic handling, but got %v", finalStatus)
	}
}

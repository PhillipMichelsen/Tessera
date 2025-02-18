package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/constants"
	"log"
	"sync"

	"github.com/google/uuid"
)

// TODO: Rework with new instance module liaison system when it is implemented. Remove the worker channels, they are to be implemented in the concrete modules.

type Module interface {
	Initialize()
	Start()
	Stop()
	GetStatus() constants.ModuleStatus
	GetModuleUUID() uuid.UUID
}

type BaseModule struct {
	moduleUUID uuid.UUID
	status     constants.ModuleStatus

	instanceAPI api.InstanceAPI

	stopSignalChannel            chan struct{}

	stopCompletionWaitGroup sync.WaitGroup
	mu                      sync.Mutex
}

// NewBaseModule creates a new BaseModule with directional channels.
func NewBaseModule(
	moduleUUID uuid.UUID,
	outputChannel chan<- interface{},
	updateAlertChannel chan<- uuid.UUID,
) *BaseModule {
	return &BaseModule{
		moduleUUID: moduleUUID,
		status:     constants.Uninitialized,
		mu:                 sync.Mutex{},
	}
}

// Initialize initializes the module. Does not implement the Initialize method of the Module interface, to be implemented by concrete modules that should call this method.
func (b *BaseModule) Initialize(initializationFunc func()) {
	b.setStatus(constants.Initializing)
	log.Printf("[INFO] Module [%s] | State: %s | Module initializing. \n", b.moduleUUID, b.GetStatus())

	defer b.recoverFromPanic()
	initializationFunc()

	b.setStatus(constants.Initialized)
	log.Printf("[INFO] Module [%s] | State: %s | Module initialized succesfully. \n", b.moduleUUID, b.GetStatus())
}

// Start begins the module's worker function. Does not implement the Start method of the Module interface, to be implemented by concrete modules that should call this method.
func (b *BaseModule) Start(runWorker func(api.InstanceAPI, <-chan struct{})) {
	if b.status != constants.Initialized {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not in initialized state, cannot start. \n", b.moduleUUID, b.GetStatus())
		return
	}

	b.setStatus(constants.Starting)
	log.Printf("[INFO] Module [%s] | State: %s | Module starting... \n", b.moduleUUID, b.GetStatus())
	b.stopSignalChannel = make(chan struct{})
	go func() {
		defer b.recoverFromPanic()
		b.stopCompletionWaitGroup.Add(1)
		defer b.stopCompletionWaitGroup.Done()

		runWorker(b.instanceAPI, b.stopSignalChannel)
	}()
	b.setStatus(constants.Started)
	log.Printf("[INFO] Module [%s] | State: %s | Module started succesfully. \n", b.moduleUUID, b.GetStatus())
}

// Stop gracefully stops the module's execution.
func (b *BaseModule) Stop() {
	if b.status != constants.Started {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not in started state, cannot stop. \n", b.moduleUUID, b.GetStatus())
		return
	}

	b.setStatus(constants.Stopping)
	log.Printf("[INFO] Module [%s] | State: %s | Module stopping... \n", b.moduleUUID, b.GetStatus())
	close(b.stopSignalChannel)
	b.stopCompletionWaitGroup.Wait() // Wait for the worker goroutine to finish, should be pretty much instantaneous
	b.setStatus(constants.Stopped)
	log.Printf("[INFO] Module [%s] | State: %s | Module stopped succesfully. \n", b.moduleUUID, b.GetStatus())
}

// GetStatus returns the module's status.
func (b *BaseModule) GetStatus() constants.ModuleStatus {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.status
}

// GetModuleUUID returns the module's UUID.
func (b *BaseModule) GetModuleUUID() uuid.UUID {
	return b.moduleUUID
}

// recoverFromPanic handles recovery from a panic in the module's worker function.
func (b *BaseModule) recoverFromPanic() {
	if r := recover(); r != nil {
		b.setStatus(constants.Error)
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recovery triggered: %v\n", b.moduleUUID, b.GetStatus(), r)
		close(b.stopSignalChannel)
		b.stopCompletionWaitGroup.Wait() // Wait for the worker goroutine to finish
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recoverd, module stopped with error.", b.moduleUUID, b.GetStatus())
	}
}

// setStatus sets the module's status.
func (b *BaseModule) setStatus(newStatus constants.ModuleStatus) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status = newStatus
}

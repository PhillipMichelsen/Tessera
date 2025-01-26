package modules

import (
	"github.com/google/uuid"
	"log"
	"sync"
)

type Module interface {
	Initialize()
	Start()
	Stop()
	GetStatus() ModuleStatus
	GetModuleUUID() uuid.UUID
	GetInputChannel() chan interface{}
}

type ModuleStatus string

const (
	Uninitialized ModuleStatus = "Uninitialized"
	Initializing  ModuleStatus = "Initializing"
	Initialized   ModuleStatus = "Initialized"
	Starting      ModuleStatus = "Starting"
	Started       ModuleStatus = "Started"
	Stopping      ModuleStatus = "Stopping"
	Stopped       ModuleStatus = "Stopped"
	Error         ModuleStatus = "Error"
)

type BaseModule struct {
	moduleUUID uuid.UUID
	status     ModuleStatus

	workerChannels WorkerChannels

	updateAlertChannel      chan<- uuid.UUID
	stopCompletionWaitGroup sync.WaitGroup
	mu                      sync.Mutex
}

type WorkerChannels struct {
	InputChannel  chan interface{}
	OutputChannel chan<- interface{}
	StopSignal    chan struct{}
}

// NewBaseModule creates a new BaseModule with directional channels.
func NewBaseModule(
	moduleUUID uuid.UUID,
	outputChannel chan<- interface{},
	updateAlertChannel chan<- uuid.UUID,
) *BaseModule {
	return &BaseModule{
		moduleUUID: moduleUUID,
		status:     Uninitialized,
		workerChannels: WorkerChannels{
			InputChannel:  make(chan interface{}),
			OutputChannel: outputChannel,
			StopSignal:    make(chan struct{}),
		},
		updateAlertChannel: updateAlertChannel,
		mu:                 sync.Mutex{},
	}
}

// Initialize initializes the module. Does not implement the Initialize method of the Module interface, to be implemented by concrete modules that should call this method.
func (b *BaseModule) Initialize(initializationFunc func()) {
	b.setStatus(Initializing)
	log.Printf("[INFO] Module [%s] | State: %s | Module initializing. \n", b.moduleUUID, b.GetStatus())

	defer b.recoverFromPanic()
	initializationFunc()

	b.setStatus(Initialized)
	b.sendUpdateAlert()
	log.Printf("[INFO] Module [%s] | State: %s | Module initialized succesfully. \n", b.moduleUUID, b.GetStatus())
}

// Start begins the module's worker function. Does not implement the Start method of the Module interface, to be implemented by concrete modules that should call this method.
func (b *BaseModule) Start(runWorker func(WorkerChannels)) {
	if b.status != Initialized {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not in initialized state, cannot start. \n", b.moduleUUID, b.GetStatus())
		return
	}

	b.setStatus(Starting)
	log.Printf("[INFO] Module [%s] | State: %s | Module starting... \n", b.moduleUUID, b.GetStatus())
	b.workerChannels.StopSignal = make(chan struct{})
	go func() {
		defer b.recoverFromPanic()
		b.stopCompletionWaitGroup.Add(1)
		defer b.stopCompletionWaitGroup.Done()

		runWorker(b.workerChannels)
	}()
	b.setStatus(Started)
	b.sendUpdateAlert()
	log.Printf("[INFO] Module [%s] | State: %s | Module started succesfully. \n", b.moduleUUID, b.GetStatus())
}

// Stop gracefully stops the module's execution.
func (b *BaseModule) Stop() {
	if b.status != Started {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not in started state, cannot stop. \n", b.moduleUUID, b.GetStatus())
		return
	}

	b.setStatus(Stopping)
	log.Printf("[INFO] Module [%s] | State: %s | Module stopping... \n", b.moduleUUID, b.GetStatus())
	close(b.workerChannels.StopSignal)
	b.stopCompletionWaitGroup.Wait() // Wait for the worker goroutine to finish, should be pretty much instantaneous
	b.setStatus(Stopped)
	b.sendUpdateAlert()
	log.Printf("[INFO] Module [%s] | State: %s | Module stopped succesfully. \n", b.moduleUUID, b.GetStatus())
}

// GetStatus returns the module's status.
func (b *BaseModule) GetStatus() ModuleStatus {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.status
}

// GetInputChannel returns the module's input channel.
func (b *BaseModule) GetInputChannel() chan<- interface{} {
	return b.workerChannels.InputChannel
}

// GetModuleUUID returns the module's UUID.
func (b *BaseModule) GetModuleUUID() uuid.UUID {
	return b.moduleUUID
}

// recoverFromPanic handles recovery from a panic in the module's worker function.
func (b *BaseModule) recoverFromPanic() {
	if r := recover(); r != nil {
		b.setStatus(Error)
		b.sendUpdateAlert()
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recovery triggered: %v\n", b.moduleUUID, b.GetStatus(), r)
		close(b.workerChannels.StopSignal)
		b.stopCompletionWaitGroup.Wait() // Wait for the worker goroutine to finish
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recoverd, module stopped with error.", b.moduleUUID, b.GetStatus())
	}
}

// sendUpdateAlert sends an update alert to the update alert channel. This is used to notify the instance that the module's status has changed.
func (b *BaseModule) sendUpdateAlert() {
	b.updateAlertChannel <- b.moduleUUID
}

// setStatus sets the module's status.
func (b *BaseModule) setStatus(newStatus ModuleStatus) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.status = newStatus
}

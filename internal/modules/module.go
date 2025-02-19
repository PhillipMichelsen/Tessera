package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/constants"
	"log"
	"sync"

	"github.com/google/uuid"
)

type ModuleHandler interface {
	Initialize(rawConfig map[string]interface{}) error
	Run(instanceAPI api.InstanceAPIInternal, runtimeErrorReceiver func(error))
	Stop() error
}

type Module struct {
	moduleUUID uuid.UUID
	status     constants.ModuleStatus
	handler    ModuleHandler

	instanceAPI       api.InstanceAPIInternal
	stopSignalChannel chan struct{}
	workerWaitGroup   sync.WaitGroup
	mu                sync.Mutex
}

func NewModule(moduleUUID uuid.UUID, worker ModuleHandler, instanceAPI api.InstanceAPIInternal) *Module {
	return &Module{
		moduleUUID:  moduleUUID,
		status:      constants.Uninitialized,
		handler:     worker,
		instanceAPI: instanceAPI,
		mu:          sync.Mutex{},
	}
}

func (m *Module) Initialize(config map[string]interface{}) {
	m.setStatus(constants.Initializing)
	log.Printf("[INFO] Module [%s] | State: %s | Loading configuration...\n", m.moduleUUID, m.GetStatus())

	defer m.recoverFromPanic()
	err := m.handler.Initialize(config)

	if err != nil {
		m.setStatus(constants.Error)
		log.Printf("[ERROR] Module [%s] | State: %s | Error: %v\n", m.moduleUUID, m.GetStatus(), err)
		return
	}

	m.setStatus(constants.Initialized)
	log.Printf("[INFO] Module [%s] | State: %s | Configuration loaded.\n", m.moduleUUID, m.GetStatus())
}

func (m *Module) Start() {
	if m.status != constants.Initialized {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not initialized, cannot start.\n", m.moduleUUID, m.GetStatus())
		return
	}

	m.setStatus(constants.Starting)
	log.Printf("[INFO] Module [%s] | State: %s | Starting...\n", m.moduleUUID, m.GetStatus())

	m.stopSignalChannel = make(chan struct{})
	m.workerWaitGroup.Add(1)

	go func() {
		defer m.recoverFromPanic()
		defer m.workerWaitGroup.Done()
		m.handler.Run(m.instanceAPI, m.receiveRuntimeError)
	}()

	m.setStatus(constants.Started)
	log.Printf("[INFO] Module [%s] | State: %s | Started successfully.\n", m.moduleUUID, m.GetStatus())
}

func (m *Module) Stop() {
	if m.status != constants.Started {
		log.Printf("[WARNING] Module [%s] | State: %s | Module not running, cannot stop.\n", m.moduleUUID, m.GetStatus())
		return
	}

	m.setStatus(constants.Stopping)
	log.Printf("[INFO] Module [%s] | State: %s | Stopping...\n", m.moduleUUID, m.GetStatus())

	close(m.stopSignalChannel)
	m.workerWaitGroup.Wait()
	err := m.handler.Stop()

	if err != nil {
		m.setStatus(constants.Error)
		log.Printf("[ERROR] Module [%s] | State: %s | Error: %v\n", m.moduleUUID, m.GetStatus(), err)
		return
	}

	m.setStatus(constants.Stopped)
	log.Printf("[INFO] Module [%s] | State: %s | Stopped successfully.\n", m.moduleUUID, m.GetStatus())
}

func (m *Module) GetStatus() constants.ModuleStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *Module) GetModuleUUID() uuid.UUID {
	return m.moduleUUID
}

func (m *Module) setStatus(newStatus constants.ModuleStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = newStatus
}

func (m *Module) receiveRuntimeError(err error) {
	if err != nil {
		m.setStatus(constants.Error)
		log.Printf("[ERROR] Module [%s] | State: %s | Error: %v\n", m.moduleUUID, m.GetStatus(), err)
		if m.stopSignalChannel != nil {
			close(m.stopSignalChannel)
		}
		m.workerWaitGroup.Wait()
		log.Printf("[ERROR] Module [%s] | State: %s | Module stopped with error.\n", m.moduleUUID, m.GetStatus())
	}
}

func (m *Module) recoverFromPanic() {
	if r := recover(); r != nil {
		m.setStatus(constants.Error)
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recovery triggered: %v\n", m.moduleUUID, m.GetStatus(), r)
		if m.stopSignalChannel != nil {
			close(m.stopSignalChannel)
		}
		m.workerWaitGroup.Wait()
		log.Printf("[CRITICAL] Module [%s] | State: %s | Panic recovered, module stopped with error.\n", m.moduleUUID, m.GetStatus())
	}
}

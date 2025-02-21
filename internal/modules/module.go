package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/constants"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Module struct {
	moduleUUID uuid.UUID
	status     constants.ModuleStatus
	core       Core

	instanceServicesAPI api.InstanceServicesAPI

	stopSignalChannel chan struct{}
	coreWaitGroup     sync.WaitGroup
	mu                sync.Mutex
}

func NewModule(moduleUUID uuid.UUID, core Core) *Module {
	return &Module{
		moduleUUID: moduleUUID,
		status:     constants.Uninitialized,
		core:       core,
		mu:         sync.Mutex{},
	}
}

func (m *Module) Initialize(config map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) {
	m.instanceServicesAPI = instanceServicesAPI
	m.setStatus(constants.Initializing)
	m.sendLog(zerolog.DebugLevel, "Initializing module...", nil)

	err := m.core.Initialize(config, m.instanceServicesAPI)

	if err != nil {
		m.setStatus(constants.Error)
		m.sendLog(zerolog.ErrorLevel, "Error initializing core", err)

		return
	}

	m.setStatus(constants.Initialized)
	m.sendLog(zerolog.DebugLevel, "Initialized module successfully!", nil)
}

func (m *Module) Start() {
	if m.status != constants.Initialized {
		m.sendLog(zerolog.WarnLevel, "Cannot start module that is not initialized", nil)
		return
	}

	m.setStatus(constants.Starting)
	m.sendLog(zerolog.DebugLevel, "Starting module...", nil)
	

	m.stopSignalChannel = make(chan struct{})
	m.coreWaitGroup.Add(1)

	go func() {
		defer m.coreWaitGroup.Done()
		m.core.Run()
	}()

	m.setStatus(constants.Started)
	m.sendLog(zerolog.DebugLevel, "Started module successfully!", nil)
}

func (m *Module) Stop() {
	if m.status != constants.Started {
		m.sendLog(zerolog.WarnLevel, "Cannot stop module that is not started", nil)
		return
	}

	m.setStatus(constants.Stopping)
	m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)

	close(m.stopSignalChannel)
	m.coreWaitGroup.Wait()
	err := m.core.Stop()

	if err != nil {
		m.setStatus(constants.Error)
		m.sendLog(zerolog.ErrorLevel, "Error stopping core", err)
		return
	}

	m.setStatus(constants.Stopped)
	m.sendLog(zerolog.DebugLevel, "Stopped module successfully!", nil)
}

func (m *Module) GetModuleUUID() uuid.UUID {
	return m.moduleUUID
}

func (m *Module) GetStatus() constants.ModuleStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *Module) setStatus(newStatus constants.ModuleStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = newStatus
}

func (m *Module) sendLog(level zerolog.Level, message string, err error) {
	log.WithLevel(level).Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetCoreType()).Str("module_status", string(m.status)).Err(err).Msg(message)
}
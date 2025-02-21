package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/constants"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Module struct {
	moduleUUID uuid.UUID
	status     constants.ModuleStatus
	core       Core

	instanceAPIInternal api.InstanceAPIInternal

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

func (m *Module) Initialize(config map[string]interface{}) {
	m.setStatus(constants.Initializing)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Initializing module...")

	err := m.core.Initialize(config)

	if err != nil {
		m.setStatus(constants.Error)
		log.Error().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Err(err).Msg("Error initializing module")

		return
	}

	m.setStatus(constants.Initialized)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Initialized module successfully!")
}

func (m *Module) Start() {
	if m.status != constants.Initialized {
		log.Warn().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Cannot start module that is not initialized")
		return
	}

	m.setStatus(constants.Starting)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Starting module...")

	m.stopSignalChannel = make(chan struct{})
	m.coreWaitGroup.Add(1)

	go func() {
		defer m.coreWaitGroup.Done()
		m.core.Run()
	}()

	m.setStatus(constants.Started)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Started module successfully!")
}

func (m *Module) Stop() {
	if m.status != constants.Started {
		log.Warn().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Cannot stop module that is not started")
		return
	}

	m.setStatus(constants.Stopping)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Stopping module...")

	close(m.stopSignalChannel)
	m.coreWaitGroup.Wait()
	err := m.core.Stop()

	if err != nil {
		m.setStatus(constants.Error)
		log.Error().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Err(err).Msg("Error stopping core")
		return
	}

	m.setStatus(constants.Stopped)
	log.Debug().Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetType()).Str("module_status", string(m.status)).Msg("Stopped module successfully!")
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

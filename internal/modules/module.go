package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/constants"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Module struct {
	moduleUUID         uuid.UUID
	status             constants.ModuleStatus
	timeOfStatusChange time.Time

	coreContainer CoreContainer

	instanceServicesAPI api.InstanceServicesAPI

	mu sync.Mutex
}

type CoreContainer struct {
	core     Core
	config   map[string]interface{}
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewModule(moduleUUID uuid.UUID, core Core) *Module {
	return &Module{
		moduleUUID: moduleUUID,
		status:     constants.UninitializedModuleStatus,
		coreContainer: CoreContainer{
			core: core,
		},
		mu: sync.Mutex{},
	}
}

func (m *Module) Initialize(coreConfig map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) {
	if m.status == constants.StartedModuleStatus || m.status == constants.StartingModuleStatus {
		m.sendLog(zerolog.WarnLevel, "Cannot initialize module that is already started or starting", nil)
		return
	}

	if m.status == constants.StoppingModuleStatus {
		m.sendLog(zerolog.WarnLevel, "Cannot initialize module that is stopping", nil)
		return
	}

	m.instanceServicesAPI = instanceServicesAPI
	m.setStatus(constants.InitializingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Initializing module...", nil)

	m.coreContainer.config = coreConfig

	m.setStatus(constants.InitializedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Initialized module successfully!", nil)
}

func (m *Module) Start() {
	if m.status != constants.InitializedModuleStatus {
		m.sendLog(zerolog.WarnLevel, "Cannot start module that is not initialized", nil)
		return
	}

	m.setStatus(constants.StartingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Starting module...", nil)

	ctx, cancel := context.WithCancel(context.Background())
	m.coreContainer.ctx = ctx
	m.coreContainer.cancel = cancel

	m.coreContainer.wg.Add(1)

	m.runCore()

	m.setStatus(constants.StartedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Started module successfully!", nil)
}

func (m *Module) Stop() {
	if m.status != constants.StartedModuleStatus {
		m.sendLog(zerolog.WarnLevel, "Cannot stop module that is not started", nil)
		return
	}

	m.stopModule(true)
}

func (m *Module) GetModuleUUID() uuid.UUID {
	return m.moduleUUID
}

func (m *Module) GetStatus() constants.ModuleStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *Module) stopModule(stopCore bool) {
	m.setStatus(constants.StoppingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)

	if stopCore {
		m.stopCore()
		m.coreContainer.wg.Wait()
	}

	m.setStatus(constants.StoppedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Stopped module successfully!", nil)
}

func (m *Module) runCore() {
	m.sendLog(zerolog.DebugLevel, "Running core in goroutine...", nil)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				m.sendLog(zerolog.ErrorLevel, "Core's Run function panicked! Recovered and sent to core exit handler.", err)
				m.handleCoreExit(err)
			}
		}()
		err := m.coreContainer.core.Run(m.coreContainer.ctx, m.coreContainer.config, m.instanceServicesAPI)
		m.handleCoreExit(err)
	}()
	m.sendLog(zerolog.DebugLevel, "Core running in goroutine!", nil)
}

func (m *Module) stopCore() {
	m.sendLog(zerolog.DebugLevel, "Stopping core...", nil)
	m.coreContainer.cancel()
}

func (m *Module) handleCoreExit(err error) {
	m.coreContainer.wg.Done()

	if err == nil {
		m.sendLog(zerolog.DebugLevel, "Core has exited successfully.", nil)

		if m.status != constants.StoppingModuleStatus {
			m.sendLog(zerolog.DebugLevel, "Core's exit was self-initiated. Will stop module.", nil)
			m.stopModule(false)
		}
	} else {
		m.sendLog(zerolog.WarnLevel, "Core exited with error. Will stop module.", err)
		m.stopModule(false)
	}
}

func (m *Module) setStatus(newStatus constants.ModuleStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = newStatus
	m.timeOfStatusChange = time.Now()
}

func (m *Module) sendLog(level zerolog.Level, message string, err error) {
	log.WithLevel(level).Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.coreContainer.core.GetCoreName()).Str("module_status", string(m.status)).Err(err).Msg(message)
}

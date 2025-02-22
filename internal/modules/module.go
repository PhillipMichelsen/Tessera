package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/constants"
	"context"
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
	core   Core
	config map[string]interface{}
	ctx    context.Context
	cancel context.CancelFunc
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

func (m *Module) Initialize(config map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) {
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

	m.coreContainer.config = config

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

	m.coreContainer.core.Run(ctx, m.coreContainer.config, m.receiveCoreError, m.instanceServicesAPI)

	m.setStatus(constants.StartedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Started module successfully!", nil)
}

func (m *Module) Stop() {
	if m.status != constants.StartedModuleStatus {
		m.sendLog(zerolog.WarnLevel, "Cannot stop module that is not started", nil)
		return
	}

	m.setStatus(constants.StoppingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)

	m.coreContainer.cancel()

	m.setStatus(constants.StoppedModuleStatus)
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

func (m *Module) receiveCoreError(err error) {
	m.setStatus(constants.ErrorModuleStatus)
	m.sendLog(zerolog.ErrorLevel, "Received error from core, stopping and reporting to instance...", err)

	m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)
	m.core.Stop()
	m.sendLog(zerolog.DebugLevel, "Stopped module successfully!", nil)

	m.instanceServicesAPI.ReceiveModuleError(m.moduleUUID, err)
}

func (m *Module) setStatus(newStatus constants.ModuleStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = newStatus
	m.timeOfStatusChange = time.Now()
}

func (m *Module) sendLog(level zerolog.Level, message string, err error) {
	log.WithLevel(level).Str("module_uuid", m.moduleUUID.String()).Str("core_type", m.core.GetCoreTypeName()).Str("module_status", string(m.status)).Err(err).Msg(message)
}

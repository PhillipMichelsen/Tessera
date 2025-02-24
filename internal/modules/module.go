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

// Module encapsulates lifecycle management for a core.
type Module struct {
	moduleUUID         uuid.UUID // immutable after construction, no lock needed
	status             constants.ModuleStatus
	timeOfStatusChange time.Time

	coreContainer CoreContainer

	instanceServicesAPI api.InstanceServicesAPI

	mu sync.Mutex // protects mutable fields: status, timeOfStatusChange, instanceServicesAPI, and coreContainer fields when needed
}

// CoreContainer bundles the core with its runtime context.
type CoreContainer struct {
	core   Core
	config map[string]interface{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewModule creates a new Module with the provided UUID and core.
func NewModule(moduleUUID uuid.UUID, core Core) *Module {
	return &Module{
		moduleUUID: moduleUUID,
		status:     constants.UninitializedModuleStatus,
		coreContainer: CoreContainer{
			core: core,
		},
	}
}

// Initialize sets up the module. It is only allowed if the module is not already running or stopping.
func (m *Module) Initialize(coreConfig map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) {
	m.mu.Lock()
	currentStatus := m.status
	if currentStatus == constants.StartedModuleStatus || currentStatus == constants.StartingModuleStatus || currentStatus == constants.StoppingModuleStatus {
		m.mu.Unlock()
		m.sendLog(zerolog.WarnLevel, fmt.Sprintf("Cannot initialize module in status %s", currentStatus), nil)
		return
	}
	m.mu.Unlock()

	m.instanceServicesAPI = instanceServicesAPI

	m.setStatus(constants.InitializingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Initializing module...", nil)

	m.coreContainer.config = coreConfig

	m.setStatus(constants.InitializedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Initialized module successfully!", nil)
}

// Start launches the core if the module is initialized.
func (m *Module) Start() {
	m.mu.Lock()
	if m.status != constants.InitializedModuleStatus {
		m.mu.Unlock()
		m.sendLog(zerolog.WarnLevel, "Cannot start module that is not initialized", nil)
		return
	}
	m.mu.Unlock()

	m.setStatus(constants.StartingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Starting module...", nil)

	ctx, cancel := context.WithCancel(context.Background())
	m.mu.Lock()
	m.coreContainer.ctx = ctx
	m.coreContainer.cancel = cancel
	m.mu.Unlock()

	m.runCore()

	m.setStatus(constants.StartedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Started module successfully!", nil)
}

// Stop initiates a graceful shutdown if the module is started.
func (m *Module) Stop() {
	m.mu.Lock()
	if m.status != constants.StartedModuleStatus {
		m.mu.Unlock()
		m.sendLog(zerolog.WarnLevel, "Cannot stop module that is not started", nil)
		return
	}
	m.mu.Unlock()

	m.stopModule(true)
}

// GetModuleUUID returns the module's immutable UUID.
func (m *Module) GetModuleUUID() uuid.UUID {
	return m.moduleUUID
}

// GetStatus returns the module's current status in a thread-safe way.
func (m *Module) GetStatus() constants.ModuleStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

// stopModule coordinates the shutdown of the module and its core.
func (m *Module) stopModule(stopCore bool) {
	m.setStatus(constants.StoppingModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)

	if stopCore {
		m.stopCore() // now blocking until core fully stops
	}

	m.setStatus(constants.StoppedModuleStatus)
	m.sendLog(zerolog.DebugLevel, "Stopped module successfully!", nil)
}

// runCore starts the core in its own goroutine with proper error and panic handling.
func (m *Module) runCore() {
	m.sendLog(zerolog.DebugLevel, "Starting core in goroutine...", nil)

	m.coreContainer.wg.Add(1)
	go func() {
		defer m.coreContainer.wg.Done()
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
				m.sendLog(zerolog.ErrorLevel, "Core panicked; recovered in goroutine.", err)
			}
			m.handleCoreExit(err)
		}()
		err = m.coreContainer.core.Run(m.coreContainer.ctx, m.coreContainer.config, m.instanceServicesAPI)
	}()

	m.sendLog(zerolog.DebugLevel, "Core goroutine started.", nil)
}

// stopCore cancels the core context and blocks until the core goroutine has exited.
func (m *Module) stopCore() {
	m.sendLog(zerolog.DebugLevel, "Stopping core...", nil)
	m.mu.Lock()
	cancel := m.coreContainer.cancel
	m.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	m.coreContainer.wg.Wait()
}

// handleCoreExit deals with core termination, logging the result and triggering module shutdown if needed.
func (m *Module) handleCoreExit(err error) {
	if err == nil {
		m.sendLog(zerolog.DebugLevel, "Core exited successfully.", nil)
		m.mu.Lock()
		isStopping := m.status == constants.StoppingModuleStatus
		m.mu.Unlock()
		if !isStopping {
			m.sendLog(zerolog.DebugLevel, "Core exited on its own; initiating module stop.", nil)
			m.stopModule(false)
		}
	} else {
		m.sendLog(zerolog.WarnLevel, "Core exited with error; initiating module stop.", err)
		m.stopModule(false)
	}
}

// setStatus updates the module's status and records the change time.
func (m *Module) setStatus(newStatus constants.ModuleStatus) {
	m.mu.Lock()
	m.status = newStatus
	m.timeOfStatusChange = time.Now()
	m.mu.Unlock()
}

// sendLog enriches log messages with module context.
func (m *Module) sendLog(level zerolog.Level, message string, err error) {
	// Immutable values can be read without locking.
	uuidStr := m.moduleUUID.String()
	coreName := m.coreContainer.core.GetCoreName()

	// For current status, we grab it in a thread-safe way.
	m.mu.Lock()
	currentStatus := m.status
	m.mu.Unlock()

	log.WithLevel(level).
		Str("module_uuid", uuidStr).
		Str("core_type", coreName).
		Str("module_status", string(currentStatus)).
		Err(err).
		Msg(message)
}

package modules

import (
    "AlgorithmicTraderDistributed/internal/api"
    "context"
    "fmt"
    "sync"

    "github.com/google/uuid"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Module encapsulates lifecycle management for a core.
type Module struct {
    moduleUUID uuid.UUID
    isActive   bool
    isStopping bool

    coreContainer CoreContainer

    mu sync.Mutex
}

// CoreContainer bundles the core with its runtime context.
type CoreContainer struct {
    core      Core
    config    map[string]interface{}
    ctx       context.Context
    cancel    context.CancelFunc
    wg        sync.WaitGroup
}

// NewModule creates a new Module with the provided UUID and core.
func NewModule(moduleUUID uuid.UUID, core Core) *Module {
    return &Module{
        moduleUUID: moduleUUID,
        isActive:   false,
        isStopping: false,
        coreContainer: CoreContainer{
            core: core,
        },
    }
}

// Start launches the core if the module is initialized.
func (m *Module) Start(config map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) {
    m.mu.Lock()
    
    // Check if module can be started
    if m.isActive {
        m.sendLog(zerolog.WarnLevel, "Cannot start module that is already active", nil)
        m.mu.Unlock()
        return
    }
    
    if m.isStopping {
        m.sendLog(zerolog.WarnLevel, "Cannot start module that is currently stopping", nil)
        m.mu.Unlock()
        return
    }
    
    m.sendLog(zerolog.DebugLevel, "Starting module...", nil)
    
    // Set up context and configuration
    ctx, cancel := context.WithCancel(context.Background())
    m.coreContainer.ctx = ctx
    m.coreContainer.cancel = cancel
    m.coreContainer.config = config
    
    // Mark as active before launching goroutine
    m.isActive = true
    m.mu.Unlock()
    
    // Launch core in a goroutine
    m.coreContainer.wg.Add(1)
    go func() {
        defer m.coreContainer.wg.Done()
        
        var err error
        defer func() {
            // Handle panics and natural exits
            if r := recover(); r != nil {
                err = fmt.Errorf("panic: %v", r)
                
                m.mu.Lock()
                m.sendLog(zerolog.ErrorLevel, "Core panicked; recovered in calling goroutine", err)
                m.mu.Unlock()
            }
            
            m.mu.Lock()
            // Only attempt to stop if we're still active and not already stopping
            shouldInitiateStop := m.isActive && !m.isStopping
            
            if err == nil {
                m.sendLog(zerolog.DebugLevel, "Core exited without an error", nil)
            } else {
                m.sendLog(zerolog.ErrorLevel, "Core exited with an error", err)
            }
            m.mu.Unlock()
            
            // If core exited naturally and module is still marked as active,
            // we need to stop the module
            if shouldInitiateStop {
                m.Stop()
            }
        }()
        
        // Run the core
        err = m.coreContainer.core.Run(m.coreContainer.ctx, m.coreContainer.config, instanceServicesAPI)
    }()
    
    m.mu.Lock()
    m.sendLog(zerolog.DebugLevel, "Started module successfully", nil)
    m.mu.Unlock()
}

// Stop initiates a graceful shutdown if the module is started.
func (m *Module) Stop() {
    m.mu.Lock()
    
    // Check if module can be stopped
    if !m.isActive {
        m.sendLog(zerolog.WarnLevel, "Cannot stop module that is not started", nil)
        m.mu.Unlock()
        return
    }
    
    if m.isStopping {
        m.sendLog(zerolog.WarnLevel, "Module is already stopping", nil)
        m.mu.Unlock()
        return
    }
    
    // Mark as stopping to prevent concurrent stops
    m.isStopping = true
    m.sendLog(zerolog.DebugLevel, "Stopping module...", nil)
    m.sendLog(zerolog.DebugLevel, "Stopping core...", nil)
    
    // Save cancel function for use outside the lock
    cancelFunc := m.coreContainer.cancel
    m.mu.Unlock()
    
    // Cancel context without holding the lock
    cancelFunc()
    
    // Wait without holding the lock to avoid deadlocks
    m.coreContainer.wg.Wait()
    
    // Update state after core is fully stopped
    m.mu.Lock()
    m.isActive = false
    m.isStopping = false
    m.sendLog(zerolog.DebugLevel, "Module stopped successfully", nil)
    m.mu.Unlock()
}

// GetModuleUUID returns the module's immutable UUID.
func (m *Module) GetModuleUUID() uuid.UUID {
    // No need for locking since this is an immutable value
    return m.moduleUUID
}

// IsActive returns the module's current status in a thread-safe way.
func (m *Module) IsActive() bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    return m.isActive
}

// IsStopping returns whether the module is in the process of stopping.
func (m *Module) IsStopping() bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    return m.isStopping
}

// sendLog enriches log messages with module context.
// Must be called with the mutex held.
func (m *Module) sendLog(level zerolog.Level, message string, err error) {
    // Since this is called with the lock held, we can safely access state
    log.WithLevel(level).
        Str("module_uuid", m.moduleUUID.String()).
        Str("core_type", m.coreContainer.core.GetCoreName()).
        Bool("module_isactive", m.isActive).
        Bool("module_isstopping", m.isStopping).
        Err(err).
        Msg(message)
}
package instance

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/constants"
	"AlgorithmicTraderDistributed/internal/common/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"os"
	"time"

	"github.com/google/uuid"
)

type Instance struct {
	instanceUUID uuid.UUID

	dispatcher       *Dispatcher
	controlInterface *ControlInterface

	moduleFactoryFunction func(string, uuid.UUID) *modules.Module
	modules               map[uuid.UUID]ModuleInfo
}

type ModuleInfo struct {
	ModuleAPI    api.ModuleAPI
	CreationTime time.Time
}

func NewInstance() *Instance {
	return &Instance{
		instanceUUID:          uuid.New(),
		dispatcher:            NewDispatcher(100),
		//controlInterface:      NewControlInterface(), // TODO: Add control interface
		moduleFactoryFunction: modules.InstantiateModule,
	}
}

// EXTERNAL API METHODS

// CreateModule creates a new module with the given name and UUID. Uses the provided module factory function to create the modules.
func (i *Instance) CreateModule(moduleName string, moduleUUID uuid.UUID) {
	module := i.moduleFactoryFunction(moduleName, moduleUUID)

	i.modules[module.GetModuleUUID()] = ModuleInfo{
		ModuleAPI:    module,
		CreationTime: time.Now(),
	}
}

// RemoveModule stops and removes the module with the given UUID.
func (i *Instance) RemoveModule(moduleUUID uuid.UUID) {
	i.DeregisterModuleInputChannel(moduleUUID)
	delete(i.modules, moduleUUID)
}

// InitializeModule initializes the module with the given UUID with the provided configuration.
func (i *Instance) InitializeModule(moduleUUID uuid.UUID, config map[string]interface{}) {
	i.modules[moduleUUID].ModuleAPI.Initialize(config)
}

// StartModule starts the module with the given UUID.
func (i *Instance) StartModule(moduleUUID uuid.UUID) {
	i.modules[moduleUUID].ModuleAPI.Start()
}

// StopModule stops the module with the given UUID.
func (i *Instance) StopModule(moduleUUID uuid.UUID) {
	i.modules[moduleUUID].ModuleAPI.Stop()
}

// Halt ungracefully shuts down the instance. Does not call Stop modules. 
func (i *Instance) Halt() {
	os.Exit(1)
}

// Shutdown gracefully shuts down the instance. Calls Stop on all modules.
func (i *Instance) Shutdown() {
	for _, moduleInfo := range i.modules {
		moduleInfo.ModuleAPI.Stop()
	}
	os.Exit(0)
}

// GetModules returns a list of all module UUIDs.
func (i *Instance) GetModules() []uuid.UUID {
	moduleUUIDs := make([]uuid.UUID, 0, len(i.modules))
	for moduleUUID := range i.modules {
		moduleUUIDs = append(moduleUUIDs, moduleUUID)
	}
	return moduleUUIDs
}

// GetModuleStatus returns the status of the module with the given UUID.
func (i *Instance) GetModuleStatus(moduleUUID uuid.UUID) constants.ModuleStatus {
	return i.modules[moduleUUID].ModuleAPI.GetStatus()
}

// INTERNAL API METHODS

// DispatchPacket dispatches the given packet with the dispatcher.
func (i *Instance) DispatchPacket(packet models.Packet) {
	i.dispatcher.Dispatch(packet)
}

// RegisterModuleInputChannel registers the given input channel with the dispatcher for the module with the given UUID.
func (i *Instance) RegisterModuleInputChannel(moduleUUID uuid.UUID, inputChannel chan interface{}) {
	i.dispatcher.AddMapping(moduleUUID, inputChannel)
}


func (i *Instance) DeregisterModuleInputChannel(moduleUUID uuid.UUID) {
	i.dispatcher.RemoveMapping(moduleUUID)
}

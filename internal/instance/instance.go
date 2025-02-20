package instance

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/constants"
	"AlgorithmicTraderDistributed/internal/common/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"github.com/google/uuid"
	"log"
	"os"
)

type Instance struct {
	instanceUUID uuid.UUID
	controllers  []Controller

	moduleFactoryFunction func(string, uuid.UUID, api.InstanceAPIInternal) *modules.Module
	modules               map[uuid.UUID]*ModuleRecord

	packetDispatchQueue      chan *models.Packet
	packetDispatchStopSignal chan struct{}
}

type Controller interface {
	Start()
	Stop()
}

type ModuleRecord struct {
	ModuleAPI          api.ModuleAPI
	ModuleInputChannel chan *models.Packet
}

// NewInstance initializes a new instance with internal message routing.
func NewInstance() *Instance {
	instance := &Instance{
		instanceUUID:             uuid.New(),
		moduleFactoryFunction:    modules.InstantiateModule,
		modules:                  make(map[uuid.UUID]*ModuleRecord),
		packetDispatchQueue:      make(chan *models.Packet, 100),
		packetDispatchStopSignal: make(chan struct{}),
	}

	// Start internal message processing
	go instance.packetDispatchWorker()

	return instance
}

// EXTERNAL API METHODS

// CreateModule creates a new module and registers it for communication.
func (i *Instance) CreateModule(moduleName string, moduleUUID uuid.UUID) {
	module := i.moduleFactoryFunction(moduleName, moduleUUID, i)

	i.modules[module.GetModuleUUID()] = &ModuleRecord{
		ModuleAPI:          module,
		ModuleInputChannel: nil,
	}
}

func (i *Instance) AddModule(module *modules.Module) {
	i.modules[module.GetModuleUUID()] = &ModuleRecord{
		ModuleAPI:          module,
		ModuleInputChannel: nil,
	}
}

// RemoveModule stops and removes the module.
func (i *Instance) RemoveModule(moduleUUID uuid.UUID) {
	i.UnregisterModuleInputChannel(moduleUUID)
	delete(i.modules, moduleUUID)
}

// InitializeModule initializes the module with the given configuration.
func (i *Instance) InitializeModule(moduleUUID uuid.UUID, config map[string]interface{}) {
	i.modules[moduleUUID].ModuleAPI.Initialize(config)
}

// StartModule starts the module.
func (i *Instance) StartModule(moduleUUID uuid.UUID) {
	i.modules[moduleUUID].ModuleAPI.Start()
}

// StopModule stops the module.
func (i *Instance) StopModule(moduleUUID uuid.UUID) {
	i.modules[moduleUUID].ModuleAPI.Stop()
}

// Halt ungracefully shuts down the instance.
func (i *Instance) Halt() {
	os.Exit(1)
}

// Shutdown gracefully stops all modules and shuts down the instance.
func (i *Instance) Shutdown() {
	for _, module := range i.modules {
		module.ModuleAPI.Stop()
	}
	close(i.packetDispatchStopSignal)
	os.Exit(0)
}

// GetModules returns all module UUIDs.
func (i *Instance) GetModules() []uuid.UUID {
	moduleUUIDs := make([]uuid.UUID, 0, len(i.modules))
	for moduleUUID := range i.modules {
		moduleUUIDs = append(moduleUUIDs, moduleUUID)
	}
	return moduleUUIDs
}

// GetModuleStatus returns the module's status.
func (i *Instance) GetModuleStatus(moduleUUID uuid.UUID) constants.ModuleStatus {
	return i.modules[moduleUUID].ModuleAPI.GetStatus()
}

// INTERNAL API METHODS

// DispatchPacket enqueues a packet for module-to-module communication.
func (i *Instance) DispatchPacket(packet *models.Packet) {
	i.packetDispatchQueue <- packet
}

// RegisterModuleInputChannel registers a module's input channel for communication.
func (i *Instance) RegisterModuleInputChannel(moduleUUID uuid.UUID, inputChannel chan *models.Packet) {
	i.modules[moduleUUID].ModuleInputChannel = inputChannel
}

// UnregisterModuleInputChannel unregisters a module's input channel.
func (i *Instance) UnregisterModuleInputChannel(moduleUUID uuid.UUID) {
	i.modules[moduleUUID].ModuleInputChannel = nil
}

func (i *Instance) packetDispatchWorker() {
	for {
		select {
		case <-i.packetDispatchStopSignal:
			return
		case packet := <-i.packetDispatchQueue:
			i.dispatchMessage(packet)
		}
	}
}

// dispatchMessage routes packets to their respective modules.
func (i *Instance) dispatchMessage(packet *models.Packet) {
	destinationModuleChannel := i.modules[packet.DestinationModuleUUID].ModuleInputChannel
	if destinationModuleChannel == nil {
		log.Printf("Destination module [%s] not found.\n", packet.DestinationModuleUUID)
		return
	}

	destinationModuleChannel <- packet
}

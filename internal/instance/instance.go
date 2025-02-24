package instance

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/constants"
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Instance struct {
	instanceUUID uuid.UUID
	controllers  []Controller

	modules map[uuid.UUID]*ModuleContainer

	packetDispatchQueue      chan *models.Packet
	packetDispatchStopSignal chan struct{}
}

type ModuleContainer struct {
	ModuleControlAPI   api.ModuleControlAPI
	ModuleInputChannel chan *models.Packet
}

// NewInstance initializes a new instance with internal message routing.
func NewInstance() *Instance {
	instance := &Instance{
		instanceUUID:             uuid.New(),
		modules:                  make(map[uuid.UUID]*ModuleContainer),
		packetDispatchQueue:      make(chan *models.Packet, 100),
		packetDispatchStopSignal: make(chan struct{}),
	}
	return instance
}

// EXTERNAL API METHODS

// CreateModule creates a new module and registers it for communication.
// TODO: Add support for plugin cores. Currently only uses the DefaultCoreFactor method exposed by modules.
func (i *Instance) CreateModule(coreName string, moduleUUID uuid.UUID) {
	module := modules.NewModule(
		moduleUUID,
		modules.DefaultCoreFactory(coreName),
	)

	i.modules[module.GetModuleUUID()] = &ModuleContainer{
		ModuleControlAPI:   module,
		ModuleInputChannel: nil,
	}

	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Created module with [%s] core and moduleUUID [%s]", coreName, moduleUUID), nil)
}

// RemoveModule stops and removes the module.
func (i *Instance) RemoveModule(moduleUUID uuid.UUID) {
	i.UnregisterModuleInputChannel(moduleUUID)
	delete(i.modules, moduleUUID)

	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Removed module with moduleUUID [%s]", moduleUUID), nil)
}

// InitializeModule initializes the module with the given configuration.
func (i *Instance) InitializeModule(moduleUUID uuid.UUID, config map[string]interface{}) {
	i.modules[moduleUUID].ModuleControlAPI.Initialize(config, i)

	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Initialized module with moduleUUID [%s]", moduleUUID), nil)
}

// StartModule starts the module.
func (i *Instance) StartModule(moduleUUID uuid.UUID) {
	if i.modules[moduleUUID].ModuleControlAPI.GetStatus() == constants.StartedModuleStatus {
		log.Warn().Str("instance_uuid", i.instanceUUID.String()).Msg(fmt.Sprintf("Cannot start module [%s] as it is already started.", moduleUUID))
		return
	}
	i.modules[moduleUUID].ModuleControlAPI.Start()

	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Started module with moduleUUID [%s]", moduleUUID), nil)
}

// StopModule stops the module.
func (i *Instance) StopModule(moduleUUID uuid.UUID) {
	if i.modules[moduleUUID].ModuleControlAPI.GetStatus() != constants.StartedModuleStatus {
		log.Warn().Str("instance_uuid", i.instanceUUID.String()).Msg(fmt.Sprintf("Cannot stop module [%s] as it is not started.", moduleUUID))
		return
	}
	i.modules[moduleUUID].ModuleControlAPI.Stop()

	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Stopped module with moduleUUID [%s]", moduleUUID), nil)
}

// Halt ungracefully shuts down the instance.
func (i *Instance) Halt() {
	i.sendLog(zerolog.ErrorLevel, "Halting instance...", nil)
	os.Exit(444)
}

// Shutdown gracefully stops all modules and shuts down the instance.
func (i *Instance) Shutdown() {
	i.sendLog(zerolog.InfoLevel, "Shutting down instance...", nil)

	for _, module := range i.modules {
		if module.ModuleControlAPI.GetStatus() == constants.StartedModuleStatus {
			i.StopModule(module.ModuleControlAPI.GetModuleUUID())
		}
	}
	close(i.packetDispatchStopSignal)

	for _, controller := range i.controllers {
		controller.Stop()
	}

	i.sendLog(zerolog.InfoLevel, "Instance shut down successfully.", nil)
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
	return i.modules[moduleUUID].ModuleControlAPI.GetStatus()
}

// GetInstanceUUID returns the instance's UUID.
func (i *Instance) GetInstanceUUID() uuid.UUID {
	return i.instanceUUID
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

func (i *Instance) SignalStatusUpdate(moduleUUID uuid.UUID) {
	i.sendLog(zerolog.InfoLevel, fmt.Sprintf("Received status update signal from module with moduleUUID [%s], new status: %s", moduleUUID, i.GetModuleStatus(moduleUUID)), nil)
}

// NON-API METHODS

func (i *Instance) Start() {
	for _, controller := range i.controllers {
		controller.Start()
	}

	go i.packetDispatchWorker()

	i.sendLog(zerolog.InfoLevel, "Instance started successfully.", nil)
}

func (i *Instance) AddController(controller Controller) {
	i.controllers = append(i.controllers, controller)
}

// INTERNAL METHODS

// packetDispatchWorker listens for packets and dispatches them to their respective modules. To be run as a goroutine.
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
		i.sendLog(zerolog.ErrorLevel, fmt.Sprintf("Failed to dispatch packet to module with moduleUUID [%s]: module not found", packet.DestinationModuleUUID), nil)
		return
	}

	destinationModuleChannel <- packet
}

func (i *Instance) sendLog(level zerolog.Level, message string, err error) {
	log.WithLevel(level).Str("instance_uuid", i.instanceUUID.String()).Err(err).Msg(message)
}

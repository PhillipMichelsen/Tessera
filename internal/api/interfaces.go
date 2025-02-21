package api

import (
	"AlgorithmicTraderDistributed/internal/constants"
	"AlgorithmicTraderDistributed/internal/models"
	"github.com/google/uuid"
)

type InstanceAPIExternal interface {
	CreateModule(moduleName string, moduleUUID uuid.UUID)
	RemoveModule(moduleUUID uuid.UUID)
	InitializeModule(moduleUUID uuid.UUID, config map[string]interface{})
	StartModule(moduleUUID uuid.UUID)
	StopModule(moduleUUID uuid.UUID)

	Halt()     // Ungraceful shutdown, for emergency purposes.
	Shutdown() // Graceful shutdown.

	GetModules() []uuid.UUID
	GetModuleStatus(moduleUUID uuid.UUID) constants.ModuleStatus

	GetInstanceUUID() uuid.UUID
}

type InstanceAPIInternal interface {
	DispatchPacket(packet *models.Packet)
	RegisterModuleInputChannel(moduleUUID uuid.UUID, inputChannel chan *models.Packet)
	UnregisterModuleInputChannel(moduleUUID uuid.UUID)
	ReceiveRuntimeErrorAlert(moduleUUID uuid.UUID, err error)
}

type ModuleAPI interface {
	Initialize(map[string]interface{})
	Start()
	Stop()
	GetModuleUUID() uuid.UUID
	GetStatus() constants.ModuleStatus
}

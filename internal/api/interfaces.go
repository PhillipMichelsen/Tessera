package api

import (
	"AlgorithmicTraderDistributed/internal/common/constants"
	"AlgorithmicTraderDistributed/internal/common/models"
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
}

type InstanceAPIInternal interface {
	DispatchPacket(packet *models.Packet)
	RegisterModuleInputChannel(moduleUUID uuid.UUID, inputChannel chan *models.Packet)
	UnregisterModuleInputChannel(moduleUUID uuid.UUID)
}

type ModuleAPI interface {
	Initialize(map[string]interface{})
	Start()
	Stop()
	GetStatus() constants.ModuleStatus
	GetModuleUUID() uuid.UUID
}

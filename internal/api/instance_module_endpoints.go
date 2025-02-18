package api

import (
	"AlgorithmicTraderDistributed/internal/common/constants"
	"AlgorithmicTraderDistributed/internal/common/models"
	"github.com/google/uuid"
)

type ModuleAPI interface {
	Initialize(map[string]interface{}) error
	Start() error
	Stop() error
	GetStatus() constants.ModuleStatus
	GetModuleUUID() uuid.UUID
}

type InstanceAPI interface {
	DispatchPacket(packet models.Packet)
	RegisterModuleInputChannel(inputChannel chan interface{})
	DeregisterModuleInputChannel()
}

package api

import (
	"AlgorithmicTraderDistributed/internal/common/constants"
	"AlgorithmicTraderDistributed/internal/common/models"
	"AlgorithmicTraderDistributed/internal/modules"

	"github.com/google/uuid"
)

type ModuleAPI interface {
	modules.Module
}

type InstanceAPI interface {
	DispatchPacket(packet models.Packet)
	RegisterModuleInputChannel(inputChannel chan interface{})
	DeregisterModuleInputChannel()
}
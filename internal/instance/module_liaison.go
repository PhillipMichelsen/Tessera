package instance

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
)

type ModuleLiaison struct {
	instance *Instance
	module   modules.Module
}

func NewModuleLiaison() *ModuleLiaison {
	return &ModuleLiaison{}
}

func (m *ModuleLiaison) RequestPacketDispatch(packet models.Packet) {
	m.instance.DispatchPacket(packet)
}

func (m *ModuleLiaison) RegisterModuleInputChannel(inputChannel chan interface{}) {
	m.instance.RegisterModuleInputChannel(m.module.GetModuleUUID(), inputChannel)
}

func (m *ModuleLiaison) DeregisterModuleInputChannel() {
	m.instance.DeregisterModuleInputChannel(m.module.GetModuleUUID())
}

func (m *ModuleLiaison) ReportStatusUpdate() {
	// TODO: Implement status updates, to escalate to the instance where it can be handled
	return
}

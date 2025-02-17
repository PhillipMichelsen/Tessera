package instance

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"github.com/google/uuid"
	"time"
)

type Instance struct {
	instanceUUID uuid.UUID

	dispatcher       *Dispatcher
	controlInterface *ControlInterface

	moduleFactoryFunction func(string, uuid.UUID) modules.Module
	modules               map[uuid.UUID]ModuleInfo
}

type ModuleInfo struct {
	Module        modules.Module
	ModuleLiaison *ModuleLiaison
	CreationTime  time.Time
	CreatedBy     uuid.UUID
}

func NewInstance() *Instance {
	return &Instance{
		instanceUUID:     uuid.New(),
		dispatcher:       NewDispatcher(100),
		controlInterface: NewControlInterface(),
	}
}

func (i *Instance) Start() {
	i.dispatcher.Start()
	i.controlInterface.Start()
}

func (i *Instance) Stop() {
	i.dispatcher.Stop()
	i.controlInterface.Stop()
}

func (i *Instance) DispatchPacket(packet models.Packet) {
	i.dispatcher.Dispatch(packet)
}

func (i *Instance) NewModule(moduleName string, moduleUUID uuid.UUID) {
	moduleLiaison := NewModuleLiaison()
	module := i.moduleFactoryFunction(moduleName, moduleUUID, moduleLiaison)

	i.modules[module.GetModuleUUID()] = ModuleInfo{
		Module:        module,
		ModuleLiaison: moduleLiaison,
		CreationTime:  time.Now(),
		CreatedBy:     i.instanceUUID,
	}
}

func (i *Instance) RemoveModule(moduleUUID uuid.UUID) {
	i.DeregisterModuleInputChannel(moduleUUID)
	delete(i.modules, moduleUUID)
}

func (i *Instance) RegisterModuleInputChannel(moduleUUID uuid.UUID, inputChannel chan interface{}) {
	i.dispatcher.AddMapping(moduleUUID, inputChannel)
}

func (i *Instance) DeregisterModuleInputChannel(moduleUUID uuid.UUID) {
	i.dispatcher.RemoveMapping(moduleUUID)
}

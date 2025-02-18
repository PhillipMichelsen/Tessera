package instance

import (
	"AlgorithmicTraderDistributed/internal/api"
	"AlgorithmicTraderDistributed/internal/common/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"time"

	"github.com/google/uuid"
)

type Instance struct {
	instanceUUID uuid.UUID

	dispatcher       *Dispatcher
	controlInterface *ControlInterface

	moduleFactoryFunction func(string, uuid.UUID) modules.Module
	modules               map[uuid.UUID]ModuleInfo
}

type ModuleInfo struct {
	ModuleAPI    api.ModuleAPI
	CreationTime time.Time
	CreatedBy    uuid.UUID
}

func NewInstance() *Instance {
	return &Instance{
		instanceUUID:     uuid.New(),
		dispatcher:       NewDispatcher(100),
		controlInterface: NewControlInterface(),
		moduleFactoryFunction: modules.CreateModule,
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
	module := i.moduleFactoryFunction(moduleName, moduleUUID)

	i.modules[module.GetModuleUUID()] = ModuleInfo{
		ModuleAPI: module,
		CreationTime:   time.Now(),
		CreatedBy:      i.instanceUUID,
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

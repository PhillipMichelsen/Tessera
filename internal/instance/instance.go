package instance

type Instance struct {
	dispatcher       *Dispatcher
	moduleLiaison    *ModuleLiaison
	moduleManager    *ModuleManager
	controlInterface *ControlInterface
}

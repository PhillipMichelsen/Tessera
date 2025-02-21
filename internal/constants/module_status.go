package constants

type ModuleStatus string

const (
	UninitializedModuleStatus ModuleStatus = "UninitializedModuleStatus"
	InitializingModuleStatus  ModuleStatus = "InitializingModuleStatus"
	InitializedModuleStatus   ModuleStatus = "InitializedModuleStatus"
	StartingModuleStatus      ModuleStatus = "StartingModuleStatus"
	StartedModuleStatus       ModuleStatus = "StartedModuleStatus"
	StoppingModuleStatus      ModuleStatus = "StoppingModuleStatus"
	StoppedModuleStatus       ModuleStatus = "StoppedModuleStatus"
	ErrorModuleStatus         ModuleStatus = "ErrorModuleStatus"
)

package constants

type ModuleStatus string

const (
	Uninitialized ModuleStatus = "Uninitialized"
	Initializing  ModuleStatus = "Initializing"
	Initialized   ModuleStatus = "Initialized"
	Starting      ModuleStatus = "Starting"
	Started       ModuleStatus = "Started"
	Stopping      ModuleStatus = "Stopping"
	Stopped       ModuleStatus = "Stopped"
	Error         ModuleStatus = "Error"
)

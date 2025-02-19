package modules

import "github.com/google/uuid"

func InstantiateModule(moduleName string, moduleUUID uuid.UUID) *Module {
	switch moduleName {
	case "":
		return &Module{}
	default:
		return &Module{}
	}
}

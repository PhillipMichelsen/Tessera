package modules

import "github.com/google/uuid"

func CreateModule(moduleName string, moduleUUID uuid.UUID) *Module {
	switch moduleName {
	case "":
		return nil
	default:
		return nil
	}
}

package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	"github.com/google/uuid"
)

func InstantiateModule(moduleName string, moduleUUID uuid.UUID, instanceAPI api.InstanceAPIInternal) *Module {
	switch moduleName {
	case "":
		return &Module{}
	default:
		return &Module{}
	}
}

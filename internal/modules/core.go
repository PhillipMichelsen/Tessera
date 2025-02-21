package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	cores_concrete "AlgorithmicTraderDistributed/internal/modules/cores"
)

type Core interface {
	Initialize(rawConfig map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) error
	Run()
	Stop() error
	GetCoreType() string
}

func InstantiateCoreByName(coreName string) Core {
	switch coreName {
	case "TestCore":
		return &cores_concrete.TestCore{}
	}

	return nil
}

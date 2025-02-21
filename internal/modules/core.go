package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	misccores "AlgorithmicTraderDistributed/internal/modules/cores/misc"
)

type Core interface {
	Initialize(rawConfig map[string]interface{}, coreErrorReceiver func(error), instanceServicesAPI api.InstanceServicesAPI) error
	Run()
	Stop()
	GetCoreType() string
}

func InstantiateCoreByName(coreName string) Core {
	switch coreName {
	case "TestCore":
		return &misccores.TestCore{}
	}

	return nil
}

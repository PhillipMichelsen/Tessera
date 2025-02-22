package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	misccores "AlgorithmicTraderDistributed/internal/modules/cores/misc"
	"context"
)

type Core interface {
	Run(ctx context.Context, coreConfig map[string]interface{}, coreErrorReceiver func(error), instance api.InstanceServicesAPI)
	GetCoreType() string
}

func InstantiateCoreByName(coreName string) Core {
	switch coreName {
	case "TestCore":
		return &misccores.TestCore{}
	}

	return nil
}

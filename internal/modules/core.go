package modules

import (
	"AlgorithmicTraderDistributed/internal/api"
	misccores "AlgorithmicTraderDistributed/internal/modules/cores/misc"
	"context"
)

type Core interface {
	Run(ctx context.Context, coreConfig map[string]interface{}, instance api.InstanceServicesAPI) error
	GetCoreName() string
}

func DefaultCoreFactory(coreName string) Core {
	switch coreName {
	case "MockCore":
		return &misccores.MockCore{}
	case "MockPanicCore":
		return &misccores.MockPanicCore{}
	default:
		panic("Requested core not found")
	}
}

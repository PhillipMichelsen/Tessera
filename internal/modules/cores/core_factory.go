package cores

import (
	"AlgorithmicTraderDistributed/internal/api"
	misccores "AlgorithmicTraderDistributed/internal/modules/cores/concrete/misc"
)

func InstantiateCoreByName(coreName string, instanceAPI api.InstanceAPIInternal) Core {
	switch coreName {
	case "TestCore":
		return &misccores.TestCore{}
	}

	return nil
}

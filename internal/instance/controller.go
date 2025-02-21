package instance

import "AlgorithmicTraderDistributed/internal/api"

type Controller interface {
	Start()
	Stop()
}

func InstantiateControllerByName(controllerName string, instanceAPI api.InstanceAPIExternal) Controller {
	switch controllerName {
	default:
		return nil
	}
}

package controllers

import (
	"AlgorithmicTraderDistributed/internal/api"
)

func InstantiateControllerByName(controllerName string, instanceAPI api.InstanceAPIExternal) Controller {
	switch controllerName {
	default:
		return nil
	}
}

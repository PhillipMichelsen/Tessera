package modules

import (
	cores "AlgorithmicTraderDistributed/internal/concrete/cores"
)

type Core interface {
	Initialize(rawConfig map[string]interface{}) error
	Run()
	Stop() error
	GetType() string
}

func InstantiateCoreByName(coreName string) Core {
	switch coreName {
	case "TestCore":
		return &cores.TestCore{}
	}

	return nil
}

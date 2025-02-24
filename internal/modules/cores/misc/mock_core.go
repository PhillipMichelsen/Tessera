package cores

import (
	"AlgorithmicTraderDistributed/internal/api"
	"context"
)

type MockCore struct{}

func (c *MockCore) Run(ctx context.Context, coreConfig map[string]interface{}, instance api.InstanceServicesAPI) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *MockCore) GetCoreName() string {
	return "MockCore"
}

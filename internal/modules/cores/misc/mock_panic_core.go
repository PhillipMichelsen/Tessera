package cores

import (
	"AlgorithmicTraderDistributed/internal/api"
	"context"
)

type MockPanicCore struct{}

func (c *MockPanicCore) Run(ctx context.Context, coreConfig map[string]interface{}, instance api.InstanceServicesAPI) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			panic("MockPanicCore panic")
		}
	}
}

func (c *MockPanicCore) GetCoreName() string {
	return "MockPanicCore"
}

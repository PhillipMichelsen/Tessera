package cores

import (
	"AlgorithmicTraderDistributed/internal/api"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type TestCore struct{}

func (c *TestCore) Run(ctx context.Context, coreConfig map[string]interface{}, instance api.InstanceServicesAPI) error {
	log.Warn().Msg("!! Test Core will call panic in 5 seconds !!")

	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if time.Since(startTime) > 5*time.Second {
				panic("Panic as requested!!!")
			}

			log.Trace().Msg(fmt.Sprint("Test Core will panic in ", 5-int(time.Since(startTime).Seconds()), " seconds"))
		}
	}
}

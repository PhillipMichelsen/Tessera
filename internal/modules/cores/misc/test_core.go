package cores

import (
	"AlgorithmicTraderDistributed/internal/api"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type PanicTestCore struct{}

func (c *PanicTestCore) Run(ctx context.Context, coreConfig map[string]interface{}, instance api.InstanceServicesAPI) error {
	log.Warn().Msg("!! Test Core will call panic in 5 seconds !!")

	log.Debug().Msg(fmt.Sprint("Core Config: ", coreConfig))

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

func (c *PanicTestCore) GetCoreName() string {
	return "PanicTestCore"
}

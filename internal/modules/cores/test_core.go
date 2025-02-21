package concrete

import (
	"AlgorithmicTraderDistributed/internal/api"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

type TestCore struct {
	stopSignalChannel chan struct{}

	instanceServicesAPI api.InstanceServicesAPI
}

func (t *TestCore) Initialize(rawConfig map[string]interface{}, instanceServicesAPI api.InstanceServicesAPI) error {
	t.instanceServicesAPI = instanceServicesAPI

	return nil
}

func (t *TestCore) Run() {
	t.stopSignalChannel = make(chan struct{})

	go t.runSomeWorker()
}

func (t *TestCore) Stop() error {
	t.stopSignalChannel <- struct{}{}
	return nil
}

func (t *TestCore) GetCoreType() string {
	return "TestCore"
}

func (t *TestCore) runSomeWorker() {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if recErr, ok := r.(error); ok {
				err = recErr
			} else {
				err = fmt.Errorf("%v", r)
			}
			log.Error().Err(err).Msg("TestCore panic recovered")
		}
	}()
	log.Warn().Msg("-- Test Core will call panic randomly --")

	for {
		select {
		case <-t.stopSignalChannel:
			return
		default:
			if rand.Intn(3) == 0 {
				panic("TestCore panic!")
			} else {
				log.Trace().Msg("TestCore running...")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

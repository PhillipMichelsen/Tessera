package misc

import (
	"AlgorithmicTraderDistributed/internal/api"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type TestCore struct {
	stopSignalChannel chan struct{}
	waitGroup         sync.WaitGroup

	instanceServicesAPI api.InstanceServicesAPI
	coreErrorReceiver   func(error)
}

func (t *TestCore) Initialize(rawConfig map[string]interface{}, coreErrorReceiver func(error), instanceServicesAPI api.InstanceServicesAPI) error {
	t.instanceServicesAPI = instanceServicesAPI
	t.waitGroup = sync.WaitGroup{}
	t.coreErrorReceiver = coreErrorReceiver

	return nil
}

func (t *TestCore) Run() {
	t.stopSignalChannel = make(chan struct{})

	go t.runSomeWorker()
}

func (t *TestCore) Stop() {
	close(t.stopSignalChannel)
	t.waitGroup.Wait()
}

func (t *TestCore) GetCoreType() string {
	return "TestCore"
}

func (t *TestCore) runSomeWorker() {
	defer func() {
		t.waitGroup.Done()
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			t.coreErrorReceiver(err)
		}
	}()

	t.waitGroup.Add(1)
	log.Warn().Msg("-- Test Core will call panic after 5 seconds --")

	counter := 1
	for {
		select {
		case <-t.stopSignalChannel:
			return
		default:
			if counter == 5 {
				panic("TestCore panic!")
			} else {
				log.Trace().Msg(fmt.Sprintf("TestCore running, seconds: %d", counter))
				counter++
			}
			time.Sleep(1 * time.Second)
		}
	}
}

package concrete

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

type TestCore struct {
	stopSignalChannel chan struct{}
}

func (t *TestCore) Initialize(rawConfig map[string]interface{}) error {
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

func (t *TestCore) GetType() string {
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
	log.Warn().Msg("THE TEST CORE HAS BEEN RUN AND WILL CALL PANIC, STRICTLY FOR TESTING")

	for {
		select {
		case <-t.stopSignalChannel:
			return
		default:
			if rand.Intn(3) == 0 {
				panic("TestCore panic!!!")
			} else {
				log.Trace().Msg("TestCore running!!!")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

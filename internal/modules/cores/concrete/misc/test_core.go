package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type TestCore struct {
	stopSignalChannel chan struct{}
}

func (t *TestCore) Initialize(rawConfig map[string]interface{}) error {
	return nil
}

func (t *TestCore) Run(runtimeErrorReceiver func(error)) {
	t.stopSignalChannel = make(chan struct{})

	go t.runSomeWorker(runtimeErrorReceiver)
}

func (t *TestCore) Stop() error {
	t.stopSignalChannel <- struct{}{}
	return nil
}

func (t *TestCore) GetType() string {
	return "TestCore"
}

func (t *TestCore) runSomeWorker(runtimeErrorReceiver func(error)) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if recErr, ok := r.(error); ok {
				err = recErr
			} else {
				err = fmt.Errorf("%v", r)
			}
			runtimeErrorReceiver(err)
		}
	}()
	log.Warn().Msg("THE TEST CORE HAS BEEN RUN AND WILL CALL PANIC, STRICTLY FOR TESTING")

	for {
		select {
		case <-t.stopSignalChannel:
			return
		default:
			if rand.Intn(3) == 0 {
				//panic("TestCore panic!!!")
			} else {
				log.Trace().Msg("TestCore running!!!")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

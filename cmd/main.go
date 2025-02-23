package main

import (
	"AlgorithmicTraderDistributed/internal/instance"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.InfoLevel)

	inst := instance.NewInstance()
	//inst.AddController(instance.InstantiateController("tui", inst))

	inst.CreateModule("TestCore", uuid.Must(uuid.NewRandom()))

	inst.Start()

	modules := inst.GetModules()
	log.Info().Msgf("Modules: %v", modules)

	for _, module := range modules {
		inst.InitializeModule(module, map[string]interface{}{})
		inst.StartModule(module)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	inst.Shutdown()
}

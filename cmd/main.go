package main

import (
	"AlgorithmicTraderDistributed/_deprecated/instance"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	inst := instance.NewInstance()

	//inst.AddController(instance.InstantiateController("tui", inst))

	inst.Start()

	testCoreUUID := uuid.Must(uuid.NewRandom())
	inst.CreateModule("MockPanicCore", testCoreUUID)

	inst.InitializeModule(testCoreUUID, map[string]interface{}{})
	inst.StartModule(testCoreUUID)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals
	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
	inst.Shutdown()
}

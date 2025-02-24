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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	inst := instance.NewInstance()
	//inst.AddController(instance.InstantiateController("tui", inst))

	testCoreUUID := uuid.Must(uuid.NewRandom())
	inst.CreateModule("PanicTestCore", testCoreUUID)

	inst.Start()

	inst.InitializeModule(testCoreUUID, map[string]interface{}{"panicIn": 3})
	inst.StartModule(testCoreUUID)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals
	inst.Shutdown()
}

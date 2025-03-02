package main

import (
	"AlgorithmicTraderDistributed/internal/instance"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	workerRegistry := instance.NewWorkerRegistry()
	workerRegistry.RegisterWorker(nil) // Will cause a nil pointer dereference if left as nil!

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals
	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

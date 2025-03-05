package main

import (
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	// Create a new worker factory.
	workerFactory := worker.NewFactory()

	// Add workers to the factory.
	// workerFactory.RegisterWorker("workerType", workerConstructor)

	// Create a new node.
	inst := node.NewNode(workerFactory)

	// Start the node.
	inst.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals
	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

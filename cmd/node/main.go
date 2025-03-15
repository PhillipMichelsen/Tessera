package main

import (
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"AlgorithmicTraderDistributed/internal/worker/workers"
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
	workerFactory := worker.CombineFactories(
		workers.NewPrebuiltStandardWorkersFactory(),
		workers.NewPrebuiltBinanceSpotWorkersFactory(),
		workers.NewPrebuiltMEXCSpotWorkersFactory(),
		workers.NewStrategyWorkersFactory(),
	)

	// Create a new node.
	inst := node.NewNode(workerFactory)

	err := inst.StartDeployment("cmd/node/deployment_test.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to deploy workers from YAML")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals

	// TODO: Add a graceful shutdown mechanism. Need to implement deployment termination in Node.

	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

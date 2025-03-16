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
	// Set up logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	// Create a new worker factory.
	workerFactory := worker.CombineFactories(
		workers.NewPrebuiltStandardWorkersFactory(),
		workers.NewPrebuiltBinanceSpotWorkersFactory(),
		workers.NewPrebuiltMEXCSpotWorkersFactory(),
		workers.NewStrategyWorkersFactory(),
	)

	// Create a new node instance.
	nodeInst := node.NewNode(workerFactory)

	// Define the task files
	taskFiles := []string{"cmd/node/tasks/task1_create.yaml", "cmd/node/tasks/task2_start.yaml"}

	for _, taskFile := range taskFiles {
		log.Info().Msgf("Processing task file: %s", taskFile)

		// Read the YAML file
		yamlBytes, err := os.ReadFile(taskFile)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read file: %s", taskFile)
			continue
		}

		// Parse the task from YAML
		task, err := node.ParseTaskFromYaml(yamlBytes)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse task from file: %s", taskFile)
			continue
		}

		// Process the task
		if err := nodeInst.ProcessTask(task); err != nil {
			log.Error().Err(err).Msgf("Failed to process task from file: %s", taskFile)
			continue
		}

		log.Info().Msgf("Successfully processed task file: %s", taskFile)
	}

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Info().Msg("Shutting down gracefully...")
}

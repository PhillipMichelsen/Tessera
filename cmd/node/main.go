package main

import (
	"Tessera/internal/node"
	"Tessera/internal/worker"
	"Tessera/internal/worker/workers"
	_ "embed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

// Embed YAML files. Currently testing.
//
//go:embed tasks/task1_create.yaml
var task1Yaml []byte

//go:embed tasks/task2_start.yaml
var task2Yaml []byte

func main() {
	// Set up logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}).Level(zerolog.DebugLevel)

	// Create a new worker factory.
	workerFactory := worker.AggregateFactories(
		workers.NewPrebuiltStandardWorkersFactory(),
		workers.NewPrebuiltBinanceSpotWorkersFactory(),
		workers.NewPrebuiltMEXCSpotWorkersFactory(),
		workers.NewStrategyWorkersFactory(),
	)

	// Create a new node instance.
	nodeInst := node.NewNode(workerFactory)

	// Define embedded tasks in the desired order.
	tasksYamlBytes := [][]byte{
		task1Yaml,
		task2Yaml,
	}

	// Process each task in order.
	for _, yamlBytes := range tasksYamlBytes {
		log.Info().Msg("Processing new task")

		// Parse the task from YAML.
		task, err := node.ParseTaskFromYaml(yamlBytes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse task from embedded YAML")
			continue
		}

		// Process the task.
		if err := nodeInst.ProcessTask(task); err != nil {
			log.Error().Err(err).Msg("Failed to process task from embedded YAML")
			continue
		}

		log.Info().Msg("Successfully processed embedded task")
	}

	// Handle graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Info().Msg("Shutting down gracefully...")
}

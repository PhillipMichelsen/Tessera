package main

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"AlgorithmicTraderDistributed/internal/worker/workers"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	// Create a new worker factory.
	workerFactory := worker.NewFactory()
	workerFactory.RegisterWorkerCreationFunction("BinanceSpotKlineToOHLCV", func(moduleUUID uuid.UUID) worker.Worker {
		return workers.NewBinanceSpotKlineToOHLCVWorker(moduleUUID)
	})

	// Add workers to the factory.
	// workerFactory.RegisterWorker("workerType", workerConstructor)

	// Create a new node.
	inst := node.NewNode(workerFactory)

	w1UUID := uuid.New()
	err := inst.CreateWorker("BinanceSpotKlineToOHLCV", w1UUID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}

	jsonStr := `{"k": {"t": 1618317040000, "o":60000, "h":61000, "l":59000, "c":60500, "v":1200}}`
	go func() {
		timeTrack := time.Now()

		inst.CreateMailbox(uuid.MustParse("00000000-0000-0000-0000-000000000001"), func(message worker.Message) {
			fmt.Printf("Received message: %v, Time taken: %v \n", message, time.Since(timeTrack))
		})

		for {
			_ = inst.SendMessage(uuid.MustParse("00000000-0000-0000-0000-000000000002"), worker.Message{
				Tag:     "test_input",
				Payload: models.SerializedJSON{JSON: jsonStr},
			})
			timeTrack = time.Now()
			time.Sleep(1 * time.Second)
		}
	}()

	err = inst.StartWorker(w1UUID, map[string]interface{}{
		"input_output_mapping": map[string]interface{}{
			"test_input": map[string]interface{}{
				"destination_mailbox_uuid": "00000000-0000-0000-0000-000000000001",
				"tag":                      "test_output",
			},
		},
		"input_mailbox_uuid": "00000000-0000-0000-0000-000000000002",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	// Start the node.
	inst.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals

	_ = inst.StopWorker(w1UUID)
	inst.Stop()

	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

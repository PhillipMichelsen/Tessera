package main

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"AlgorithmicTraderDistributed/internal/worker/workers"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
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
		return &workers.BinanceSpotKlineToOHLCVWorker{}
	})
	workerFactory.RegisterWorkerCreationFunction("StandardOutput", func(moduleUUID uuid.UUID) worker.Worker {
		return &workers.StandardOutputWorker{}
	})

	// Create a new node.
	inst := node.NewNode(workerFactory)

	w1UUID := uuid.New()
	w2UUID := uuid.New()
	if err := inst.CreateWorker("BinanceSpotKlineToOHLCV", w1UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("StandardOutput", w2UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}

	w1config := workers.BinanceSpotKlineToOHLCVConfig{
		InputMailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		InputOutputMapping: map[string]struct {
			DestinationMailboxUUID uuid.UUID `yaml:"destination_mailbox_uuid"`
			Tag                    string    `yaml:"tag"`
		}{
			"test_input": {
				DestinationMailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Tag:                    "processed_ohlcv",
			},
		},
	}

	w2Config := workers.StandardOutputConfig{
		MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
	}

	w1ConfigBytes, _ := yaml.Marshal(w1config)
	w2ConfigBytes, _ := yaml.Marshal(w2Config)

	if err := inst.StartWorker(w1UUID, w1ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w2UUID, w2ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	// Start the node.
	inst.Start()

	jsonStr := `{"k": {"t": 1618317040000, "o":60000, "h":61000, "l":59000, "c":60500, "v":1200}}`
	go func() {
		for {
			_ = inst.SendMessage(uuid.MustParse("00000000-0000-0000-0000-000000000001"), worker.Message{
				Tag:     "test_input",
				Payload: models.SerializedJSON{JSON: jsonStr},
			})
			time.Sleep(1 * time.Second)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals

	_ = inst.StopWorker(w1UUID)
	_ = inst.StopWorker(w2UUID)
	inst.Stop()

	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

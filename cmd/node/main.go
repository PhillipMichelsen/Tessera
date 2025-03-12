package main

import (
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
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}).Level(zerolog.DebugLevel)

	// Create a new worker factory.
	workerFactory := worker.NewFactory()
	workerFactory.RegisterWorkerCreationFunction("BinanceSpotWebsocket", func(moduleUUID uuid.UUID) worker.Worker {
		return &workers.BinanceSpotWebsocketWorker{}
	})
	workerFactory.RegisterWorkerCreationFunction("BinanceSpotKlineToOHLCV", func(moduleUUID uuid.UUID) worker.Worker {
		return &workers.BinanceSpotKlineToOHLCVWorker{}
	})
	workerFactory.RegisterWorkerCreationFunction("BinanceSpotBookTickerToBookTicker", func(moduleUUID uuid.UUID) worker.Worker {
		return &workers.BinanceSpotBookTickerToBookTickerWorker{}
	})
	workerFactory.RegisterWorkerCreationFunction("StandardOutput", func(moduleUUID uuid.UUID) worker.Worker {
		return &workers.StandardOutputWorker{}
	})

	// Create a new node.
	inst := node.NewNode(workerFactory)

	w1UUID := uuid.New()
	w2UUID := uuid.New()
	w3UUID := uuid.New()
	if err := inst.CreateWorker("BinanceSpotWebsocket", w1UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("BinanceSpotBookTickerToBookTicker", w2UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("StandardOutput", w3UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}

	w1Config := workers.BinanceSpotWebsocketWorkerConfig{
		BaseURL: "stream.binance.com:9443",
		StreamsOutputMapping: map[string][]struct {
			MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
			Tag         string    `yaml:"tag"`
		}{
			"btcusdt@bookTicker": {
				{
					MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Tag:         "test_input",
				},
			},
		},
		BlockingSend: false,
	}

	w2Config := workers.BinanceSpotBookTickerToBookTickerConfig{
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		InputMailboxBuffer: 1000,
		InputOutputMapping: map[string]struct {
			MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
			Tag         string    `yaml:"tag"`
		}{
			"test_input": {
				MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Tag:         "processed_book_ticker",
			},
		},
		BlockingSend: false,
	}

	w3Config := workers.StandardOutputConfig{
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		InputMailboxBuffer: 1000,
	}

	w1ConfigBytes, _ := yaml.Marshal(w1Config)
	w2ConfigBytes, _ := yaml.Marshal(w2Config)
	w3ConfigBytes, _ := yaml.Marshal(w3Config)

	if err := inst.StartWorker(w3UUID, w3ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w2UUID, w2ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w1UUID, w1ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	// Start the node.
	//inst.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals

	_ = inst.StopWorker(w1UUID)
	_ = inst.StopWorker(w2UUID)
	_ = inst.StopWorker(w3UUID)
	inst.Stop()

	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

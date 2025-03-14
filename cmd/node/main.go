package main

import (
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"AlgorithmicTraderDistributed/internal/worker/workers"
	binancespot "AlgorithmicTraderDistributed/internal/worker/workers/binancespot"
	mexcspot "AlgorithmicTraderDistributed/internal/worker/workers/mexcspot"
	standard "AlgorithmicTraderDistributed/internal/worker/workers/standard"
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
	workerFactory := worker.CombineFactories(
		workers.NewPrebuiltStandardWorkersFactory(),
		workers.NewPrebuiltBinanceSpotWorkersFactory(),
	)
	workerFactory.RegisterWorkerCreationFunction("MEXCSpotWebsocket", func(uuid uuid.UUID) worker.Worker {
		return &mexcspot.MEXCSpotWebsocketWorker{}
	})

	// Create a new node.
	inst := node.NewNode(workerFactory)

	w0UUID := uuid.New()
	w1UUID := uuid.New()
	w2UUID := uuid.New()
	w3UUID := uuid.New()
	if err := inst.CreateWorker("MEXCSpotWebsocket", w0UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("BinanceSpotWebsocket", w1UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("BinanceSpotBookTickerToBookTicker", w2UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}
	if err := inst.CreateWorker("StandardOutput", w3UUID); err != nil {
		log.Fatal().Err(err).Msg("Failed to create worker")
	}

	// Strictly testing purposes. Configs would be constructed as yaml.

	w0Config := mexcspot.MEXCSpotWebsocketWorkerConfig{
		BaseURL: "wbs-api.mexc.com",
		StreamsOutputMapping: map[string]struct {
			MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
			Tag         string    `yaml:"tag"`
		}{
			"spot@public.aggre.bookTicker.v3.api.pb@100ms@BTCUSDT": {
				MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Tag:         "test_input",
			},
		},
		BlockingSend: false,
	}

	w1Config := binancespot.BinanceSpotWebsocketWorkerConfig{
		BaseURL: "stream.binance.com:9443",
		StreamsOutputMapping: map[string]struct {
			MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
			Tag         string    `yaml:"tag"`
		}{
			"btcusdt@bookTicker": {
				MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Tag:         "test_input",
			},
		},
		BlockingSend: false,
	}

	w2Config := binancespot.BinanceSpotBookTickerToBookTickerConfig{
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		InputMailboxBuffer: 1000,
		InputOutputMapping: map[string]struct {
			MailboxUUID uuid.UUID `yaml:"mailbox_uuid"`
			Tag         string    `yaml:"tag"`
		}{
			"test_input": {
				MailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Tag:         "processed_depth",
			},
		},
		BlockingSend: false,
	}

	w3Config := standard.StandardOutputConfig{
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		InputMailboxBuffer: 1000,
	}

	w0ConfigBytes, _ := yaml.Marshal(w0Config)
	_, _ = yaml.Marshal(w1Config)
	w2ConfigBytes, _ := yaml.Marshal(w2Config)
	w3ConfigBytes, _ := yaml.Marshal(w3Config)

	if err := inst.StartWorker(w3UUID, w3ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w2UUID, w2ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w0UUID, w0ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	// Start the node.
	// inst.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Main thread waiting for shutdown signal...")
	<-signals

	_ = inst.StopWorker(w0UUID)
	_ = inst.StopWorker(w1UUID)
	_ = inst.StopWorker(w2UUID)
	_ = inst.StopWorker(w3UUID)
	inst.Stop()

	log.Info().Msg("Shutdown signal received. Initiating shutdown...")
}

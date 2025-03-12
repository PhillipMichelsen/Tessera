package tests

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/node"
	"AlgorithmicTraderDistributed/internal/worker"
	"AlgorithmicTraderDistributed/internal/worker/workers"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"testing"
	"time"
)

func TestThroughput(t *testing.T) {
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
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		InputMailboxBuffer: 1000,
		InputOutputMapping: map[string]struct {
			DestinationMailboxUUID uuid.UUID `yaml:"destination_mailbox_uuid"`
			Tag                    string    `yaml:"tag"`
		}{
			"test_input": {
				DestinationMailboxUUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Tag:                    "processed_ohlcv",
			},
		},
		BlockingSend: true,
	}

	w2Config := workers.StandardOutputConfig{
		InputMailboxUUID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		InputMailboxBuffer: 1000,
	}

	w1ConfigBytes, _ := yaml.Marshal(w1config)
	w2ConfigBytes, _ := yaml.Marshal(w2Config)

	if err := inst.StartWorker(w2UUID, w2ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
	if err := inst.StartWorker(w1UUID, w1ConfigBytes); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}

	time.Sleep(100 * time.Millisecond)

	jsonStr := `{"k": {"t": 1618317040000, "o":60000, "h":61000, "l":59000, "c":60500, "v":1200}}`

	startTime := time.Now()
	for i := 0; i < 1000000; i++ {
		if err := inst.SendMessage(uuid.MustParse("00000000-0000-0000-0000-000000000001"), worker.Message{
			Tag:     "test_input",
			Payload: models.SerializedJSON{JSON: jsonStr},
		}, true); err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
	}
	duration := time.Since(startTime)

	_ = inst.StopWorker(w1UUID)
	_ = inst.StopWorker(w2UUID)

	fmt.Println(duration)
}

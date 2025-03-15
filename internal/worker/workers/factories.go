package workers

import (
	"AlgorithmicTraderDistributed/internal/worker"
	binancespot "AlgorithmicTraderDistributed/internal/worker/workers/binancespot"
	mexcspot "AlgorithmicTraderDistributed/internal/worker/workers/mexcspot"
	standard "AlgorithmicTraderDistributed/internal/worker/workers/standard"
	strategy "AlgorithmicTraderDistributed/internal/worker/workers/strategy"
	"github.com/google/uuid"
)

func NewPrebuiltStandardWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("StandardOutput", func(uuid uuid.UUID) worker.Worker {
		return &standard.StandardOutputWorker{}
	})
	factory.RegisterWorkerCreationFunction("Broadcast", func(uuid uuid.UUID) worker.Worker {
		return &standard.BroadcastWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

// TODO: Just realised we do not need the uuid parameters in these functions. They are unused.

func NewPrebuiltBinanceSpotWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("BinanceSpotWebsocket", func(uuid uuid.UUID) worker.Worker {
		return &binancespot.BinanceSpotWebsocketWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotKlineToOHLCV", func(uuid uuid.UUID) worker.Worker {
		return &binancespot.BinanceSpotKlineToOHLCVWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotBookTickerToBookTicker", func(uuid uuid.UUID) worker.Worker {
		return &binancespot.BinanceSpotBookTickerToBookTickerWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotDepthToOrderBookSnapshot", func(uuid uuid.UUID) worker.Worker {
		return &binancespot.BinanceSpotDepthToOrderBookWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotDepthUpdateToOrderBookUpdate", func(uuid uuid.UUID) worker.Worker {
		return &binancespot.BinanceSpotDepthUpdateToOrderBookWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

func NewPrebuiltMEXCSpotWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("MEXCSpotWebsocket", func(uuid uuid.UUID) worker.Worker {
		return &mexcspot.MEXCSpotWebsocketWorker{}
	})
	factory.RegisterWorkerCreationFunction("MEXCSpotBookTickerToBookTicker", func(uuid uuid.UUID) worker.Worker {
		return &mexcspot.MEXCSpotBookTickerToBookTickerWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

func NewStrategyWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("CrossMarketSpotArbitrageStrategy", func(uuid uuid.UUID) worker.Worker {
		return &strategy.CrossMarketSpotArbitrageStrategyWorker{}
	})

	return factory
}

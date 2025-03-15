package workers

import (
	"AlgorithmicTraderDistributed/internal/worker"
	binancespot "AlgorithmicTraderDistributed/internal/worker/workers/binancespot"
	mexcspot "AlgorithmicTraderDistributed/internal/worker/workers/mexcspot"
	standard "AlgorithmicTraderDistributed/internal/worker/workers/standard"
	strategy "AlgorithmicTraderDistributed/internal/worker/workers/strategy"
)

func NewPrebuiltStandardWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("StandardOutput", func() worker.Worker {
		return &standard.StandardOutputWorker{}
	})
	factory.RegisterWorkerCreationFunction("Broadcast", func() worker.Worker {
		return &standard.BroadcastWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

func NewPrebuiltBinanceSpotWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("BinanceSpotWebsocket", func() worker.Worker {
		return &binancespot.BinanceSpotWebsocketWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotKlineToOHLCV", func() worker.Worker {
		return &binancespot.BinanceSpotKlineToOHLCVWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotBookTickerToBookTicker", func() worker.Worker {
		return &binancespot.BinanceSpotBookTickerToBookTickerWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotDepthToOrderBookSnapshot", func() worker.Worker {
		return &binancespot.BinanceSpotDepthToOrderBookWorker{}
	})
	factory.RegisterWorkerCreationFunction("BinanceSpotDepthUpdateToOrderBookUpdate", func() worker.Worker {
		return &binancespot.BinanceSpotDepthUpdateToOrderBookWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

func NewPrebuiltMEXCSpotWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("MEXCSpotWebsocket", func() worker.Worker {
		return &mexcspot.MEXCSpotWebsocketWorker{}
	})
	factory.RegisterWorkerCreationFunction("MEXCSpotBookTickerToBookTicker", func() worker.Worker {
		return &mexcspot.MEXCSpotBookTickerToBookTickerWorker{}
	})
	// Add more worker types here as needed.

	return factory
}

func NewStrategyWorkersFactory() *worker.Factory {
	factory := worker.NewFactory()
	factory.RegisterWorkerCreationFunction("CrossMarketSpotArbitrageStrategy", func() worker.Worker {
		return &strategy.CrossMarketSpotArbitrageStrategyWorker{}
	})

	return factory
}

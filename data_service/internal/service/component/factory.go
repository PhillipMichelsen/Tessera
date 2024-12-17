package component

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/service/component/concrete_old"
	"fmt"
)

func CreateComponent(componentType string, routeFunction func(marketData models.MarketDataPiece), config map[string]string) (Component, error) {
	switch componentType {
	// -- Binance Spot
	case "BinanceSpotWebsocketFetcher":
		return concrete_old.NewBinanceSpotWebsocketFetcher(routeFunction, config), nil
	case "BinanceSpotKlineToOHLCVTransformer":
		return concrete_old.NewBinanceSpotKlineToOHLCVTransformer(routeFunction), nil // Non-config processor.
	case "BinanceSpotBookTickerToBookTickerTransformer":
		return concrete_old.NewBinanceSpotBookTickerToBookTickerTransformer(routeFunction), nil // Non-config processor.
	case "BinanceSpotDepthUpdateToOrderBookUpdateTransformer":
		return concrete_old.NewBinanceSpotDepthUpdateToOrderBookUpdateTransformer(routeFunction), nil // Non-config processor.
	case "BinanceSpotTradeToTradeTransformer":
		return concrete_old.NewBinanceSpotTradeToTradeTransformer(routeFunction), nil // Non-config processor.

	// -- Binance Futures
	case "BinanceFuturesWebsocketFetcher":
		return concrete_old.NewBinanceFuturesWebsocketFetcher(routeFunction, config), nil
	case "BinanceFuturesBookTickerToBookTickerTransformer":
		return concrete_old.NewBinanceFuturesBookTickerToBookTickerTransformer(routeFunction), nil // Non-config processor.
	case "BinanceFuturesKlineToOHLCVTransformer":
		return concrete_old.NewBinanceFuturesKlineToOHLCVTransformer(routeFunction), nil // Non-config processor.

	// -- General
	case "OrderBookSnapshotRangeFilter":
		createdProcessor, err := concrete_old.NewOrderBookSnapshotRangeFilter(routeFunction, config)
		if err != nil {
			return nil, err
		}
		return createdProcessor, nil

	// ---

	default:
		return nil, fmt.Errorf("component factory method: Unknown component name: %s", componentType)
	}
}

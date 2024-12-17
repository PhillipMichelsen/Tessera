package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"github.com/google/uuid"
	"log"
	"time"
)

// BinanceSpotDepthUpdateToOrderBookUpdateTransformer processes raw Binance order book update data
type BinanceSpotDepthUpdateToOrderBookUpdateTransformer struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID
}

func NewBinanceSpotDepthUpdateToOrderBookUpdateTransformer(routeFunction func(marketDataPiece models.MarketDataPiece)) *BinanceSpotDepthUpdateToOrderBookUpdateTransformer {
	return &BinanceSpotDepthUpdateToOrderBookUpdateTransformer{
		routeFunction: routeFunction,
		name:          "BinanceSpotDepthUpdateToOrderBookUpdateTransformer",
		uuid:          uuid.New(),
	}
}

// Process filters and converts Binance raw order book update data into OrderBookUpdate format
func (p *BinanceSpotDepthUpdateToOrderBookUpdateTransformer) Handle(marketDataPiece models.MarketDataPiece) {
	// Validate and extract required fields for order book update data
	validatedFields, err := helper.ValidateAndExtract(
		marketDataPiece.Payload.(map[string]interface{}),
		[]helper.FieldInfo{
			{"E", "time"},
			{"b", "[]interface{}"},
			{"a", "[]interface{}"},
			{"u", "int"},
			{"U", "int"},
		},
	)
	if err != nil {
		log.Println("Validation error at main field:", err)
		return
	}

	bidUpdates := helper.ParseBinanceOrderBookEntries(validatedFields["b"].([]interface{}))
	askUpdates := helper.ParseBinanceOrderBookEntries(validatedFields["a"].([]interface{}))

	orderBookUpdate := models.OrderBookUpdate{
		AskUpdates: askUpdates,
		BidUpdates: bidUpdates,
		Timestamp:  validatedFields["E"].(time.Time),
	}

	marketDataPiece.ExternalTimestamp = validatedFields["E"].(time.Time)
	marketDataPiece.Processors += p.name + "."
	marketDataPiece.Payload = orderBookUpdate

	p.routeFunction(marketDataPiece)
}

func (p *BinanceSpotDepthUpdateToOrderBookUpdateTransformer) Start() error {
	return nil
}

func (p *BinanceSpotDepthUpdateToOrderBookUpdateTransformer) Stop() error {
	return nil
}

func (p *BinanceSpotDepthUpdateToOrderBookUpdateTransformer) GetName() string {
	return p.name
}

func (p *BinanceSpotDepthUpdateToOrderBookUpdateTransformer) GetUUID() uuid.UUID {
	return p.uuid
}

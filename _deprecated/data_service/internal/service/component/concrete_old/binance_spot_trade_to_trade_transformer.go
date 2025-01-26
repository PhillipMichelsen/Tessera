package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"github.com/google/uuid"
	"log"
	"time"
)

// BinanceSpotTradeToTradeTransformer is a processor that converts raw Binance trade data into the common models.Trade format
type BinanceSpotTradeToTradeTransformer struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID
}

func NewBinanceSpotTradeToTradeTransformer(routeFunction func(marketDataPiece models.MarketDataPiece)) *BinanceSpotTradeToTradeTransformer {
	return &BinanceSpotTradeToTradeTransformer{
		routeFunction: routeFunction,
		name:          "BinanceSpotTradeToTradeTransformer",
		uuid:          uuid.New(),
	}
}

// Handle filters and converts Binance raw trade data into Trade format
func (p *BinanceSpotTradeToTradeTransformer) Handle(marketDataPiece models.MarketDataPiece) {
	// Validate and extract required fields for trade data
	validatedFields, err := helper.ValidateAndExtract(
		marketDataPiece.Payload.(map[string]interface{}),
		[]helper.FieldInfo{
			{"p", "float64"},
			{"q", "float64"},
			{"T", "time"},
			{"E", "time"},
		},
	)
	if err != nil {
		log.Println("Validation error in main fields:", err)
		return
	}

	// Map validated data to Trade model
	tradeData := models.Trade{
		Price:     validatedFields["p"].(float64),
		Quantity:  validatedFields["q"].(float64),
		Timestamp: validatedFields["T"].(time.Time),
	}

	// Update the MarketDataPiece struct
	marketDataPiece.ExternalTimestamp = validatedFields["E"].(time.Time)
	marketDataPiece.Processors += p.name + "."
	marketDataPiece.Payload = tradeData

	// Route the transformed Trade data
	p.routeFunction(marketDataPiece)
}

func (p *BinanceSpotTradeToTradeTransformer) Start() error {
	return nil
}

func (p *BinanceSpotTradeToTradeTransformer) Stop() error {
	return nil
}

func (p *BinanceSpotTradeToTradeTransformer) GetName() string {
	return p.name
}

func (p *BinanceSpotTradeToTradeTransformer) GetUUID() uuid.UUID {
	return p.uuid
}

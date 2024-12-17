package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"github.com/google/uuid"
	"log"
	"time"
)

// BinanceFuturesBookTickerToBookTickerTransformer is a processor that converts raw Binance futures book ticker data into the common models.BookTicker format
type BinanceFuturesBookTickerToBookTickerTransformer struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID
}

func NewBinanceFuturesBookTickerToBookTickerTransformer(routeFunction func(marketDataPiece models.MarketDataPiece)) *BinanceFuturesBookTickerToBookTickerTransformer {
	return &BinanceFuturesBookTickerToBookTickerTransformer{
		routeFunction: routeFunction,
		name:          "BinanceFuturesBookTickerToBookTickerTransformer",
		uuid:          uuid.New(),
	}
}

// Handle filters and converts Binance raw futures book ticker data into BookTicker format
func (p *BinanceFuturesBookTickerToBookTickerTransformer) Handle(marketDataPiece models.MarketDataPiece) {
	// Validate and extract main fields
	validatedFields, err := helper.ValidateAndExtract(
		marketDataPiece.Payload.(map[string]interface{}),
		[]helper.FieldInfo{
			{"E", "time"},
			{"s", "string"},  // Symbol
			{"b", "float64"}, // Best bid price
			{"B", "float64"}, // Best bid quantity
			{"a", "float64"}, // Best ask price
			{"A", "float64"}, // Best ask quantity
		})
	if err != nil {
		log.Println("Validation error in futures book ticker fields:", err)
		return
	}

	// Create BookTicker object
	bookTicker := models.BookTicker{
		BidPrice:    validatedFields["b"].(float64),
		BidQuantity: validatedFields["B"].(float64),
		AskPrice:    validatedFields["a"].(float64),
		AskQuantity: validatedFields["A"].(float64),
		Timestamp:   validatedFields["E"].(time.Time),
	}

	marketDataPiece.ExternalTimestamp = validatedFields["E"].(time.Time)
	marketDataPiece.Processors += p.name + "."
	marketDataPiece.Payload = bookTicker

	p.routeFunction(marketDataPiece)
}

func (p *BinanceFuturesBookTickerToBookTickerTransformer) Start() error {
	return nil
}

func (p *BinanceFuturesBookTickerToBookTickerTransformer) Stop() error {
	return nil
}

func (p *BinanceFuturesBookTickerToBookTickerTransformer) GetName() string {
	return p.name
}

func (p *BinanceFuturesBookTickerToBookTickerTransformer) GetUUID() uuid.UUID {
	return p.uuid
}

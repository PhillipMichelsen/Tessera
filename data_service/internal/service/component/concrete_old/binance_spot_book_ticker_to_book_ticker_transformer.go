package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"github.com/google/uuid"
	"log"
	"time"
)

// BinanceSpotBookTickerToBookTickerTransformer is a processor that converts raw Binance book ticker data into the common models.BookTicker format
type BinanceSpotBookTickerToBookTickerTransformer struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID
}

func NewBinanceSpotBookTickerToBookTickerTransformer(routeFunction func(marketDataPiece models.MarketDataPiece)) *BinanceSpotBookTickerToBookTickerTransformer {
	return &BinanceSpotBookTickerToBookTickerTransformer{
		routeFunction: routeFunction,
		name:          "BinanceSpotBookTickerToBookTickerTransformer",
		uuid:          uuid.New(),
	}
}

// Handle filters and converts Binance raw book ticker data into BookTicker format
func (p *BinanceSpotBookTickerToBookTickerTransformer) Handle(marketDataPiece models.MarketDataPiece) {
	// Validate and extract main fields
	validatedFields, err := helper.ValidateAndExtract(
		marketDataPiece.Payload.(map[string]interface{}),
		[]helper.FieldInfo{
			{"b", "float64"},
			{"B", "float64"},
			{"a", "float64"},
			{"A", "float64"},
		})
	if err != nil {
		log.Println("Validation error in book ticker fields:", err)
		return
	}

	// Create BookTicker object
	bookTicker := models.BookTicker{
		BidPrice:    validatedFields["b"].(float64),
		BidQuantity: validatedFields["B"].(float64),
		AskPrice:    validatedFields["a"].(float64),
		AskQuantity: validatedFields["A"].(float64),
		Timestamp:   time.Now().UTC(), // Use current UTC time since Binance data doesn't provide a timestamp
	}

	marketDataPiece.ExternalTimestamp = time.Now().UTC()
	marketDataPiece.Processors += p.name + "."
	marketDataPiece.Payload = bookTicker

	p.routeFunction(marketDataPiece)
}

func (p *BinanceSpotBookTickerToBookTickerTransformer) Start() error {
	return nil
}

func (p *BinanceSpotBookTickerToBookTickerTransformer) Stop() error {
	return nil
}

func (p *BinanceSpotBookTickerToBookTickerTransformer) GetName() string {
	return p.name
}

func (p *BinanceSpotBookTickerToBookTickerTransformer) GetUUID() uuid.UUID {
	return p.uuid
}

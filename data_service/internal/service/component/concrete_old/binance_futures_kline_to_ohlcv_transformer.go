package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"github.com/google/uuid"
	"log"
	"time"
)

type KlinePayload struct {
	EventTime time.Time   `json:"E"`
	Kline     KlineFields `json:"k"`
}

type KlineFields struct {
	Timestamp time.Time `json:"t"`
	Open      float64   `json:"o"`
	High      float64   `json:"h"`
	Low       float64   `json:"l"`
	Close     float64   `json:"c"`
	Volume    float64   `json:"v"`
	IsClosed  bool      `json:"x"`
}

// BinanceFuturesKlineToOHLCVTransformer is a processor that converts raw Binance Futures kline data into the common models.OHLCV format
type BinanceFuturesKlineToOHLCVTransformer struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID
}

func NewBinanceFuturesKlineToOHLCVTransformer(routeFunction func(marketDataPiece models.MarketDataPiece)) *BinanceFuturesKlineToOHLCVTransformer {
	return &BinanceFuturesKlineToOHLCVTransformer{
		routeFunction: routeFunction,
		name:          "BinanceFuturesKlineToOHLCVTransformer",
		uuid:          uuid.New(),
	}
}

// Handle filters and converts Binance raw kline data into OHLCVData format
func (p *BinanceFuturesKlineToOHLCVTransformer) Handle(marketDataPiece models.MarketDataPiece) {
	// Validate and extract main event fields
	validatedFields, err := helper.ValidateAndExtract(
		marketDataPiece.Payload.(map[string]interface{}),
		[]helper.FieldInfo{
			{"E", "time"},
			{"k", "map[string]interface{}"},
		})
	if err != nil {
		log.Println("Validation error in main fields:", err)
		return
	}

	// Validate and extract kline fields
	validatedKlineFields, err := helper.ValidateAndExtract(validatedFields["k"].(map[string]interface{}), []helper.FieldInfo{
		{"t", "time"},
		{"o", "float64"},
		{"h", "float64"},
		{"l", "float64"},
		{"c", "float64"},
		{"v", "float64"},
		{"x", "bool"},
	})
	if err != nil {
		log.Println("Validation error in kline fields:", err)
		return
	}

	// Only process if the kline is closed
	if !validatedKlineFields["x"].(bool) {
		return
	}

	ohlcvData := models.OHLCV{
		Timestamp: validatedKlineFields["t"].(time.Time),
		Open:      validatedKlineFields["o"].(float64),
		High:      validatedKlineFields["h"].(float64),
		Low:       validatedKlineFields["l"].(float64),
		Close:     validatedKlineFields["c"].(float64),
		Volume:    validatedKlineFields["v"].(float64),
	}

	marketDataPiece.ExternalTimestamp = validatedFields["E"].(time.Time)
	marketDataPiece.Processors += p.name + "."
	marketDataPiece.Payload = ohlcvData

	p.routeFunction(marketDataPiece)
}

func (p *BinanceFuturesKlineToOHLCVTransformer) Start() error {
	return nil
}

func (p *BinanceFuturesKlineToOHLCVTransformer) Stop() error {
	return nil
}

func (p *BinanceFuturesKlineToOHLCVTransformer) GetName() string {
	return p.name
}

func (p *BinanceFuturesKlineToOHLCVTransformer) GetUUID() uuid.UUID {
	return p.uuid
}

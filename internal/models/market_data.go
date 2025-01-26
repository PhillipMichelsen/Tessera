package models

import (
	"time"
)

// MarketData General data type
type MarketData interface{}

type OHLCV struct {
	Open      float64   `json:"O"`
	High      float64   `json:"H"`
	Low       float64   `json:"L"`
	Close     float64   `json:"C"`
	Volume    float64   `json:"V"`
	Timestamp time.Time `json:"T"`
}

type Trade struct {
	Price     float64   `json:"P"`
	Quantity  float64   `json:"Q"`
	Timestamp time.Time `json:"T"`
}

type BookTicker struct {
	BidPrice    float64   `json:"b"`
	BidQuantity float64   `json:"B"`
	AskPrice    float64   `json:"a"`
	AskQuantity float64   `json:"A"`
	Timestamp   time.Time `json:"T"`
}

type OrderBookEntry struct {
	Price    float64 `json:"P"`
	Quantity float64 `json:"Q"`
}

type OrderBookUpdate struct {
	AskUpdates []OrderBookEntry `json:"A"`
	BidUpdates []OrderBookEntry `json:"B"`
	Timestamp  time.Time        `json:"T"`
}

type OrderBookSnapshot struct {
	Asks      []OrderBookEntry `json:"A"`
	Bids      []OrderBookEntry `json:"B"`
	Timestamp time.Time        `json:"T"`
}

type SerializedJSON struct {
	JSON string `json:"data"`
}

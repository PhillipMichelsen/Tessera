package models

import (
	"time"
)

// OHLCV Regular OHLCV data, values are in float64 except for timestamp which is a time.Time
type OHLCV struct {
	Open      float64   `json:"O"`
	High      float64   `json:"H"`
	Low       float64   `json:"L"`
	Close     float64   `json:"C"`
	Volume    float64   `json:"V"`
	Timestamp time.Time `json:"T"`
}

// Trade A single (or aggregated, depends on context) trade, values are in float64 except for timestamp which is a time.Time
type Trade struct {
	Price     float64   `json:"P"`
	Quantity  float64   `json:"Q"`
	Timestamp time.Time `json:"T"`
}

// BookTicker Top of the book ticker, values are in float64 except for timestamp which is a time.Time
type BookTicker struct {
	BidPrice    float64   `json:"B"`
	BidQuantity float64   `json:"b"`
	AskPrice    float64   `json:"A"`
	AskQuantity float64   `json:"a"`
	Timestamp   time.Time `json:"T"`
}

// OrderBookEntry A single price level in the order book, values are in float64
type OrderBookEntry struct {
	Price    float64 `json:"P"`
	Quantity float64 `json:"Q"`
}

// OrderBookUpdate New order book update, values are in float64 except for timestamp which is a time.Time
type OrderBookUpdate struct {
	AskUpdates []OrderBookEntry `json:"A"`
	BidUpdates []OrderBookEntry `json:"B"`
	Timestamp  time.Time        `json:"T"`
}

// OrderBookSnapshot Order book snapshot, values are in float64 except for timestamp which is a time.Time
type OrderBookSnapshot struct {
	Asks      []OrderBookEntry `json:"A"`
	Bids      []OrderBookEntry `json:"B"`
	Timestamp time.Time        `json:"T"`
}

// SerializedJSON JSON data as a string
type SerializedJSON struct {
	JSON string `json:"D"`
}

package concrete

import (
	"github.com/google/uuid"
	"github.com/pebbe/zmq4"
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

type BinanceFuturesKlineToOHLCVTransformer struct {
	uuid           uuid.UUID
	name           string
	subscriberURLs []string
	publisherURL   string
	subscriber     *zmq4.Socket
	publisher      *zmq4.Socket
	stopSignal     chan struct{}
}

package strategy

import (
	"AlgorithimcTraderDistributed/common/models"
	"strings"
	"time"
)

type Status int

const (
	Uninitialized Status = iota
	Initialized
	Staging
	Staged
	Live
)

type Asset struct {
	Source string
	Symbol string
}

func (a Asset) String() string {
	return a.Source + ":" + a.Symbol
}

func AssetFromString(assetStr string) Asset {
	parts := strings.Split(assetStr, ":")
	if len(parts) != 2 {
		return Asset{}
	}
	return Asset{
		Source: parts[0],
		Symbol: parts[1],
	}
}

func AssetFromMarketDataPiece(data models.MarketDataPiece) Asset {
	return Asset{
		Source: data.Source,
		Symbol: data.Symbol,
	}
}

type Position struct {
	PositionID         string
	StrategyID         string
	PositionType       string
	PositionStatus     string
	PositionSize       float64
	PositionOpenPrice  float64
	PositionClosePrice float64
	PositionOpenTime   time.Time
	PositionCloseTime  time.Time
}

type Order struct {
	OrderID     string
	StrategyID  string
	OrderType   string
	OrderStatus string
	OrderSize   float64
	OrderPrice  float64
	OrderTime   time.Time
}

type Strategy interface {
	Initialize(config map[string]string) error
	Start() error
	SetLive(live bool) error
	SetAutoLive(autoLive bool)
	Stop() error
}

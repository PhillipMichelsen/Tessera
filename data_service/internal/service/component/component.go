package component

import (
	"AlgorithimcTraderDistributed/common/models"
	"github.com/google/uuid"
)

type Component interface {
	Handle(data models.MarketDataPiece)
	Start() error
	Stop() error
	GetName() string
	GetUUID() uuid.UUID
}

package models

import (
	"github.com/google/uuid"
	"time"
)

type MarketDataPacket struct {
	CurrentUUID            uuid.UUID   `json:"current_uuid"`
	TransmissionPath       []uuid.UUID `json:"transmission_path"`
	TransmissionTimestamps []time.Time `json:"transmission_timestamps"`
	Payload                MarketData  `json:"payload"`
}

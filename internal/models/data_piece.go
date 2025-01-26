package models

import (
	"github.com/google/uuid"
	"time"
)

type MarketDataPacket struct {
	CurrentUUID           uuid.UUID   `json:"current_uuid"`
	UUIDHistory           []uuid.UUID `json:"uuid_history"`
	UUIDHistoryTimestamps []time.Time `json:"uuid_history_timestamps"`
	Payload               MarketData  `json:"payload"`
}

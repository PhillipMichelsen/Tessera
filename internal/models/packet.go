package models

import (
	"github.com/google/uuid"
)

type Packet struct {
	SourceUUID      uuid.UUID   `json:"S"`
	DestinationUUID uuid.UUID   `json:"D"`
	Payload         interface{} `json:"P"`
}

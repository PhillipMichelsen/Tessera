package models

import (
	"github.com/google/uuid"
)

type Packet struct {
	SourceModuleUUID      uuid.UUID   `json:"S"`
	DestinationModuleUUID uuid.UUID   `json:"D"`
	Payload               interface{} `json:"P"`
}

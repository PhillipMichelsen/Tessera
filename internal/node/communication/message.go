package communication

import (
	"github.com/google/uuid"
	"time"
)

type IntraNodeMessage struct {
	SourceWorkerUUID      uuid.UUID
	DestinationWorkerUUID uuid.UUID
	SentTimestamp         time.Time
	Payload               interface{}
}

type InterNodeMessage struct {
	SourceNodeUUID        uuid.UUID
	SourceWorkerUUID      uuid.UUID
	DestinationNodeUUID   uuid.UUID
	DestinationWorkerUUID uuid.UUID
	SentTimestamp         time.Time
	Payload               interface{}
}

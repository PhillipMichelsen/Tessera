package worker

import (
	"github.com/google/uuid"
	"time"
)

type InboxMessage struct {
	SenderUUID uuid.UUID
	SendTime   time.Time
	Payload    interface{}
}

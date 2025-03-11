package worker

import "github.com/google/uuid"

type OutputMetadata struct {
	DestinationMailboxUUID uuid.UUID
	Tag                    string
}

type InputOutputMapping map[string]OutputMetadata

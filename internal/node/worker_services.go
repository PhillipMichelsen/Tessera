package node

import (
	"fmt"
	"github.com/PhillipMichelsen/Tessera/internal/worker"
	"github.com/google/uuid"
)

type WorkerServices struct {
	node *Node

	mailboxUUIDs []uuid.UUID
	messagesSent int
}

func NewWorkerServices(node *Node) WorkerServices {
	return WorkerServices{
		node:         node,
		mailboxUUIDs: make([]uuid.UUID, 0),
		messagesSent: 0,
	}
}

func (ws WorkerServices) SendMessage(destinationMailboxUUID uuid.UUID, message worker.Message, block bool) error {
	// Intra-node message case, can be directly pushed to mailbox.
	if ws.node.dispatcher.CheckMailboxExists(destinationMailboxUUID) {
		var err error
		if block {
			err = ws.node.dispatcher.PushMessageBlocking(destinationMailboxUUID, message)
		} else {
			err = ws.node.dispatcher.PushMessage(destinationMailboxUUID, message)
		}

		if err != nil {
			return fmt.Errorf("failed to send message to destination mailbox %s: %w", destinationMailboxUUID, err)
		}
		ws.messagesSent++
		return nil
	}

	// Inter-node message case, needs to be routed through the bridge.
	// First, resolve the destination node. Done with the cluster discovery service.
	// Then, bridge the message to the destination node with the bridge.
	// TODO: Implement the above. Start in a new branch, do alongside the orchestrator and related components.

	return fmt.Errorf("unimplemented non intra-node message routing to destination mailbox %s", destinationMailboxUUID)
}

func (ws WorkerServices) CreateMailbox(mailboxUUID uuid.UUID, bufferSize int) (<-chan any, error) {
	mailbox, err := ws.node.dispatcher.CreateMailbox(mailboxUUID, bufferSize)
	if err != nil {
		return nil, err
	}
	ws.mailboxUUIDs = append(ws.mailboxUUIDs, mailboxUUID)
	return mailbox, nil
}

func (ws WorkerServices) RemoveMailbox(mailboxUUID uuid.UUID) {
	ws.node.dispatcher.RemoveMailbox(mailboxUUID)

	for i, currentUUID := range ws.mailboxUUIDs {
		if currentUUID == mailboxUUID {
			ws.mailboxUUIDs = append(ws.mailboxUUIDs[:i], ws.mailboxUUIDs[i+1:]...)
			break
		}
	}
}

func (ws WorkerServices) cleanupMailboxes() {
	for _, mailboxUUID := range ws.mailboxUUIDs {
		ws.RemoveMailbox(mailboxUUID)
	}
}

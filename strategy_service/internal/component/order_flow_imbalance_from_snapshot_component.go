package component

import (
	"AlgorithimcTraderDistributed/common/models"
	"sync"
)

type OrderFlowImbalanceFromSnapshotComponent struct {
	mu                        sync.Mutex
	previousOrderBookSnapshot models.OrderBookSnapshot
	currentOrderBookSnapshot  models.OrderBookSnapshot
}

// NewOrderFlowImbalanceFromSnapshotComponent initializes the component with an initial snapshot
func NewOrderFlowImbalanceFromSnapshotComponent() *OrderFlowImbalanceFromSnapshotComponent {
	return &OrderFlowImbalanceFromSnapshotComponent{
		previousOrderBookSnapshot: models.OrderBookSnapshot{},
		currentOrderBookSnapshot:  models.OrderBookSnapshot{},
	}
}

// NewOrderBookSnapshot updates the snapshots, shifting current to previous
func (c *OrderFlowImbalanceFromSnapshotComponent) NewOrderBookSnapshot(newSnapshot models.OrderBookSnapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.previousOrderBookSnapshot = c.currentOrderBookSnapshot
	c.currentOrderBookSnapshot = newSnapshot
}

// GetOrderFlowImbalance calculates the imbalance based on changes between snapshots
func (c *OrderFlowImbalanceFromSnapshotComponent) GetOrderFlowImbalance() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	askChange := calculateOrderBookChange(c.previousOrderBookSnapshot.Asks, c.currentOrderBookSnapshot.Asks)
	bidChange := calculateOrderBookChange(c.previousOrderBookSnapshot.Bids, c.currentOrderBookSnapshot.Bids)

	// Avoid division by zero
	totalChange := askChange + bidChange
	if totalChange == 0 {
		return 0.0
	}

	return (askChange - bidChange) / totalChange
}

// calculateOrderBookChange calculates the total change in weighted price*quantity for an order book side
func calculateOrderBookChange(previous []models.OrderBookEntry, current []models.OrderBookEntry) float64 {
	priceMap := make(map[float64]float64)

	// Populate map with previous snapshot
	for _, entry := range previous {
		priceMap[entry.Price] -= entry.Quantity
	}

	// Update map with current snapshot
	for _, entry := range current {
		priceMap[entry.Price] += entry.Quantity
	}

	// Calculate the total change
	totalChange := 0.0
	for price, quantity := range priceMap {
		totalChange += price * quantity
	}

	return totalChange
}

func (c *OrderFlowImbalanceFromSnapshotComponent) IsReady() bool {
	prevExists := len(c.previousOrderBookSnapshot.Asks) > 0 && len(c.previousOrderBookSnapshot.Bids) > 0
	currExists := len(c.currentOrderBookSnapshot.Asks) > 0 && len(c.currentOrderBookSnapshot.Bids) > 0
	return prevExists && currExists
}

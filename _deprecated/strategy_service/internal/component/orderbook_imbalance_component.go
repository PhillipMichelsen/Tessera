package component

import "AlgorithimcTraderDistributed/common/models"

type OrderBookSnapshotImbalance struct {
}

func NewOrderBookSnapshotImbalance() *OrderBookSnapshotImbalance {
	return &OrderBookSnapshotImbalance{}
}

func (o *OrderBookSnapshotImbalance) CalculateOrderBookSnapshotImbalance(orderBookSnapshot models.OrderBookSnapshot) float64 {
	asks := orderBookSnapshot.Asks
	bids := orderBookSnapshot.Bids
	askSum := 0.0
	bidSum := 0.0
	for _, ask := range asks {
		askSum += ask.Price * ask.Quantity
	}
	for _, bid := range bids {
		bidSum += bid.Price * bid.Quantity
	}
	return (askSum - bidSum) / (askSum + bidSum)
}

package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
)

// OrderBookSnapshotRangeFilter filters order book data based on a specified range percentage.
type OrderBookSnapshotRangeFilter struct {
	routeFunction func(marketDataPiece models.MarketDataPiece)
	name          string
	uuid          uuid.UUID

	rangePercentage float64
}

// NewOrderBookSnapshotRangeFilter initializes a new OrderBookSnapshotRangeFilter.
func NewOrderBookSnapshotRangeFilter(routeFunction func(marketDataPiece models.MarketDataPiece), config map[string]string) (*OrderBookSnapshotRangeFilter, error) {
	// Validate and extract the range percentage from the config. TODO: Do this nicer, maybe with a helper function.
	rangePercentage, ok := config["range_percentage"]
	if !ok {
		return nil, fmt.Errorf("range_percentage not found in config")
	}

	rangePercentageStr, err := strconv.ParseFloat(rangePercentage, 64)
	if err != nil {
		return nil, fmt.Errorf("range_percentage must be a float64")
	}

	return &OrderBookSnapshotRangeFilter{
		routeFunction:   routeFunction,
		rangePercentage: rangePercentageStr,
		name:            "OrderBookSnapshotRangeFilter",
		uuid:            uuid.New(),
	}, nil
}

// Handle filters the order book snapshot data within a range defined by the percentage threshold.
func (p *OrderBookSnapshotRangeFilter) Handle(marketDataPiece models.MarketDataPiece) {
	// Ensure the data payload is of type OrderBookSnapshot.
	orderBookSnapshot, ok := marketDataPiece.Payload.(models.OrderBookSnapshot)
	if !ok {
		log.Println("Invalid payload type for OrderBookSnapshotRangeFilter")
		return
	}

	// Ensure there are bids and asks to process.
	if len(orderBookSnapshot.Bids) == 0 || len(orderBookSnapshot.Asks) == 0 {
		log.Println("Empty bids or asks in OrderBookSnapshot")
		return
	}

	// Calculate the price limits based on the range percentage.
	bestBid := orderBookSnapshot.Bids[0].Price
	bestAsk := orderBookSnapshot.Asks[0].Price
	bidLimit := bestBid * (1 - p.rangePercentage/100)
	askLimit := bestAsk * (1 + p.rangePercentage/100)

	// Filter bids and asks within the range limits.
	filteredBids := filterOrderBookEntries(orderBookSnapshot.Bids, bidLimit, true)
	filteredAsks := filterOrderBookEntries(orderBookSnapshot.Asks, askLimit, false)

	// Update the payload with the filtered snapshot.
	filteredSnapshot := models.OrderBookSnapshot{
		Bids:      filteredBids,
		Asks:      filteredAsks,
		Timestamp: orderBookSnapshot.Timestamp,
	}

	marketDataPiece.Processors += fmt.Sprintf("%s_%.2f%%.", p.name, p.rangePercentage)
	marketDataPiece.Payload = filteredSnapshot

	p.routeFunction(marketDataPiece)
}

// filterOrderBookEntries filters order book entries based on a limit and direction.
func filterOrderBookEntries(entries []models.OrderBookEntry, limit float64, isBid bool) []models.OrderBookEntry {
	filtered := make([]models.OrderBookEntry, 0, len(entries))
	for _, entry := range entries {
		if (isBid && entry.Price >= limit) || (!isBid && entry.Price <= limit) {
			filtered = append(filtered, entry)
		} else if (isBid && entry.Price < limit) || (!isBid && entry.Price > limit) {
			break
		}
	}
	return filtered
}

func (p *OrderBookSnapshotRangeFilter) Start() error {
	return nil
}

func (p *OrderBookSnapshotRangeFilter) Stop() error {
	return nil
}

func (p *OrderBookSnapshotRangeFilter) GetName() string {
	return fmt.Sprintf("%s_%.2f%%.", p.name, p.rangePercentage)
}

func (p *OrderBookSnapshotRangeFilter) GetUUID() uuid.UUID {
	return p.uuid
}

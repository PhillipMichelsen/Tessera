package helper

import (
	"AlgorithimcTraderDistributed/common/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ParseBinanceOrderBookEntries converts raw bid/ask data into a slice of OrderBookEntry
func ParseBinanceOrderBookEntries(entriesData []interface{}) []models.OrderBookEntry {
	var entries []models.OrderBookEntry
	for _, entry := range entriesData {
		entrySlice, ok := entry.([]interface{})
		if !ok || len(entrySlice) != 2 {
			continue
		}
		priceStr, ok1 := entrySlice[0].(string)
		quantityStr, ok2 := entrySlice[1].(string)
		if !ok1 || !ok2 {
			continue
		}
		price := ParseFloat(priceStr)
		quantity := ParseFloat(quantityStr)
		entries = append(entries, models.OrderBookEntry{Price: price, Quantity: quantity})
	}
	return entries
}

// FetchBinanceOrderBookSnapshot retrieves an order book snapshot for a given symbol
func FetchBinanceOrderBookSnapshot(symbol string) (models.OrderBookSnapshot, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/depth?symbol=%s&limit=250", strings.ToUpper(symbol))
	resp, err := http.Get(url)
	if err != nil {
		return models.OrderBookSnapshot{}, fmt.Errorf("failed to fetch order book snapshot: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return models.OrderBookSnapshot{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.OrderBookSnapshot{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var snapshotResponse struct {
		Bids []interface{} `json:"bids"`
		Asks []interface{} `json:"asks"`
	}

	if err := json.Unmarshal(body, &snapshotResponse); err != nil {
		return models.OrderBookSnapshot{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Parse bids and asks into OrderBookSnapshot format
	bids := ParseBinanceOrderBookEntries(snapshotResponse.Bids)
	asks := ParseBinanceOrderBookEntries(snapshotResponse.Asks)

	return models.OrderBookSnapshot{
		Bids:      bids,
		Asks:      asks,
		Timestamp: time.Now(),
	}, nil
}

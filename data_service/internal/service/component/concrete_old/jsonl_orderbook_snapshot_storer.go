package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type JSONLOrderBookSnapshotStorer struct {
	dirPath string
	files   map[string]*os.File // Map of open files keyed by file identifier
	mu      sync.Mutex          // Mutex to ensure thread-safe access to files
}

// NewJSONLOrderBookSnapshotStorer initializes a new JSONLOrderBookSnapshotStorer with the given directory path
func NewJSONLOrderBookSnapshotStorer(dirPath string) *JSONLOrderBookSnapshotStorer {
	// Ensure the directory exists
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dirPath, err)
	}

	return &JSONLOrderBookSnapshotStorer{
		dirPath: dirPath,
		files:   make(map[string]*os.File),
	}
}

// Store writes the OrderBookSnapshot data to a JSONL file based on metadata from MarketDataPiece
func (s *JSONLOrderBookSnapshotStorer) Store(marketDataPiece models.MarketDataPiece) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the payload is of type OrderBookSnapshot
	orderBookSnapshot, ok := marketDataPiece.Payload.(models.OrderBookSnapshot)
	if !ok {
		log.Println("Invalid data format for JSONL storage, expected OrderBookSnapshot")
		return
	}

	// Serialize the snapshot to JSON
	jsonData, err := json.Marshal(orderBookSnapshot)
	if err != nil {
		log.Printf("Error serializing order book snapshot to JSON: %v\n", err)
		return
	}

	// Generate a unique file identifier based on metadata fields
	fileIdentifier := fmt.Sprintf("%s_%s_%s_%s_%s",
		marketDataPiece.Source,
		marketDataPiece.Symbol,
		marketDataPiece.BaseType,
		marketDataPiece.Interval,
		marketDataPiece.Processors,
	)

	// Get or create a file for the identifier
	file, err := s.getOrCreateFile(fileIdentifier)
	if err != nil {
		log.Printf("Error creating file for identifier %s: %v\n", fileIdentifier, err)
		return
	}

	// Write the JSON data as a new line in the JSONL file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Printf("Error writing JSON data to file for %s: %v\n", fileIdentifier, err)
	}
	_, err = file.Write([]byte("\n"))
	if err != nil {
		log.Printf("Error writing newline to file for %s: %v\n", fileIdentifier, err)
	} // Separate entries by newlines
}

// getOrCreateFile checks if a file for the given identifier exists; if not, creates one
func (s *JSONLOrderBookSnapshotStorer) getOrCreateFile(fileIdentifier string) (*os.File, error) {
	if file, exists := s.files[fileIdentifier]; exists {
		return file, nil
	}

	// Generate a unique file path with a timestamp in the filename
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.jsonl", fileIdentifier, timestamp) // Use .jsonl extension
	filePath := filepath.Join(s.dirPath, fileName)

	// Open a new file for the given identifier
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open or create file: %w", err)
	}

	// Handle the file for the duration of the session
	s.files[fileIdentifier] = file

	return file, nil
}

func (s *JSONLOrderBookSnapshotStorer) Start() error {
	return nil
}

// Stop Close releases all open files, called at the end of the session
func (s *JSONLOrderBookSnapshotStorer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for identifier, file := range s.files {
		log.Printf("Closing file for identifier %s", identifier)
		err := file.Close()
		if err != nil {
			return err
		}
		delete(s.files, identifier)
	}
	return nil
}

package concrete_old

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/helper"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CSVTradeStorer struct {
	dirPath string
	files   map[string]*os.File    // Map of open files keyed by file identifier
	writers map[string]*csv.Writer // Map of CSV writers keyed by file identifier
	mu      sync.Mutex             // Mutex to ensure thread-safe access to files and writers
}

// NewCSVTradeStorer initializes a new CSVTradeStorer with the given directory path
func NewCSVTradeStorer(dirPath string) *CSVTradeStorer {
	// Ensure the directory exists
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory %s: %v", dirPath, err)
	}

	return &CSVTradeStorer{
		dirPath: dirPath,
		files:   make(map[string]*os.File),
		writers: make(map[string]*csv.Writer),
	}
}

// Store writes the Trade data to a CSV file based on metadata from MarketDataPiece
func (s *CSVTradeStorer) Store(marketDataPiece models.MarketDataPiece) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the payload is of type Trade
	tradeData, ok := marketDataPiece.Payload.(models.Trade)
	if !ok {
		log.Println("Invalid data format for CSV storage, expected Trade")
		return
	}

	// Generate a unique identifier for the file based on metadata fields
	fileIdentifier := fmt.Sprintf("%s_%s_%s_%s",
		marketDataPiece.Source,
		marketDataPiece.Symbol,
		marketDataPiece.BaseType,
		marketDataPiece.Interval,
	)

	// Get or create CSV writer for the identifier
	writer, err := s.getOrCreateWriter(fileIdentifier)
	if err != nil {
		log.Printf("Error creating writer for identifier %s: %v\n", fileIdentifier, err)
		return
	}

	// Prepare the data to write
	record := []string{
		tradeData.Timestamp.Format(time.RFC3339), // Format timestamp as ISO 8601
		helper.FormatFloat(tradeData.Price),
		helper.FormatFloat(tradeData.Quantity),
	}

	// Write the record
	if err := writer.Write(record); err != nil {
		log.Printf("Error writing record to CSV file for %s: %v\n", fileIdentifier, err)
	}
	writer.Flush() // Ensure data is written to disk
}

// getOrCreateWriter checks if a writer for the given identifier exists; if not, creates one
func (s *CSVTradeStorer) getOrCreateWriter(fileIdentifier string) (*csv.Writer, error) {
	if writer, exists := s.writers[fileIdentifier]; exists {
		return writer, nil
	}

	// Generate a unique file path for this session with a timestamp in the filename
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.csv", fileIdentifier, timestamp)
	filePath := filepath.Join(s.dirPath, fileName)

	// Open a new file for the given identifier
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to open or create file: %w", err)
	}

	// Write the header if the file is new
	writer := csv.NewWriter(file)
	if fileStat, _ := file.Stat(); fileStat.Size() == 0 {
		header := []string{"Timestamp", "Price", "Quantity"}
		if err := writer.Write(header); err != nil {
			log.Printf("Error writing header to CSV file: %v\n", err)
		}
		writer.Flush()
	}

	// Handle the file and writer for the duration of the session
	s.files[fileIdentifier] = file
	s.writers[fileIdentifier] = writer

	return writer, nil
}

func (s *CSVTradeStorer) Start() error {
	return nil
}

// Stop releases all open files and writers, called at the end of the session
func (s *CSVTradeStorer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for identifier, file := range s.files {
		log.Printf("Closing file for identifier %s", identifier)
		err := file.Close()
		if err != nil {
			return err
		}
		delete(s.files, identifier)
		delete(s.writers, identifier)
	}
	return nil
}

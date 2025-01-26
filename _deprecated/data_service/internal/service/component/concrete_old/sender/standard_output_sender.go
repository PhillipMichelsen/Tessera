package sender

import (
	"AlgorithimcTraderDistributed/common/models"
	"fmt"
	"time"
)

type StandardOutputSender struct {
	Verbose bool
}

func NewStandardOutputSender(verbose bool) *StandardOutputSender {
	return &StandardOutputSender{
		Verbose: verbose,
	}
}

func (s *StandardOutputSender) Send(marketDataPiece models.MarketDataPiece) {
	marketDataPiece.SendTimestamp = time.Now().UTC()
	if !s.Verbose {
		marketDataPiece.Payload = "*Omitted*"

	}
	fmt.Printf("StandardOutputSender: %+v\n", marketDataPiece)
}

func (s *StandardOutputSender) Start() error {
	return nil
}

func (s *StandardOutputSender) Stop() error {
	return nil
}

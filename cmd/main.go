package main

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules/concrete"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

func main() {
	// Output and update alert channels
	outputChannel := make(chan interface{}, 100)
	updateAlertChannel := make(chan uuid.UUID, 10)

	var wg sync.WaitGroup

	go func() {
		defer wg.Done() // Mark output goroutine as done
		for output := range outputChannel {
			switch output.(type) {
			case models.MarketDataPacket:
			default:
				fmt.Printf("Unknown output type: %T\n", output)
			}
		}
	}()

	go func() {
		for alert := range updateAlertChannel {
			fmt.Printf("Received update alert for module: %s\n", alert)
		}
	}()

	// Define rawConfig with mappings (simulated YAML or JSON input)
	rawConfig := map[string]interface{}{
		"input_output_uuid_mappings": map[interface{}]interface{}{
			"00000000-0000-0000-0000-000000000000": "11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222": "33333333-3333-3333-3333-333333333333",
		},
	}

	// Create a new transformer module
	transformer := concrete.NewBinanceSpotKlineToOHLCVTransformer(uuid.New(), rawConfig, outputChannel, updateAlertChannel)

	// Initialize and start the module
	transformer.Initialize()
	transformer.Start()

	// Add a count to the WaitGroup for input and output processing
	const numPackets = 10
	wg.Add(2) // One for input goroutine, one for output goroutine

	// Simulate sending packets to the input channel
	go func() {
		defer wg.Done() // Mark input goroutine as done
		inputChannel := transformer.GetInputChannel()
		for i := 0; i < numPackets; i++ {
			packetUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
			inputChannel <- models.MarketDataPacket{
				CurrentUUID:           packetUUID,
				UUIDHistory:           nil,
				UUIDHistoryTimestamps: nil,
				Payload: models.SerializedJSON{
					JSON: fmt.Sprintf(`{
						"E": %d,
						"k": {
							"t": %d,
							"o": 100.0,
							"h": 200.0,
							"l": 90.0,
							"c": 150.0,
							"v": 5000.0
						}
					}`, time.Now().UnixMilli(), time.Now().UnixMilli()),
				},
			}
		}
		close(inputChannel) // Close input channel after sending packets
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Stop the module and close remaining channels
	transformer.Stop()
	close(outputChannel)
	close(updateAlertChannel)

	fmt.Println("All packets processed, program terminated.")
}

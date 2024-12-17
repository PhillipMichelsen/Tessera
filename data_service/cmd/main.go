package main

import (
	"AlgorithimcTraderDistributed/data_service/internal/service/manager"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mgr := manager.NewManager(500, 2)

	// Define a map for receivers and their configurations
	receivers := map[string]map[string]string{
		"BinanceSpotWebsocketReceiver": {
			"streams":  "btcusdt@kline_1s",
			"base_url": "stream.binance.com:9443",
		},
	}

	// Define a map for components and their configurations
	components := map[string]map[string]string{
		"BinanceSpotKlineToOHLCVTransformerProcessor": {},
	}

	// Register all receivers from the map
	for name, config := range receivers {
		if err := mgr.CreateAndRegisterReceiver(name, mgr.Router.Route, config); err != nil {
			log.Printf("Error registering receiver %s: %v", name, err)
			return
		}
	}

	// Register all components from the map
	for name, config := range components {
		if err := mgr.CreateAndRegisterComponent(name, mgr.Router.Route, config); err != nil {
			log.Printf("Error registering component %s: %v", name, err)
			return
		}
	}

	// START SERVICE
	if err := mgr.Start(); err != nil {
		log.Fatalf("Failed to start manager: %v", err)
	}

	// STOP SERVICE (GRACEFULLY)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Service started. Press Ctrl+C to stop.")
	<-stop

	if err := mgr.Stop(); err != nil {
		log.Fatalf("Error stopping manager: %v", err)
	}

	log.Println("Service stopped.")
}

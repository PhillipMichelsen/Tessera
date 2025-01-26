package main

import (
	"AlgorithimcTraderDistributed/strategy_service/internal/manager"
	strategy "AlgorithimcTraderDistributed/strategy_service/internal/strategy/concrete"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sessionUUID := uuid.NewString()
	mgr, err := manager.NewManager(sessionUUID, "amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Printf("failed to create mgr: %v", err)
		return
	}

	strategyUUID := uuid.NewString()
	channel, err := mgr.GetRabbitMQChannel()
	if err != nil {
		log.Printf("failed to get RabbitMQ channel: %v", err)
		return
	}
	strat := strategy.NewOrderFlowImbalanceStrategy(strategyUUID, sessionUUID, channel)
	err = strat.Initialize(map[string]string{
		"zScorePeriod": "30",
		"asset":        "BinanceSpot:btcusdt",
	})
	if err != nil {
		log.Printf("failed to initialize strategy %v", err)
		return
	}

	mgr.AddStrategy(strategyUUID, strat)
	mgr.StartAllStrategiesAndAutoLive()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Service started. Press Ctrl+C to stop.")
	<-stop

	log.Println("Service stopped.")
}

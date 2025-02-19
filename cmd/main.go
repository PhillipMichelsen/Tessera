package main

import (
	"AlgorithmicTraderDistributed/internal/instance"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	inst := instance.NewInstance()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Println("[INFO] *Main* | Starting instance...")
	<-signals
	log.Println("[INFO] *Main* | Received signal, shutting down...")
	inst.Shutdown()
	log.Println("[INFO] *Main* | Shutdown complete.")
}

package main

import (
	"AlgorithmicTraderDistributed/internal/instance"
	"AlgorithmicTraderDistributed/internal/instance/controllers"
	"AlgorithmicTraderDistributed/internal/modules"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func main() {
	inst := instance.NewInstance()

	module := modules.NewModule(uuid.Must(uuid.NewUUID()), nil, inst)
	inst.AddModule(module)

	tuiController := controllers.NewTUIController(inst)
	go tuiController.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Println("[INFO] *Main* | Starting instance...")
	<-signals
	log.Println("[INFO] *Main* | Received signal, shutting down...")
	inst.Shutdown()
	log.Println("[INFO] *Main* | Shutdown complete.")
}

package main

import (
	"AlgorithmicTraderDistributed/internal/instance"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	inst := instance.NewInstance()
	//inst.AddController(instance.InstantiateController("tui", inst))

	inst.CreateModule("TestCore", uuid.Must(uuid.NewRandom()))

	inst.Start()
	inst.InitializeModule(inst.GetModules()[0], map[string]interface{}{})
	inst.StartModule(inst.GetModules()[0])

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	inst.Shutdown()
}

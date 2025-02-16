package instance

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"github.com/google/uuid"
	"log"
)

type Dispatcher struct {
	channelMappings map[uuid.UUID]chan<- interface{}
	buffer          chan interface{}

	stopChannel chan struct{}
}

func NewDispatcher(bufferSize int) *Dispatcher {
	return &Dispatcher{
		channelMappings: make(map[uuid.UUID]chan<- interface{}),
		buffer:          make(chan interface{}, bufferSize),
		stopChannel:     make(chan struct{}),
	}
}

func (d *Dispatcher) AddMapping(module modules.Module) {
	log.Printf("[INFO] *Dispatcher* | Adding mapping for module: %s \n", module.GetModuleUUID())
	d.channelMappings[module.GetModuleUUID()] = module.GetInputChannel()
	log.Printf("[INFO] *Dispatcher* | Mapping added for module: %s \n", module.GetModuleUUID())
}

func (d *Dispatcher) RemoveMapping(module modules.Module) {
	log.Printf("[INFO] *Dispatcher* | Removing mapping for module: %s \n", module.GetModuleUUID())
	delete(d.channelMappings, module.GetModuleUUID())
	log.Printf("[INFO] *Dispatcher* | Mapping removed for module: %s \n", module.GetModuleUUID())
}

func (d *Dispatcher) Dispatch(packet models.Packet) {
	d.buffer <- packet
}

func (d *Dispatcher) Start() {
	log.Println("[INFO] *Dispatcher* | Starting...")
	go d.runWorker()
	log.Println("[INFO] *Dispatcher* | Started.")
}

func (d *Dispatcher) Stop() {
	log.Println("[INFO] *Dispatcher* | Stopping...")
	close(d.stopChannel)
	log.Println("[INFO] *Dispatcher* | Stopped.")
}

func (d *Dispatcher) runWorker() {
	for {
		select {
		case <-d.stopChannel:
			return
		case message := <-d.buffer:
			packet, ok := message.(models.Packet)
			if !ok {
				log.Printf("[ERROR] *Dispatcher* | Invalid message received (Expected Packet type): %v \n", message)
				continue
			}
			destinationChannel, ok := d.channelMappings[packet.DestinationUUID]
			if !ok {
				log.Printf("[ERROR] *Dispatcher* | Destination channel not found for Packet: %v \n", packet)
				continue
			}

			destinationChannel <- packet.Payload
		}
	}
}

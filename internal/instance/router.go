package instance

import (
	"AlgorithmicTraderDistributed/internal/models"
	"AlgorithmicTraderDistributed/internal/modules"
	"github.com/google/uuid"
	"sync"
)

type Router struct {
	inputChan chan interface{}

	modules                       map[uuid.UUID]ModuleResources
	packetUUIDToModuleUUIDMapping map[uuid.UUID]uuid.UUID

	wg      sync.WaitGroup
	stopAll chan struct{}
}

type ModuleResources struct {
	ModuleInputChannel  chan<- interface{}
	PacketBufferChannel chan interface{}
	WorkerStopChannel   chan struct{}
}

// NewRouter initializes a new Router instance
func NewRouter(routerInputBufferSize int) *Router {
	return &Router{
		inputChan:                     make(chan interface{}, routerInputBufferSize),
		modules:                       make(map[uuid.UUID]ModuleResources),
		packetUUIDToModuleUUIDMapping: make(map[uuid.UUID]uuid.UUID),
		stopAll:                       make(chan struct{}),
	}
}

// AddModule adds a module to the router and sets up its buffer channels
func (r *Router) AddModule(module modules.Module, moduleBufferSize int) {
	r.modules[module.GetModuleUUID()] = ModuleResources{
		ModuleInputChannel:  module.GetInputChannel(),
		PacketBufferChannel: make(chan interface{}, moduleBufferSize),
		WorkerStopChannel:   make(chan struct{}),
	}
}

// Start initializes the main router loop
func (r *Router) Start() {
	// For each module, start a worker
	for moduleUUID := range r.modules {
		r.startModuleWorker(moduleUUID)
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case data := <-r.inputChan:
				r.handlePacket(data)
			case <-r.stopAll:
				return
			}
		}
	}()
}

// Stop stops the router and all its workers gracefully
func (r *Router) Stop() {
	close(r.stopAll)
	r.wg.Wait()
	for _, module := range r.modules {
		close(module.WorkerStopChannel)
		close(module.PacketBufferChannel)
	}
}

// startModuleWorker starts a worker for a specific module
func (r *Router) startModuleWorker(moduleUUID uuid.UUID) {
	moduleResources, exists := r.modules[moduleUUID]
	if !exists {
		panic("Cannot start worker: no resources for module UUID: " + moduleUUID.String())
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runWorker(moduleResources, moduleUUID)
	}()
}

// handlePacket processes a single input packet
func (r *Router) handlePacket(data interface{}) {
	var packetUUID uuid.UUID

	// Validate the packet type
	switch packet := data.(type) {
	case models.MarketDataPacket:
		packetUUID = packet.CurrentUUID
	default:
		panic("Invalid packet type received in router")
	}

	// Use the packet UUID to route to the appropriate module
	moduleUUID, exists := r.packetUUIDToModuleUUIDMapping[packetUUID]
	if !exists {
		panic("No module mapping found for packet UUID: " + packetUUID.String())
	}

	moduleResources, exists := r.modules[moduleUUID]
	if !exists {
		panic("No module resources found for module UUID: " + moduleUUID.String())
	}

	// Send the validated packet to the module's buffer channel
	select {
	case moduleResources.PacketBufferChannel <- data:
		// Successfully pushed to the buffer
	default:
		panic("Buffer full for module UUID: " + moduleUUID.String())
	}
}

// runWorker processes packets from the buffer channel to the module input channel
func (r *Router) runWorker(resources ModuleResources, moduleUUID uuid.UUID) {
	for {
		select {
		case packet := <-resources.PacketBufferChannel:
			// Send packet to the module's input channel
			select {
			case resources.ModuleInputChannel <- packet:
				// Successfully forwarded to the module
			default:
				panic("Module input channel full for module UUID: " + moduleUUID.String())
			}
		case <-resources.WorkerStopChannel:
			return
		case <-r.stopAll:
			return
		}
	}
}

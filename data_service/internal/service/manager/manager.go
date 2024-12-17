package manager

import (
	"AlgorithimcTraderDistributed/common/models"
	"AlgorithimcTraderDistributed/data_service/internal/service/component"

	"github.com/google/uuid"
)

type Manager struct {
	Router     *Router
	Components map[string]map[uuid.UUID]component.Component
}

// NewManager creates a new Manager instance.
func NewManager(routerBufferSize, routerWorkerCount int) *Manager {
	return &Manager{
		Router:     NewRouter(routerBufferSize, routerWorkerCount),
		Components: make(map[string]map[uuid.UUID]component.Component),
	}
}

// CreateAndRegisterComponent creates a component using the factory and registers it with the Manager and Router.
func (m *Manager) CreateAndRegisterComponent(componentType string, routeFunction func(marketData models.MarketDataPiece), config map[string]string) error {
	comp, err := component.CreateComponent(componentType, routeFunction, config)
	if err != nil {
		return err
	}
	m.RegisterComponent(comp)

	return nil
}

// RegisterComponent registers a new component (Processor, Sender, Storer) to the Manager and Router.
func (m *Manager) RegisterComponent(comp component.Component) {
	if _, exists := m.Components[comp.GetName()]; !exists {
		m.Components[comp.GetName()] = make(map[uuid.UUID]component.Component)
	}
	m.Components[comp.GetName()][comp.GetUUID()] = comp
	m.Router.RegisterComponent(comp)
}

// RegisterMapping adds a routing rule to the Router for a specific component.
func (m *Manager) RegisterMapping(routingPattern models.RoutingPattern, componentName string, componentUUID uuid.UUID) {
	m.Router.RegisterRoutingRule(routingPattern, componentName, componentUUID)
}

// Start starts all components and receivers.
func (m *Manager) Start() error {
	for _, components := range m.Components {
		for _, comp := range components {
			if err := comp.Start(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Stop stops all components and receivers.
func (m *Manager) Stop() error {
	for _, components := range m.Components {
		for _, comp := range components {
			if err := comp.Stop(); err != nil {
				return err
			}
		}
	}

	m.Router.Shutdown()
	return nil
}

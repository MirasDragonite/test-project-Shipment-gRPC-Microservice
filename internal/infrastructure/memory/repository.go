package memory

import (
	"context"
	"sync"
	"test-task-miras/internal/domain"
)

// InMemory realization of DB through map. It realized repository interface in Use-case 
type ShipmentRepository struct {
	mu        sync.RWMutex
	shipments map[string]*domain.Shipment
	events    []*domain.ShipmentEvent
}

func NewShipmentRepository() *ShipmentRepository {
	return &ShipmentRepository{
		shipments: make(map[string]*domain.Shipment),
		events:    make([]*domain.ShipmentEvent, 0),
	}
}

func (r *ShipmentRepository) Save(ctx context.Context, shipment *domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shipments[shipment.ID] = shipment
	return nil
}

func (r *ShipmentRepository) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shipment, exists := r.shipments[id]
	if !exists {
		return nil, domain.ErrShipmentNotFound
	}
	return shipment, nil
}

func (r *ShipmentRepository) SaveEvent(ctx context.Context, event *domain.ShipmentEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}

func (r *ShipmentRepository) GetEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.ShipmentEvent
	for _, e := range r.events {
		if e.ShipmentID == shipmentID {
			result = append(result, e)
		}
	}
	return result, nil
}

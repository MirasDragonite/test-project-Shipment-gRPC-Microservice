package usecase

import (
	"context"
	"test-task-miras/internal/domain"
)

// Repo interface, so use-case could dictate the storage conditions.
type ShipmentRepository interface {
	Save(ctx context.Context, shipment *domain.Shipment) error
	GetByID(ctx context.Context, id string) (*domain.Shipment, error)
	SaveEvent(ctx context.Context, event *domain.ShipmentEvent) error
	GetEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error)
}

// Base Use Case struct with constructor to create it.
type ShipmentUseCase struct {
	repo ShipmentRepository
}

func NewShipmentUseCase(repo ShipmentRepository) *ShipmentUseCase {
	return &ShipmentUseCase{repo: repo}
}

// Cargo creation process
func (uc *ShipmentUseCase) CreateShipment(ctx context.Context, ref, origin, dest, driver, unit string, amount, rev float64) (*domain.Shipment, error) {
	shipment, err := domain.NewShipment(ref, origin, dest, driver, unit, amount, rev)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, shipment); err != nil {
		return nil, err
	}

	return shipment, nil
}

// Cargo status updating process
func (uc *ShipmentUseCase) UpdateStatus(ctx context.Context, id string, newStatus domain.Status, note string) (*domain.Shipment, *domain.ShipmentEvent, error) {
	shipment, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	event, err := shipment.ChangeStatus(newStatus, note)
	if err != nil {
		return nil, nil, err
	}

	// saving the data, realization in two function, 
	// but if we will sure about using sql type of DB, we could realized it with transactions
	if err := uc.repo.Save(ctx, shipment); err != nil {
		return nil, nil, err
	}
	if err := uc.repo.SaveEvent(ctx, event); err != nil {
		return nil, nil, err
	}

	return shipment, event, nil
}

// Gets the shipment obj by id
func (uc *ShipmentUseCase) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	return uc.repo.GetByID(ctx, id)
}

// Gets hisoty of statuses
func (uc *ShipmentUseCase) GetHistory(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	return uc.repo.GetEventsByShipmentID(ctx, shipmentID)
}

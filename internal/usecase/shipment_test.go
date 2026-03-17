package usecase_test

import (
	"context"
	"errors"
	"test-task-miras/internal/domain"
	"test-task-miras/internal/usecase"
	"testing"
)

// mock repo to isolates Use Case tests from the actual database.
type mockShipmentRepository struct {
	mockSave                  func(ctx context.Context, shipment *domain.Shipment) error
	mockGetByID               func(ctx context.Context, id string) (*domain.Shipment, error)
	mockSaveEvent             func(ctx context.Context, event *domain.ShipmentEvent) error
	mockGetEventsByShipmentID func(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error)
}

func (m *mockShipmentRepository) Save(ctx context.Context, shipment *domain.Shipment) error {
	return m.mockSave(ctx, shipment)
}

func (m *mockShipmentRepository) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	return m.mockGetByID(ctx, id)
}

func (m *mockShipmentRepository) SaveEvent(ctx context.Context, event *domain.ShipmentEvent) error {
	return m.mockSaveEvent(ctx, event)
}

func (m *mockShipmentRepository) GetEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	return m.mockGetEventsByShipmentID(ctx, shipmentID)
}

// Testing the orchestration of shipment creation.
func TestShipmentUseCase_CreateShipment(t *testing.T) {
	// Setup mock to simulate a successful database save
	mockRepo := &mockShipmentRepository{
		mockSave: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
	}

	uc := usecase.NewShipmentUseCase(mockRepo)
	ctx := context.Background()

	shipment, err := uc.CreateShipment(ctx, "REF-001", "A", "B", "Driver", "Truck", 100, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipment == nil {
		t.Fatal("expected shipment, got nil")
	}
}

// Testing fetching, updating, and saving logic.
func TestShipmentUseCase_UpdateStatus(t *testing.T) {
	existingShipment := &domain.Shipment{
		ID:            "test-id",
		CurrentStatus: domain.StatusPending,
	}

	// Simulating successful db
	mockRepo := &mockShipmentRepository{
		mockGetByID: func(ctx context.Context, id string) (*domain.Shipment, error) {
			if id == "test-id" {
				return existingShipment, nil
			}
			return nil, errors.New("not found")
		},
		mockSave: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
		mockSaveEvent: func(ctx context.Context, event *domain.ShipmentEvent) error {
			return nil
		},
	}

	uc := usecase.NewShipmentUseCase(mockRepo)
	ctx := context.Background()

	//  Test valid transition
	updated, event, err := uc.UpdateStatus(ctx, "test-id", domain.StatusPickedUp, "Picked up by driver")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.CurrentStatus != domain.StatusPickedUp {
		t.Errorf("expected status %v, got %v", domain.StatusPickedUp, updated.CurrentStatus)
	}
	if event == nil {
		t.Error("expected event to be generated")
	}

	// Test invalid transition
	_, _, err = uc.UpdateStatus(ctx, "test-id", domain.StatusDelivered, "Skipping transit")
	if err != domain.ErrInvalidStatusTransition {
		t.Errorf("expected domain error for invalid transition, got: %v", err)
	}
}

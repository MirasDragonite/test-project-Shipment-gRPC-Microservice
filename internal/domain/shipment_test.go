package domain_test

import (
	"test-task-miras/internal/domain"
	"testing"
)

// Testing creating shipment
func TestNewShipment(t *testing.T) {
	ref := "REF-123"
	origin := "Almaty"
	destination := "Astana"

	s, err := domain.NewShipment(ref, origin, destination, "Mike Tyson", "Truck-01", 500.0, 150.0)

	if err != nil {
		t.Fatalf("Error during creating, got: %v", err)
	}

	// Checking the starting status
	if s.CurrentStatus != domain.StatusPending {
		t.Errorf("Expected starting status %v, got %v", domain.StatusPending, s.CurrentStatus)
	}

	if s.ReferenceNumber != ref {
		t.Errorf("Expected reference %s, got %s", ref, s.ReferenceNumber)
	}
}

// Testing logic of the lifecycle and status transitions
func TestShipment_ChangeStatus(t *testing.T) {
	// list of tests
	tests := []struct {
		name          string
		initialStatus domain.Status
		targetStatus  domain.Status
		wantErr       bool
	}{
		{
			name:          "Successful transition: Pending -> Picked Up",
			initialStatus: domain.StatusPending,
			targetStatus:  domain.StatusPickedUp,
			wantErr:       false,
		},
		{
			name:          "Successful transition: Picked Up -> In Transit",
			initialStatus: domain.StatusPickedUp,
			targetStatus:  domain.StatusInTransit,
			wantErr:       false,
		},
		{
			name:          "Successful transition: In Transit -> Delivered",
			initialStatus: domain.StatusInTransit,
			targetStatus:  domain.StatusDelivered,
			wantErr:       false,
		},
		{
			name:          "Error: Переход из Pending сразу в Delivered",
			initialStatus: domain.StatusPending,
			targetStatus:  domain.StatusDelivered,
			wantErr:       true,
		},
		{
			name:          "Error: Переход из Delivered обратно в Pending",
			initialStatus: domain.StatusDelivered,
			targetStatus:  domain.StatusPending,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating object for test with correct status
			s := &domain.Shipment{
				CurrentStatus: tt.initialStatus,
			}

			event, err := s.ChangeStatus(tt.targetStatus, "updating through test")

			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Rechecking if our status changed
				if s.CurrentStatus != tt.targetStatus {
					t.Errorf("Status didn't changed: expected %v, got %v", tt.targetStatus, s.CurrentStatus)
				}
				// Checking if event created correctly
				if event.NewStatus != tt.targetStatus || event.PreviousStatus != tt.initialStatus {
					t.Errorf("Event was generated incorrectly")
				}
			}
		})
	}
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusPickedUp  Status = "PICKED_UP"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
)

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	CurrentStatus   Status
	DriverDetails   string
	UnitDetails     string
	ShipmentAmount  float64
	DriverRevenue   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ShipmentEvent struct {
	ID             string
	ShipmentID     string
	PreviousStatus Status
	NewStatus      Status
	Note           string
	CreatedAt      time.Time
}

// Сreates a new shipment object with initial "pending" status.
func NewShipment(ref, origin, dest, driver, unit string, amount, revenue float64) (*Shipment, error) {
	if ref == "" || origin == "" || dest == "" {
		return nil, ErrValidation
	}

	return &Shipment{
		ID:              uuid.New().String(),
		ReferenceNumber: ref,
		Origin:          origin,
		Destination:     dest,
		CurrentStatus:   StatusPending,
		DriverDetails:   driver,
		UnitDetails:     unit,
		ShipmentAmount:  amount,
		DriverRevenue:   revenue,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// Updates the shipment status and returns a new tracking event if the transition is valid.
func (s *Shipment) ChangeStatus(newStatus Status, note string) (*ShipmentEvent, error) {
	if !s.isValidTransition(newStatus) {
		return nil, ErrInvalidStatusTransition
	}

	event := &ShipmentEvent{
		ID:             uuid.New().String(),
		ShipmentID:     s.ID,
		PreviousStatus: s.CurrentStatus,
		NewStatus:      newStatus,
		Note:           note,
		CreatedAt:      time.Now(),
	}

	s.CurrentStatus = newStatus
	s.UpdatedAt = time.Now()

	return event, nil
}

// Ensures compliance with business rules for the shipment lifecycle.
func (s *Shipment) isValidTransition(newStatus Status) bool {
	switch s.CurrentStatus {
	case StatusPending:
		return newStatus == StatusPickedUp
	case StatusPickedUp:
		return newStatus == StatusInTransit
	case StatusInTransit:
		return newStatus == StatusDelivered
	case StatusDelivered:
		return false
	default:
		return false
	}
}

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"test-task-miras/internal/domain"

	_ "github.com/lib/pq"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Save(ctx context.Context, s *domain.Shipment) error {
	query := `
		INSERT INTO shipments (
			id, reference_number, origin, destination, current_status, 
			driver_details, unit_details, shipment_amount, driver_revenue, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			current_status = EXCLUDED.current_status,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.ReferenceNumber, s.Origin, s.Destination, string(s.CurrentStatus),
		s.DriverDetails, s.UnitDetails, s.ShipmentAmount, s.DriverRevenue, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *ShipmentRepository) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	query := `
		SELECT id, reference_number, origin, destination, current_status, 
		       driver_details, unit_details, shipment_amount, driver_revenue, created_at, updated_at
		FROM shipments WHERE id = $1
	`

	s := &domain.Shipment{}
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.ReferenceNumber, &s.Origin, &s.Destination, &statusStr,
		&s.DriverDetails, &s.UnitDetails, &s.ShipmentAmount, &s.DriverRevenue, &s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}

	s.CurrentStatus = domain.Status(statusStr)
	return s, nil
}

func (r *ShipmentRepository) SaveEvent(ctx context.Context, e *domain.ShipmentEvent) error {
	query := `
		INSERT INTO shipment_events (id, shipment_id, previous_status, new_status, note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.ShipmentID, string(e.PreviousStatus), string(e.NewStatus), e.Note, e.CreatedAt,
	)
	return err
}

func (r *ShipmentRepository) GetEventsByShipmentID(ctx context.Context, id string) ([]*domain.ShipmentEvent, error) {
	query := `
		SELECT id, shipment_id, previous_status, new_status, note, created_at 
		FROM shipment_events 
		WHERE shipment_id = $1 
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.ShipmentEvent
	for rows.Next() {
		e := &domain.ShipmentEvent{}
		var prevStr, newStr string

		if err := rows.Scan(&e.ID, &e.ShipmentID, &prevStr, &newStr, &e.Note, &e.CreatedAt); err != nil {
			return nil, err
		}

		e.PreviousStatus = domain.Status(prevStr)
		e.NewStatus = domain.Status(newStr)
		events = append(events, e)
	}

	return events, rows.Err()
}

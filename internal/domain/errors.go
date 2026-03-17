package domain

import "errors"

var (
	ErrShipmentNotFound        = errors.New("shipment not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrValidation              = errors.New("validation error")
)

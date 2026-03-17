package grpc

import (
	"context"
	"test-task-miras/internal/domain"
	"test-task-miras/internal/infrastructure/logger"
	"test-task-miras/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "test-task-miras/api/proto/gen"
)

// Handler struct with generated gRPC server interface.
type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer
	usecase *usecase.ShipmentUseCase
}

// Constructor to create new handler
func NewShipmentHandler(uc *usecase.ShipmentUseCase) *ShipmentHandler {
	return &ShipmentHandler{usecase: uc}
}

// CreateShipment handles the incoming request to create a new shipment.
func (h *ShipmentHandler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	shipment, err := h.usecase.CreateShipment(
		ctx, req.ReferenceNumber, req.Origin, req.Destination,
		req.DriverDetails, req.UnitDetails, req.ShipmentAmount, req.DriverRevenue,
	)

	// Translate domain errors into gRPC status codes
	if err != nil {
		if err == domain.ErrValidation {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		logger.L.Error("failed to create shipment in db", "error", err)
		return nil, status.Error(codes.Internal, "failed to create shipment")
	}

	return &pb.CreateShipmentResponse{
		Shipment: mapShipmentToPB(shipment),
	}, nil
}

// Updating shipment status transitions.
func (h *ShipmentHandler) UpdateShipmentStatus(ctx context.Context, req *pb.UpdateShipmentStatusRequest) (*pb.UpdateShipmentStatusResponse, error) {
	// Convert gRPC status to domain status
	domainStatus := mapStatusToDomain(req.NewStatus)

	shipment, event, err := h.usecase.UpdateStatus(ctx, req.Id, domainStatus, req.Note)
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if err == domain.ErrInvalidStatusTransition {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		logger.L.Error("failed to update status in db", "error", err)
		return nil, status.Error(codes.Internal, "failed to update status")
	}

	return &pb.UpdateShipmentStatusResponse{
		Shipment: mapShipmentToPB(shipment),
		Event:    mapEventToPB(event),
	}, nil
}

// GetShipment retrieves a single shipment by ID.
func (h *ShipmentHandler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	shipment, err := h.usecase.GetShipment(ctx, req.Id)
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		logger.L.Error("failed to get shipment from db", "error", err)
		return nil, status.Error(codes.Internal, "failed to get shipment")
	}

	return &pb.GetShipmentResponse{
		Shipment: mapShipmentToPB(shipment),
	}, nil
}

// GetShipmentHistory retrieves the event log for a shipment.
func (h *ShipmentHandler) GetShipmentHistory(ctx context.Context, req *pb.GetShipmentHistoryRequest) (*pb.GetShipmentHistoryResponse, error) {
	events, err := h.usecase.GetHistory(ctx, req.Id)
	if err != nil {
		logger.L.Error("failed to get history from db", "error", err)
		return nil, status.Error(codes.Internal, "failed to get history")
	}

	var pbEvents []*pb.ShipmentEvent
	for _, e := range events {
		pbEvents = append(pbEvents, mapEventToPB(e))
	}

	return &pb.GetShipmentHistoryResponse{
		Events: pbEvents,
	}, nil
}

// =====================================================================
// Mappers: Isolating Domain from Protobuf
// =====================================================================

func mapShipmentToPB(s *domain.Shipment) *pb.Shipment {
	return &pb.Shipment{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		CurrentStatus:   mapStatusToPB(s.CurrentStatus),
		DriverDetails:   s.DriverDetails,
		UnitDetails:     s.UnitDetails,
		ShipmentAmount:  s.ShipmentAmount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       timestamppb.New(s.CreatedAt),
		UpdatedAt:       timestamppb.New(s.UpdatedAt),
	}
}

func mapEventToPB(e *domain.ShipmentEvent) *pb.ShipmentEvent {
	return &pb.ShipmentEvent{
		Id:             e.ID,
		ShipmentId:     e.ShipmentID,
		PreviousStatus: mapStatusToPB(e.PreviousStatus),
		NewStatus:      mapStatusToPB(e.NewStatus),
		Note:           e.Note,
		CreatedAt:      timestamppb.New(e.CreatedAt),
	}
}

func mapStatusToPB(s domain.Status) pb.Status {
	switch s {
	case domain.StatusPending:
		return pb.Status_STATUS_PENDING
	case domain.StatusPickedUp:
		return pb.Status_STATUS_PICKED_UP
	case domain.StatusInTransit:
		return pb.Status_STATUS_IN_TRANSIT
	case domain.StatusDelivered:
		return pb.Status_STATUS_DELIVERED
	default:
		return pb.Status_STATUS_UNSPECIFIED
	}
}

func mapStatusToDomain(s pb.Status) domain.Status {
	switch s {
	case pb.Status_STATUS_PENDING:
		return domain.StatusPending
	case pb.Status_STATUS_PICKED_UP:
		return domain.StatusPickedUp
	case pb.Status_STATUS_IN_TRANSIT:
		return domain.StatusInTransit
	case pb.Status_STATUS_DELIVERED:
		return domain.StatusDelivered
	default:
		return ""
	}
}

package court

import (
	"context"
	"fmt"
	"time"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type useCase struct {
	courtRepo   interfaces.CourtRepository
	venueRepo   interfaces.VenueRepository
	bookingRepo interfaces.BookingRepository
}

func NewCourtUseCase(
	courtRepo interfaces.CourtRepository,
	venueRepo interfaces.VenueRepository,
	bookingRepo interfaces.BookingRepository,
) UseCase {
	return &useCase{
		courtRepo:   courtRepo,
		venueRepo:   venueRepo,
		bookingRepo: bookingRepo,
	}
}

func (uc *useCase) CreateCourt(ctx context.Context, req requests.CreateCourtRequest) (*responses.CourtResponse, error) {
	// Validate venue exists and is active
	venueID, err := uuid.Parse(req.VenueID)
	if err != nil {
		return nil, fmt.Errorf("invalid venue ID: %w", err)
	}

	venue, err := uc.venueRepo.GetByID(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("venue not found: %w", err)
	}

	if venue.Status != models.VenueStatusActive {
		return nil, fmt.Errorf("cannot create court for inactive venue")
	}

	court := &models.Court{
		ID:           uuid.New(),
		VenueID:      venueID,
		Name:         req.Name,
		Description:  req.Description,
		PricePerHour: req.PricePerHour,
		Status:       models.CourtStatusAvailable,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.courtRepo.Create(ctx, court); err != nil {
		return nil, fmt.Errorf("failed to create court: %w", err)
	}

	// Get complete court details
	createdCourt, err := uc.courtRepo.GetByID(ctx, court.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created court: %w", err)
	}

	return uc.toCourtResponse(createdCourt), nil
}

func (uc *useCase) GetCourt(ctx context.Context, id uuid.UUID) (*responses.CourtResponse, error) {
	court, err := uc.courtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("court not found: %w", err)
	}

	return uc.toCourtResponse(court), nil
}

func (uc *useCase) UpdateCourt(ctx context.Context, id uuid.UUID, req requests.UpdateCourtRequest) (*responses.CourtResponse, error) {
	court, err := uc.courtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("court not found: %w", err)
	}

	if req.Name != "" {
		court.Name = req.Name
	}
	if req.Description != "" {
		court.Description = req.Description
	}
	if req.PricePerHour > 0 {
		court.PricePerHour = req.PricePerHour
	}
	if req.Status != "" {
		court.Status = models.CourtStatus(req.Status)
	}

	court.UpdatedAt = time.Now()

	if err := uc.courtRepo.Update(ctx, court); err != nil {
		return nil, fmt.Errorf("failed to update court: %w", err)
	}

	return uc.toCourtResponse(court), nil
}

func (uc *useCase) DeleteCourt(ctx context.Context, id uuid.UUID) error {
	// Check if court has any future bookings
	now := time.Now()
	bookings, err := uc.bookingRepo.GetCourtBookings(ctx, id, now)
	if err != nil {
		return fmt.Errorf("failed to check court bookings: %w", err)
	}

	for _, booking := range bookings {
		if booking.Status != models.BookingStatusCancelled {
			return fmt.Errorf("cannot delete court with active bookings")
		}
	}

	if err := uc.courtRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete court: %w", err)
	}

	return nil
}

func (uc *useCase) ListCourts(ctx context.Context, req requests.ListCourtsRequest) (*responses.CourtListResponse, error) {
	filters := make(map[string]interface{})

	if req.VenueID != "" {
		venueID, err := uuid.Parse(req.VenueID)
		if err != nil {
			return nil, fmt.Errorf("invalid venue ID: %w", err)
		}
		filters["venue_id"] = venueID
	}

	if req.Status != "" {
		filters["status"] = models.CourtStatus(req.Status)
	}

	if req.Location != "" {
		filters["location"] = req.Location
	}

	if req.PriceMin > 0 {
		filters["price_min"] = req.PriceMin
	}

	if req.PriceMax > 0 {
		filters["price_max"] = req.PriceMax
	}

	// Get total count
	total, err := uc.courtRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Set pagination
	limit := 10
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	offset := 0
	if req.Offset > 0 {
		offset = req.Offset
	}

	// Get courts
	courts, err := uc.courtRepo.List(ctx, filters, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list courts: %w", err)
	}

	// Convert to response
	courtResponses := make([]responses.CourtResponse, len(courts))
	for i, court := range courts {
		courtResponses[i] = *uc.toCourtResponse(&court)
	}

	return &responses.CourtListResponse{
		Courts: courtResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (uc *useCase) GetVenueCourts(ctx context.Context, venueID uuid.UUID) ([]responses.CourtResponse, error) {
	// Validate venue exists
	venue, err := uc.venueRepo.GetByID(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("venue not found: %w", err)
	}

	courts, err := uc.courtRepo.GetByVenue(ctx, venue.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get venue courts: %w", err)
	}

	responses := make([]responses.CourtResponse, len(courts))
	for i, court := range courts {
		responses[i] = *uc.toCourtResponse(&court)
	}

	return responses, nil
}

func (uc *useCase) UpdateCourtStatus(ctx context.Context, id uuid.UUID, status string) error {

	if !isValidCourtStatus(status) {
		return fmt.Errorf("invalid court status: %s", status)
	}

	newStatus := models.CourtStatus(status)
	if newStatus == models.CourtStatusMaintenance {
		// Check for future bookings if setting to maintenance
		now := time.Now()
		bookings, err := uc.bookingRepo.GetCourtBookings(ctx, id, now)
		if err != nil {
			return fmt.Errorf("failed to check court bookings: %w", err)
		}

		for _, booking := range bookings {
			if booking.Status == models.BookingStatusConfirmed {
				return fmt.Errorf("cannot set court to maintenance: has confirmed future bookings")
			}
		}
	}

	if err := uc.courtRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		return fmt.Errorf("failed to update court status: %w", err)
	}

	return nil
}

// Helper methods

func (uc *useCase) toCourtResponse(court *models.Court) *responses.CourtResponse {
	description := ""
	if court.Description != "" {
		description = court.Description
	}

	return &responses.CourtResponse{
		ID:           court.ID.String(),
		Name:         court.Name,
		Description:  description,
		PricePerHour: court.PricePerHour,
		Status:       string(court.Status),
	}
}

func isValidCourtStatus(status string) bool {
	validStatuses := map[string]bool{
		string(models.CourtStatusAvailable):   true,
		string(models.CourtStatusOccupied):    true,
		string(models.CourtStatusMaintenance): true,
	}
	return validStatuses[status]
}

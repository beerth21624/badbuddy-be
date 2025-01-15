package interfaces

import (
	"context"
	"time"

	"badbuddy/internal/domain/models"

	"github.com/google/uuid"
)

type CourtRepository interface {
	Create(ctx context.Context, court *models.Court) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Court, error)
	GetCourtWithVenueByID(ctx context.Context, id uuid.UUID) (*models.CourtWithVenue, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Court, error)
	Update(ctx context.Context, court *models.Court) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByVenue(ctx context.Context, venueID uuid.UUID) ([]models.Court, error)
	GetCourtWithVenueByVenue(ctx context.Context, venueID uuid.UUID) ([]models.CourtWithVenue, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.CourtStatus) error
	GetAvailableCourts(ctx context.Context, venueID uuid.UUID, date time.Time, startTime, endTime time.Time) ([]models.Court, error)
	Count(ctx context.Context, filters map[string]interface{}) (int, error)
}

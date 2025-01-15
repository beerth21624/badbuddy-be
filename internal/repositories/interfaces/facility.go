package interfaces

import (
	"badbuddy/internal/domain/models"
	"context"
	"github.com/google/uuid"
)

type FacilityRepository interface {
	GetFacilities(ctx context.Context) ([]models.Facility, error)
	GetFacilityByID(ctx context.Context, id uuid.UUID) (*models.Facility, error)
	CreateFacility(ctx context.Context, facility *models.Facility) error
	UpdateFacility(ctx context.Context, facility *models.Facility) error
	DeleteFacility(ctx context.Context, id uuid.UUID) error
}

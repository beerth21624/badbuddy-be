package facility

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"context"
	"github.com/google/uuid"
)


type UseCase interface {
	ListFacilities(ctx context.Context) (*responses.FacilityListResponse, error)
	GetFacilityByID(ctx context.Context, id uuid.UUID) (*responses.FacilityResponse, error)
	CreateFacility(ctx context.Context, req requests.CreateAndUpdateFacilityRequest) (*responses.FacilityResponse, error)
	UpdateFacility(ctx context.Context, id uuid.UUID, req requests.CreateAndUpdateFacilityRequest) (*responses.FacilityResponse, error)
	DeleteFacility(ctx context.Context, id uuid.UUID) error
}

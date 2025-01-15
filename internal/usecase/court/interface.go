package court

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	CreateCourt(ctx context.Context, req requests.CreateCourtRequest) (*responses.CourtResponse, error)
	GetCourt(ctx context.Context, id uuid.UUID) (*responses.CourtResponse, error)
	UpdateCourt(ctx context.Context, id uuid.UUID, req requests.UpdateCourtRequest) (*responses.CourtResponse, error)
	DeleteCourt(ctx context.Context, id uuid.UUID) error
	ListCourts(ctx context.Context, req requests.ListCourtsRequest) (*responses.CourtListResponse, error)
	GetVenueCourts(ctx context.Context, venueID uuid.UUID) ([]responses.CourtResponse, error)
	UpdateCourtStatus(ctx context.Context, id uuid.UUID, status string) error
}

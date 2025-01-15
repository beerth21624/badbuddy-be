// internal/usecases/venue/interface.go
package venue

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	CreateVenue(ctx context.Context, ownerID uuid.UUID, req requests.CreateVenueRequest) (*responses.VenueResponse, error)
	GetVenue(ctx context.Context, id uuid.UUID) (*responses.VenueResponse, error)
	UpdateVenue(ctx context.Context, id uuid.UUID, req requests.UpdateVenueRequest) error
	ListVenues(ctx context.Context, location string, limit, offset int) ([]responses.ListVenueResponse, error)
	SearchVenues(ctx context.Context, query string, limit, offset int, minPrice int, maxPrice int, location string, facilities []string) (responses.VenueResponseDTO, error)
	AddCourt(ctx context.Context, venueID uuid.UUID, req requests.CreateCourtRequest) (*responses.CourtResponse, error)
	UpdateCourt(ctx context.Context, venueID uuid.UUID, req requests.UpdateCourtRequest) error
	DeleteCourt(ctx context.Context, venueID uuid.UUID, courtID uuid.UUID) error
	AddReview(ctx context.Context, venueID uuid.UUID, userID uuid.UUID, req requests.AddReviewRequest) error
	GetReviews(ctx context.Context, venueID uuid.UUID, limit, offset int) ([]responses.ReviewResponse, error)
	GetFacilities(ctx context.Context, venueID uuid.UUID) (*responses.FacilityListResponse, error)
	IsOwner(ctx context.Context, venueID uuid.UUID, ownerID uuid.UUID) (bool, error)
}

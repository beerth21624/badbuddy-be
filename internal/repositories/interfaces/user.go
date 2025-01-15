// interfaces/user_repository.go
package interfaces

import (
	"badbuddy/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type UserSearchFilters struct {
	PlayLevel models.PlayerLevel
	Location  string
	Limit     int
	Offset    int
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	UpdateLastActive(ctx context.Context, userID uuid.UUID) error
	SearchUsers(ctx context.Context, query string, filters UserSearchFilters) ([]models.User, error)
	GetVenueUserOwn(ctx context.Context, userID uuid.UUID) ([]models.VenueUserOwn, error)
	IsUserExist(ctx context.Context, userID uuid.UUID) (bool, error)
}

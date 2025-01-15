package user

import (
	"context"
	"errors"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrInvalidPlayLevel   = errors.New("invalid play level")
	ErrInvalidPassword    = errors.New("password does not meet requirements")
)

type UseCase interface {
	Register(ctx context.Context, req requests.RegisterRequest) error
	Login(ctx context.Context, req requests.LoginRequest) (*responses.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*responses.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req requests.UpdateProfileRequest) error
	SearchUsers(ctx context.Context, query string, filters requests.SearchFilters) ([]responses.UserResponse, error)
	RefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	GetVenueUserOwn(ctx context.Context, userID uuid.UUID) ([]responses.Venue, error)
	UpdateRoles(ctx context.Context, adminID uuid.UUID, req requests.UpdateRolesRequest) error
}

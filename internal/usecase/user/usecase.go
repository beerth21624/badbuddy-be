// usecase/user/user.go
package user

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"
	"context"
	"fmt"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type useCase struct {
	userRepo    interfaces.UserRepository
	jwtSecret   []byte
	jwtDuration time.Duration
}

func NewUserUseCase(userRepo interfaces.UserRepository, jwtSecret string, jwtDuration time.Duration) UseCase {
	return &useCase{
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
		jwtDuration: jwtDuration,
	}
}

func (uc *useCase) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidPassword)
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("%w: password must contain at least one uppercase letter, one lowercase letter, and one number",
			ErrInvalidPassword)
	}

	return nil
}

func (uc *useCase) validatePlayLevel(level string) error {
	switch models.PlayerLevel(level) {
	case models.PlayerLevelBeginner, models.PlayerLevelIntermediate, models.PlayerLevelAdvanced:
		return nil
	default:
		return ErrInvalidPlayLevel
	}
}

func (uc *useCase) Register(ctx context.Context, req requests.RegisterRequest) error {
	// Validate password
	if err := uc.validatePassword(req.Password); err != nil {
		return err
	}

	// Validate play level
	if err := uc.validatePlayLevel(req.PlayLevel); err != nil {
		return err
	}

	// Check if email exists
	if _, err := uc.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return ErrDuplicateEmail
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		PlayLevel: models.PlayerLevel(req.PlayLevel),
		Location:  req.Location,
		Bio:       req.Bio,
		Status:    models.UserStatusActive,
		CreatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (uc *useCase) generateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(uc.jwtDuration).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString(uc.jwtSecret)
}

func (uc *useCase) Login(ctx context.Context, req requests.LoginRequest) (*responses.LoginResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println(err)
		return nil, ErrInvalidCredentials
	} 

	if user.Status != models.UserStatusActive {
		return nil, fmt.Errorf("account is not active")
	}

	// Update last active
	if err := uc.userRepo.UpdateLastActive(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to update last active: %w", err)
	}

	// Generate JWT token
	tokenString, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &responses.LoginResponse{
		AccessToken: tokenString,
		User:        uc.mapUserToResponse(user),
	}, nil
}

func (uc *useCase) RefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	if user.Status != models.UserStatusActive {
		return "", fmt.Errorf("account is not active")
	}

	return uc.generateToken(userID)
}

func (uc *useCase) GetProfile(ctx context.Context, userID uuid.UUID) (*responses.UserProfileResponse, error) {
	profile, err := uc.userRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &responses.UserProfileResponse{
		UserResponse:    uc.mapUserToResponse(&profile.User),
		HostedSessions:  profile.HostedSessions,
		JoinedSessions:  profile.JoinedSessions,
		AverageRating:   profile.AverageRating,
		TotalReviews:    profile.TotalReviews,
		RegularPartners: profile.RegularPartners,
	}, nil
}

func (uc *useCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req requests.UpdateProfileRequest) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if req.PlayLevel != "" {
		if err := uc.validatePlayLevel(req.PlayLevel); err != nil {
			return err
		}
		user.PlayLevel = models.PlayerLevel(req.PlayLevel)
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Location != "" {
		user.Location = req.Location
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (uc *useCase) SearchUsers(ctx context.Context, query string, filters requests.SearchFilters) ([]responses.UserResponse, error) {
	repoFilters := interfaces.UserSearchFilters{
		PlayLevel: models.PlayerLevel(filters.PlayLevel),
		Location:  filters.Location,
		Limit:     filters.Limit,
		Offset:    filters.Offset,
	}

	users, err := uc.userRepo.SearchUsers(ctx, query, repoFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	userResponses := make([]responses.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = uc.mapUserToResponse(&user)
	}

	return userResponses, nil
}

func (uc *useCase) GetVenueUserOwn(ctx context.Context, userID uuid.UUID) ([]responses.Venue, error) {
	venues, err := uc.userRepo.GetVenueUserOwn(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get venue owners: %w", err)
	}

	venueOwners := make([]responses.Venue, len(venues))
	for i, venue := range venues {
		venueOwners[i] = responses.Venue{
			ID: venue.ID,
		}
	}

	return venueOwners, nil
}

func (uc *useCase) UpdateRoles(ctx context.Context, adminID uuid.UUID, req requests.UpdateRolesRequest) error {
	isAdmin, err := uc.IsAdmin(ctx, adminID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return fmt.Errorf("unauthorized")
	}

	user, err := uc.userRepo.GetByID(ctx, uuid.MustParse(req.UserID))
	if err != nil {
		return ErrUserNotFound
	}

	user.Role = req.Role
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (uc *useCase) mapUserToResponse(user *models.User) responses.UserResponse {
	var userID string

	if user.ID != uuid.Nil {
		userID = user.ID.String()
	}

	return responses.UserResponse{
		ID:           userID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		PlayLevel:    string(user.PlayLevel),
		Gender:       user.Gender,
		PlayHand:     user.PlayHand,
		Location:     user.Location,
		Bio:          user.Bio,
		AvatarURL:    user.AvatarURL,
		LastActiveAt: user.LastActiveAt,
		Role:         user.Role,
	}
}

func (uc *useCase) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, ErrUserNotFound
	}

	return user.Role == string(models.UserRoleAdmin), nil
}

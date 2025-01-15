package venue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
)

type useCase struct {
	venueRepo interfaces.VenueRepository
	userRepo  interfaces.UserRepository
}

func NewVenueUseCase(venueRepo interfaces.VenueRepository, userRepo interfaces.UserRepository) UseCase {
	return &useCase{
		venueRepo: venueRepo,
		userRepo:  userRepo,
	}
}

func (uc *useCase) CreateVenue(ctx context.Context, ownerID uuid.UUID, req requests.CreateVenueRequest) (*responses.VenueResponse, error) {

	venue := &models.Venue{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Location:    req.Location,
		Phone:       req.Phone,
		Email:       req.Email,
		OpenRange:   models.NullRawMessage{RawMessage: mustMarshalJSON(req.OpenRange)},
		Rules:       models.NullRawMessage{RawMessage: mustMarshalJSON(req.Rules)},
		ImageURLs:   req.ImageURLs,
		Status:      models.VenueStatusActive,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := uc.venueRepo.Create(ctx, venue); err != nil {
		return nil, fmt.Errorf("failed to create venue: %w", err)
	}

	facilityUUIDs := make([]uuid.UUID, len(req.Facilities))
	for i, facility := range req.Facilities {
		facilityUUID, err := uuid.Parse(facility.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid facility ID: %w", err)
		}
		facilityUUIDs[i] = facilityUUID
	}

	if err := uc.venueRepo.AddFacilities(ctx, venue.ID, facilityUUIDs); err != nil {
		return nil, fmt.Errorf("failed to add facilities: %w", err)
	}

	return &responses.VenueResponse{
		ID:           venue.ID.String(),
		Name:         venue.Name,
		Description:  venue.Description,
		Address:      venue.Address,
		Location:     venue.Location,
		Phone:        venue.Phone,
		Email:        venue.Email,
		OpenRange:    convertToOpenRangeResponse(req.OpenRange),
		ImageURLs:    venue.ImageURLs,
		Status:       string(venue.Status),
		Rating:       venue.Rating,
		TotalReviews: venue.TotalReviews,
		Facilities:   convertToFacilityResponse(convertToModelFacilities(req.Facilities)),
		Rules:        convertToRuleResponse(req.Rules),
		Courts:       []responses.CourtResponse{},
		Latitude:     venue.Latitude,
		Longitude:    venue.Longitude,
	}, nil
}

func (uc *useCase) GetVenue(ctx context.Context, id uuid.UUID) (*responses.VenueResponse, error) {
	venueWithCourts, err := uc.venueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get venue: %w", err)
	}

	courts := make([]responses.CourtResponse, len(venueWithCourts.Courts))
	for i, court := range venueWithCourts.Courts {
		courts[i] = responses.CourtResponse{
			ID:           court.ID.String(),
			Name:         court.Name,
			Description:  court.Description,
			PricePerHour: court.PricePerHour,
			Status:       string(court.Status),
		}
	}

	openRange := []responses.OpenRangeResponse{}
	if unMarshalJSON(venueWithCourts.OpenRange.RawMessage, &openRange) != nil {
		return nil, fmt.Errorf("error decoding enroll response: %v", err)
	}
	rules := []responses.RuleResponse{}
	if unMarshalJSON(venueWithCourts.Rules.RawMessage, &rules) != nil {
		return nil, fmt.Errorf("error decoding enroll response: %v", err)
	}

	return &responses.VenueResponse{
		ID:           venueWithCourts.ID.String(),
		Name:         venueWithCourts.Name,
		Description:  venueWithCourts.Description,
		Address:      venueWithCourts.Address,
		Location:     venueWithCourts.Location,
		Phone:        venueWithCourts.Phone,
		Email:        venueWithCourts.Email,
		OpenRange:    openRange,
		ImageURLs:    venueWithCourts.ImageURLs,
		Status:       string(venueWithCourts.Status),
		Rating:       venueWithCourts.Rating,
		TotalReviews: venueWithCourts.TotalReviews,
		Courts:       courts,
		Facilities:   convertToFacilityResponse(venueWithCourts.Facilities),
		Rules:        rules,
		Latitude:     venueWithCourts.Latitude,
		Longitude:    venueWithCourts.Longitude,
	}, nil
}

func (uc *useCase) UpdateVenue(ctx context.Context, id uuid.UUID, req requests.UpdateVenueRequest) error {
	venue, err := uc.venueRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get venue: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		venue.Name = req.Name
	}
	if req.Description != "" {
		venue.Description = req.Description
	}
	if req.Address != "" {
		venue.Address = req.Address
	}

	if req.Phone != "" {
		venue.Phone = req.Phone
	}
	if req.Email != "" {
		venue.Email = req.Email
	}
	if req.OpenRange != nil {
		openRangeJSON, err := json.Marshal(req.OpenRange)
		if err != nil {
			return fmt.Errorf("failed to marshal open range: %w", err)
		}
		venue.OpenRange.RawMessage = openRangeJSON
	}
	if req.ImageURLs != "" {
		venue.ImageURLs = req.ImageURLs
	}
	if req.Status != "" {
		venue.Status = models.VenueStatus(req.Status)
	}
	if req.Rules != nil {
		rulesJSON, err := json.Marshal(req.Rules)
		if err != nil {
			return fmt.Errorf("failed to marshal rules: %w", err)
		}
		venue.Rules.RawMessage = rulesJSON
	}
	venue.Latitude = req.Latitude
	venue.Longitude = req.Longitude

	facilityUUIDs := make([]uuid.UUID, len(req.Facilities))
	for i, facility := range req.Facilities {
		facilityUUID, err := uuid.Parse(facility.ID)
		if err != nil {
			return fmt.Errorf("invalid facility ID: %w", err)
		}
		facilityUUIDs[i] = facilityUUID
	}

	if err := uc.venueRepo.UpdateFacilities(ctx, venue.ID, facilityUUIDs); err != nil {
		return fmt.Errorf("failed to update facilities: %w", err)
	}

	venue.UpdatedAt = time.Now()
	if err := uc.venueRepo.Update(ctx, &venue.Venue); err != nil {
		return fmt.Errorf("failed to update venue: %w", err)
	}

	return nil
}

func (uc *useCase) ListVenues(ctx context.Context, location string, limit, offset int) ([]responses.ListVenueResponse, error) {
	venues, err := uc.venueRepo.List(ctx, location, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list venues: %w", err)
	}

	venueResponses := make([]responses.ListVenueResponse, 0)

	for _, venue := range venues {
		venueResponses = append(venueResponses, responses.ListVenueResponse{
			ID:   venue.ID.String(),
			Name: venue.Name,
		})
	}
	return venueResponses, nil
}

func (uc *useCase) SearchVenues(ctx context.Context, query string, limit, offset int, minPrice int, maxPrice int, location string, facilities []string) (responses.VenueResponseDTO, error) {
	venues, err := uc.venueRepo.Search(ctx, query, limit, offset, minPrice, maxPrice, location, facilities)
	if err != nil {
		return responses.VenueResponseDTO{}, fmt.Errorf("failed to search venues: %w", err)
	}

	venueResponses := make([]responses.VenueResponse, len(venues))
	for i, venue := range venues {
		venueResponses[i] = responses.VenueResponse{
			ID:          venue.ID.String(),
			Name:        venue.Name,
			Description: venue.Description,
			Address:     venue.Address,
			Location:    venue.Location,
			Phone:       venue.Phone,
			Email:       venue.Email,
			OpenRange: func() []responses.OpenRangeResponse {
				var openRange []responses.OpenRangeResponse
				if err := unMarshalJSON(venue.OpenRange.RawMessage, &openRange); err != nil {
					return nil
				}
				return openRange
			}(),
			ImageURLs:    venue.ImageURLs,
			Status:       string(venue.Status),
			Rating:       venue.Rating,
			TotalReviews: venue.TotalReviews,
			Facilities:   convertToFacilityResponse(venue.Facilities),
			Rules: func() []responses.RuleResponse {
				var rules []responses.RuleResponse
				if err := unMarshalJSON(venue.Rules.RawMessage, &rules); err != nil {
					return nil
				}
				return rules
			}(),
			Courts:    convertToCourtResponse(venue.Courts),
			Latitude:  venue.Latitude,
			Longitude: venue.Longitude,
		}
	}

	// total, err := uc.venueRepo.CountVenues(ctx)
	// if err != nil {
	// 	return responses.VenueResponseDTO{}, fmt.Errorf("failed to count venues: %w", err)
	// }

	total, err := uc.venueRepo.CountSearch(ctx, query, minPrice, maxPrice, location, facilities)
	if err != nil {
		return responses.VenueResponseDTO{}, fmt.Errorf("failed to count venues: %w", err)
	}

	return responses.VenueResponseDTO{
		Venues: venueResponses,
		Total:  total,
	}, nil
}

func (uc *useCase) AddCourt(ctx context.Context, venueID uuid.UUID, req requests.CreateCourtRequest) (*responses.CourtResponse, error) {

	courts, err := uc.venueRepo.GetCourts(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courts: %w", err)
	}

	for _, court := range courts {
		if court.Name == req.Name {
			return nil, fmt.Errorf("court name already exists")
		}
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

	if err := uc.venueRepo.AddCourt(ctx, court); err != nil {
		return nil, fmt.Errorf("failed to add court: %w", err)
	}

	return &responses.CourtResponse{
		ID:           court.ID.String(),
		Name:         court.Name,
		Description:  court.Description,
		PricePerHour: court.PricePerHour,
		Status:       string(court.Status),
	}, nil
}

func (uc *useCase) UpdateCourt(ctx context.Context, venueID uuid.UUID, req requests.UpdateCourtRequest) error {

	courts, err := uc.venueRepo.GetCourts(ctx, venueID)
	if err != nil {
		return fmt.Errorf("failed to get court: %w", err)
	}
	courtUUID, err := uuid.Parse(req.CourtID)
	if err != nil {
		return fmt.Errorf("invalid court ID: %w", err)
	}

	var court *models.Court
	for i := range courts {

		if courts[i].ID == courtUUID {
			court = &courts[i]
			break
		}
	}

	if court == nil {
		return fmt.Errorf("court not found")
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

	if err := uc.venueRepo.UpdateCourt(ctx, court); err != nil {
		return fmt.Errorf("failed to update court: %w", err)
	}

	return nil
}

func (uc *useCase) DeleteCourt(ctx context.Context, venueID uuid.UUID, courtID uuid.UUID) error {

	courts, err := uc.venueRepo.GetCourts(ctx, venueID)
	if err != nil {
		return fmt.Errorf("failed to get court: %w", err)
	}

	var court *models.Court
	for i := range courts {
		if courts[i].ID == courtID {
			court = &courts[i]
			break
		}
	}

	if court == nil {
		return fmt.Errorf("court not found")
	}

	if err := uc.venueRepo.DeleteCourt(ctx, courtID); err != nil {
		return fmt.Errorf("failed to delete court: %w", err)
	}

	return nil

}

func (uc *useCase) AddReview(ctx context.Context, venueID uuid.UUID, userID uuid.UUID, req requests.AddReviewRequest) error {
	review := &models.VenueReview{
		ID:        uuid.New(),
		VenueID:   venueID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
	}

	fmt.Println("review added before")

	if err := uc.venueRepo.AddReview(ctx, review); err != nil {
		return fmt.Errorf("failed to add review: %w", err)
	}

	fmt.Println("review added")

	return nil
}
func (uc *useCase) GetReviews(ctx context.Context, venueID uuid.UUID, limit, offset int) ([]responses.ReviewResponse, error) {
	// Input validation
	if venueID == uuid.Nil {
		return nil, fmt.Errorf("invalid venue ID")
	}

	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}

	// Get reviews
	reviews, err := uc.venueRepo.GetReviews(ctx, venueID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	// Handle empty results
	if len(reviews) == 0 {
		return []responses.ReviewResponse{}, nil
	}

	// Collect all unique user IDs
	userIDs := make(map[uuid.UUID]struct{})
	for _, review := range reviews {
		userIDs[review.UserID] = struct{}{}
	}

	// Convert map to slice for batch query
	uniqueUserIDs := make([]uuid.UUID, 0, len(userIDs))
	for userID := range userIDs {
		uniqueUserIDs = append(uniqueUserIDs, userID)
	}

	// Batch fetch all users
	users, err := uc.userRepo.GetUsersByIDs(ctx, uniqueUserIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}

	// Create user map for efficient lookups
	userMap := make(map[uuid.UUID]models.User, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	reviewResponses := make([]responses.ReviewResponse, len(reviews))

	for i, review := range reviews {
		user, exists := userMap[review.UserID]
		if !exists {
			return nil, fmt.Errorf("user not found for review %s", review.ID)
		}

		reviewResponses[i] = responses.ReviewResponse{
			ID:        review.ID.String(),
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format(time.RFC3339),
			Reviewer: responses.ReviewerResponse{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				AvatarURL: user.AvatarURL,
			},
		}
	}

	return reviewResponses, nil
}

func (uc *useCase) GetFacilities(ctx context.Context, venueID uuid.UUID) (*responses.FacilityListResponse, error) {
	facilities, err := uc.venueRepo.GetFacilities(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %w", err)
	}

	facilityResponses := make([]responses.FacilityResponse, len(facilities))
	for i := range facilities {
		facilityResponses[i] = responses.FacilityResponse{
			ID:   facilities[i].ID.String(),
			Name: facilities[i].Name,
		}
	}

	return &responses.FacilityListResponse{
		Facilities: facilityResponses,
	}, nil
}

func (uc *useCase) IsOwner(ctx context.Context, venueID uuid.UUID, ownerID uuid.UUID) (bool, error) {
	venue, err := uc.venueRepo.GetByID(ctx, venueID)
	if err != nil {
		return false, fmt.Errorf("failed to get venue: %w", err)
	}

	return venue.OwnerID == ownerID, nil
}

func convertToOpenRangeResponse(openRanges []requests.OpenRange) []responses.OpenRangeResponse {
	var openRangeResponses []responses.OpenRangeResponse
	for _, openRange := range openRanges {
		openRangeResponses = append(openRangeResponses, responses.OpenRangeResponse{
			Day:       openRange.Day,
			IsOpen:    openRange.IsOpen,
			OpenTime:  openRange.OpenTime,
			CloseTime: openRange.CloseTime,
		})
	}
	return openRangeResponses
}

func convertToRuleResponse(rules []requests.Rule) []responses.RuleResponse {
	ruleResponses := make([]responses.RuleResponse, len(rules))
	for i, rule := range rules {
		ruleResponses[i] = responses.RuleResponse{
			Rule: rule.Rule,
		}
	}
	return ruleResponses
}

func convertToModelFacilities(facilities []requests.Facility) []models.Facility {
	modelFacilities := make([]models.Facility, len(facilities))
	for i, facility := range facilities {
		modelFacilities[i] = models.Facility{
			ID:   uuid.MustParse(facility.ID),
			Name: "",
		}
	}
	return modelFacilities
}

func convertToFacilityResponse(facilities []models.Facility) []responses.FacilityResponse {
	facilityResponses := make([]responses.FacilityResponse, len(facilities))
	for i, facility := range facilities {
		facilityResponses[i] = responses.FacilityResponse{
			ID:   facility.ID.String(),
			Name: facility.Name,
		}
	}
	return facilityResponses
}

func convertToCourtResponse(courts []models.Court) []responses.CourtResponse {
	courtResponses := make([]responses.CourtResponse, len(courts))
	for i, court := range courts {
		courtResponses[i] = responses.CourtResponse{
			ID:           court.ID.String(),
			Name:         court.Name,
			Description:  court.Description,
			PricePerHour: court.PricePerHour,
			Status:       string(court.Status),
		}
	}
	return courtResponses
}

func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}

func unMarshalJSON(data json.RawMessage, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

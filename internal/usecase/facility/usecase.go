package facility

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")

	ErrValidation = errors.New("validation error")

	ErrFacilityNotFound = errors.New("facility not found")
)

type useCase struct {
	facilityRepo interfaces.FacilityRepository
}

func NewFacilityUseCase(facilityRepo interfaces.FacilityRepository) UseCase {
	return &useCase{
		facilityRepo: facilityRepo,
	}
}

func (uc *useCase) ListFacilities(ctx context.Context) (*responses.FacilityListResponse, error) {
	facilities, err := uc.facilityRepo.GetFacilities(ctx)
	if err != nil {
		return nil, err
	}

	 
	facilityResponses := []responses.FacilityResponse{}
	for _, facility := range facilities {
		facilityResponses = append(facilityResponses, responses.FacilityResponse{
			ID:          facility.ID.String(),
			Name:        facility.Name,
		})
	}

	return &responses.FacilityListResponse{
		Facilities: facilityResponses,
	}, nil
}

func (uc *useCase) GetFacilityByID(ctx context.Context, id uuid.UUID) (*responses.FacilityResponse, error) {
	facility, err := uc.facilityRepo.GetFacilityByID(ctx, id)
	if err != nil {
		return nil, ErrFacilityNotFound
	}

	if facility == nil {
		return nil, ErrFacilityNotFound
	}

	return &responses.FacilityResponse{
		ID:          facility.ID.String(),
		Name:        facility.Name,
	}, nil
}

func (uc *useCase) CreateFacility(ctx context.Context, req requests.CreateAndUpdateFacilityRequest) (*responses.FacilityResponse, error) {
	facility := &models.Facility{
		ID:          uuid.New(),
		Name:        req.Name,
	}

	err := uc.facilityRepo.CreateFacility(ctx, facility)
	if err != nil {
		return nil, err
	}

	return &responses.FacilityResponse{
		ID:          facility.ID.String(),
		Name:        facility.Name,
	}, nil
}

func (uc *useCase) UpdateFacility(ctx context.Context, id uuid.UUID, req requests.CreateAndUpdateFacilityRequest) (*responses.FacilityResponse, error) {
	facility, err := uc.facilityRepo.GetFacilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if facility == nil {
		return nil, ErrFacilityNotFound
	}

	facility.Name = req.Name

	err = uc.facilityRepo.UpdateFacility(ctx, facility)
	if err != nil {
		return nil, err
	}

	return &responses.FacilityResponse{
		ID:          facility.ID.String(),
		Name:        facility.Name,
	}, nil
}

func (uc *useCase) DeleteFacility(ctx context.Context, id uuid.UUID) error {
	facility, err := uc.facilityRepo.GetFacilityByID(ctx, id)
	if err != nil {
		return err
	}

	if facility == nil {
		return ErrFacilityNotFound
	}

	return uc.facilityRepo.DeleteFacility(ctx, id)
}
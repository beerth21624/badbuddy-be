package rest

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/usecase/facility"
	"badbuddy/internal/usecase/user"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FacilityHandler struct {
	facilityUseCase facility.UseCase
	userUseCase     user.UseCase
}

func NewFacilityHandler(facilityUseCase facility.UseCase, userUseCase user.UseCase) *FacilityHandler {
	return &FacilityHandler{
		facilityUseCase: facilityUseCase,
		userUseCase:     userUseCase,
	}
}

func (h *FacilityHandler) SetupFacilityRoutes(app *fiber.App) {
	facility := app.Group("/api/facilities")

	// Public routes

	// Protected routes
	facility.Use(middleware.AuthRequired())
	facility.Get("/", h.ListFacilities)
	facility.Get("/:id", h.GetFacility)
	facility.Post("/", h.CreateFacility)
	facility.Put("/:id", h.UpdateFacility)
	facility.Delete("/:id", h.DeleteFacility)
}

func (h *FacilityHandler) ListFacilities(c *fiber.Ctx) error {
	facilities, err := h.facilityUseCase.ListFacilities(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(facilities)
}

func (h *FacilityHandler) GetFacility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	isAdmin, err := h.userUseCase.IsAdmin(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	if !isAdmin {
		return h.handleError(c, facility.ErrUnauthorized)

	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.handleError(c, errors.New("invalid facility ID format"))
	}

	facility, err := h.facilityUseCase.GetFacilityByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(facility)
}

func (h *FacilityHandler) CreateFacility(c *fiber.Ctx) error {
	var req requests.CreateAndUpdateFacilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Facility name cannot be empty",
			Code:        "INVALID_REQUEST",
			Description: "Facility name cannot be empty",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	isAdmin, err := h.userUseCase.IsAdmin(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	if !isAdmin {
		return h.handleError(c, facility.ErrUnauthorized)
	}

	facility, err := h.facilityUseCase.CreateFacility(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(facility)
}

func (h *FacilityHandler) UpdateFacility(c *fiber.Ctx) error {
	var req requests.CreateAndUpdateFacilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Facility name cannot be empty",
			Code:        "INVALID_REQUEST",
			Description: "Facility name cannot be empty",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	isAdmin, err := h.userUseCase.IsAdmin(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	if !isAdmin {
		return h.handleError(c, facility.ErrUnauthorized)
	}

	id, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return h.handleError(c, errors.New("invalid facility ID format"))
	}

	facility, err := h.facilityUseCase.UpdateFacility(c.Context(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(facility)

}

func (h *FacilityHandler) DeleteFacility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	isAdmin, err := h.userUseCase.IsAdmin(c.Context(), userID)

	if err != nil {
		return h.handleError(c, err)
	}

	if !isAdmin {
		return h.handleError(c, facility.ErrUnauthorized)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return h.handleError(c, errors.New("invalid facility ID format"))
	}

	err = h.facilityUseCase.DeleteFacility(c.Context(), id)

	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Facility deleted successfully",
	})

}

func (h *FacilityHandler) handleError(c *fiber.Ctx, err error) error {
	var status int
	var errorResponse responses.ErrorResponse

	// Add specific error type checks here if needed
	switch {
	case errors.Is(err, facility.ErrFacilityNotFound):
		status = fiber.StatusNotFound
		errorResponse = responses.ErrorResponse{
			Error: "Chat not found",
			Code:  "CHAT_NOT_FOUND",
		}
	case errors.Is(err, facility.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		errorResponse = responses.ErrorResponse{
			Error: "Unauthorized",
			Code:  "UNAUTHORIZED",
		}
	case errors.Is(err, facility.ErrValidation):
		status = fiber.StatusBadRequest
		errorResponse = responses.ErrorResponse{
			Error: "Validation error",
			Code:  "VALIDATION_ERROR",
		}
	default:
		status = fiber.StatusInternalServerError
		errorResponse = responses.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		}
	}

	errorResponse.Description = err.Error()
	return c.Status(status).JSON(errorResponse)
}

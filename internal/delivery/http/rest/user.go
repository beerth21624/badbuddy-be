package rest

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/usecase/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase user.UseCase
}

func NewUserHandler(userUseCase user.UseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}
func (h *UserHandler) SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/api/users")

	userGroup.Post("/register", h.Register)
	userGroup.Post("/login", h.Login)

	// Protected routes
	userGroup.Use(middleware.AuthRequired())
	userGroup.Get("/profile", h.GetProfile)
	userGroup.Put("/profile", h.UpdateProfile)
	userGroup.Get("/search", h.SearchUsers)
	userGroup.Put("/update/role", h.UpdateRoles)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.userUseCase.Register(c.Context(), req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.userUseCase.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, err := uuid.Parse(response.User.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	venues, err := h.userUseCase.GetVenueUserOwn(c.Context(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response.User.Venues = venues

	return c.JSON(response)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	profile, err := h.userUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	venues, err := h.userUseCase.GetVenueUserOwn(c.Context(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	profile.Venues = venues

	return c.JSON(profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req requests.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.userUseCase.UpdateProfile(c.Context(), userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	filters := requests.SearchFilters{
		Limit:  c.QueryInt("limit", 10),
		Offset: c.QueryInt("offset", 0),
	}

	if filters.Limit <= 0 || filters.Limit > 100 {
		filters.Limit = 10
	}

	if filters.Offset < 0 {
		filters.Offset = 0
	}

	users, err := h.userUseCase.SearchUsers(c.Context(), query, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

func (h *UserHandler) UpdateRoles(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req requests.UpdateRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.userUseCase.UpdateRoles(c.Context(), userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Roles updated successfully",
	})
}



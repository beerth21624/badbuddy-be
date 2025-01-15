package rest

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/usecase/facility"
	"badbuddy/internal/usecase/user"
	"badbuddy/internal/usecase/venue"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VenueHandler struct {
	venueUseCase    venue.UseCase
	facilityUseCase facility.UseCase
	userUseCase     user.UseCase
}

func NewVenueHandler(venueUseCase venue.UseCase, facilityUseCase facility.UseCase, userUseCase user.UseCase) *VenueHandler {
	return &VenueHandler{
		venueUseCase:    venueUseCase,
		facilityUseCase: facilityUseCase,
		userUseCase:     userUseCase,
	}
}

func (h *VenueHandler) SetupVenueRoutes(app *fiber.App) {
	venueGroup := app.Group("/api/venues")

	// Public routes
	venueGroup.Get("/", h.ListVenues)
	venueGroup.Get("/search", h.SearchVenues)
	venueGroup.Get("/:id", h.GetVenue)
	venueGroup.Get("/:id/reviews", h.GetReviews)
	venueGroup.Get("/:id/facilities", h.GetFacilitiesOfVenue)

	// Protected routes
	venueGroup.Use(middleware.AuthRequired())
	venueGroup.Post("/", h.CreateVenue)
	//update court
	venueGroup.Put("/:id/courts/:courtId", h.UpdateCourt)
	venueGroup.Put("/:id", h.UpdateVenue)
	venueGroup.Post("/:id/courts", h.AddCourt)
	venueGroup.Post("/:id/reviews", h.AddReview)

	// delete court
	venueGroup.Delete("/:id/courts/:courtId", h.DeleteCourt)
}

func (h *VenueHandler) CreateVenue(c *fiber.Ctx) error {
	var req requests.CreateVenueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	ownerID := c.Locals("userID").(uuid.UUID)

	facility := req.Facilities

	if !h.validateFacilities(facility, c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid facility ID",
		})
	}

	venue, err := h.venueUseCase.CreateVenue(c.Context(), ownerID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(venue)
}

func (h *VenueHandler) GetVenue(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	venue, err := h.venueUseCase.GetVenue(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(venue)
}

// เพิ่ม method UpdateVenue
func (h *VenueHandler) UpdateVenue(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	ownerID := c.Locals("userID").(uuid.UUID)

	isAdmin, err := h.userUseCase.IsAdmin(c.Context(), ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isOwner, err := h.venueUseCase.IsOwner(c.Context(), id, ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// check isAdmin and pass owner check
	if !isAdmin {
		if !isOwner {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
	}

	var req requests.UpdateVenueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	facility := req.Facilities

	if !h.validateFacilities(facility, c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid facility ID",
		})
	}

	if err := h.venueUseCase.UpdateVenue(c.Context(), id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Venue updated successfully",
	})
}

func (h *VenueHandler) ListVenues(c *fiber.Ctx) error {
	location := c.Query("location", "")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	venues, err := h.venueUseCase.ListVenues(c.Context(), location, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"venues": venues,
	})
}

func (h *VenueHandler) SearchVenues(c *fiber.Ctx) error {
	query := c.Query("q")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)
	minPrice := c.QueryInt("minPrice", 0)
	maxPrice := c.QueryInt("maxPrice", 1000)
	location := c.Query("location", "")
	facility := c.Query("facility")
	facilityList := strings.Split(facility, ",")

	if facilityList[0] == "" && len(facilityList) == 1 {
		facilityList = []string{}
	}

	venues, err := h.venueUseCase.SearchVenues(c.Context(), query, limit, offset, minPrice, maxPrice, location, facilityList)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(venues)
}

func (h *VenueHandler) AddCourt(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	// check ownerID is owner or not
	ownerID := c.Locals("userID").(uuid.UUID)
	isOwner, err := h.venueUseCase.IsOwner(c.Context(), venueID, ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !isOwner {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req requests.CreateCourtRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	court, err := h.venueUseCase.AddCourt(c.Context(), venueID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(court)
}

func (h *VenueHandler) UpdateCourt(c *fiber.Ctx) error {
	vendorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	// check ownerID is owner or not
	ownerID := c.Locals("userID").(uuid.UUID)
	isOwner, err := h.venueUseCase.IsOwner(c.Context(), vendorID, ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !isOwner {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	courtID, err := uuid.Parse(c.Params("courtId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid court ID",
		})
	}

	var req requests.UpdateCourtRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.CourtID = courtID.String()

	if err := h.venueUseCase.UpdateCourt(c.Context(), vendorID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Court updated successfully",
	})
}

func (h *VenueHandler) DeleteCourt(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	// check ownerID is owner or not
	ownerID := c.Locals("userID").(uuid.UUID)
	isOwner, err := h.venueUseCase.IsOwner(c.Context(), venueID, ownerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !isOwner {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	courtID, err := uuid.Parse(c.Params("courtId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid court ID",
		})
	}

	if err := h.venueUseCase.DeleteCourt(c.Context(), venueID, courtID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Court deleted successfully",
	})
}

// เพิ่ม method GetReviews
func (h *VenueHandler) GetReviews(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	reviews, err := h.venueUseCase.GetReviews(c.Context(), venueID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"reviews": reviews,
	})
}

func (h *VenueHandler) AddReview(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	var req requests.AddReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.venueUseCase.AddReview(c.Context(), venueID, userID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Review added successfully",
	})
}

func (h *VenueHandler) GetFacilitiesOfVenue(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid venue ID",
		})
	}

	facilities, err := h.venueUseCase.GetFacilities(c.Context(), venueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(facilities)
}

func (h *VenueHandler) validateFacilities(facility []requests.Facility, c *fiber.Ctx) bool {
	for _, f := range facility {
		facilityID, err := uuid.Parse(f.ID)
		if err != nil {
			return false
		}
		_, err = h.facilityUseCase.GetFacilityByID(c.Context(), facilityID)
		if err != nil {
			return false
		}
	}
	return true
}

package rest

import (
	"time"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/usecase/booking"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BookingHandler struct {
	bookingUseCase booking.UseCase
}

func NewBookingHandler(bookingUseCase booking.UseCase) *BookingHandler {
	return &BookingHandler{
		bookingUseCase: bookingUseCase,
	}
}

func (h *BookingHandler) SetupBookingRoutes(app *fiber.App) {
	bookings := app.Group("/api/bookings")

	// Public routes
	bookings.Get("/availability", h.CheckAvailability)

	// Protected routes
	bookings.Use(middleware.AuthRequired())
	bookings.Post("/", h.CreateBooking)
	bookings.Get("/", h.ListBookings)
	bookings.Get("/:id", h.GetBooking)
	bookings.Put("/:id", h.UpdateBooking)
	bookings.Post("/:id/cancel", h.CancelBooking)
	bookings.Get("/user/me", h.GetUserBookings)
	bookings.Get("/:id/payment", h.GetPayment)
	bookings.Post("/:id/payment", h.CreatePayment)
	bookings.Put("/:id/payment", h.UpdatePayment)

	bookings.Post("/test", h.ChangeCourtStatus)
}

// CreateBooking handles the creation of a new booking
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req requests.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	booking, err := h.bookingUseCase.CreateBooking(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Booking created successfully",
		Data:    booking,
	})
}

// GetBooking handles retrieving a single booking
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid booking ID",
			Code:        "INVALID_ID",
			Description: "The provided booking ID is not in a valid format",
		})
	}

	booking, err := h.bookingUseCase.GetBooking(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Data: booking,
	})
}

// ListBookings handles listing bookings with filters
func (h *BookingHandler) ListBookings(c *fiber.Ctx) error {
	var req requests.ListBookingsRequest

	// Parse query parameters
	req.CourtID = c.Query("court_id")
	req.VenueID = c.Query("venue_id")
	req.DateFrom = c.Query("date_from")
	req.DateTo = c.Query("date_to")
	req.Status = c.Query("status")
	req.Limit = c.QueryInt("limit", 10)
	req.Offset = c.QueryInt("offset", 0)

	userID := c.Locals("userID").(uuid.UUID)

	bookings, err := h.bookingUseCase.ListBookings(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(bookings)
}

// UpdateBooking handles updating a booking
func (h *BookingHandler) UpdateBooking(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid booking ID",
			Code:        "INVALID_ID",
			Description: "The provided booking ID is not in a valid format",
		})
	}

	var req requests.UpdateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	booking, err := h.bookingUseCase.UpdateBooking(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Booking updated successfully",
		Data:    booking,
	})
}

// CancelBooking handles cancelling a booking
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid booking ID",
			Code:        "INVALID_ID",
			Description: "The provided booking ID is not in a valid format",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	if err := h.bookingUseCase.CancelBooking(c.Context(), id, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Booking cancelled successfully",
	})
}

// GetUserBookings handles retrieving user's bookings
func (h *BookingHandler) GetUserBookings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	includeHistory := c.QueryBool("include_history", false)

	bookings, err := h.bookingUseCase.GetUserBookings(c.Context(), userID, includeHistory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Data: bookings,
	})
}

// CheckAvailability handles checking court availability
func (h *BookingHandler) CheckAvailability(c *fiber.Ctx) error {
	var req requests.CheckAvailabilityRequest

	req.CourtID = c.Query("court_id")
	req.Date = c.Query("date")
	req.StartTime = c.Query("start_time")
	req.EndTime = c.Query("end_time")

	availability, err := h.bookingUseCase.CheckAvailability(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Data: availability,
	})
}

// get payment for booking
func (h *BookingHandler) GetPayment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid Court Booking ID",
			Code:        "INVALID_ID",
			Description: "The provided court booking ID is not in a valid format",
		})
	}

	booking, err := h.bookingUseCase.GetPayment(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(responses.SuccessResponse{
		Data: booking,
	})
}

// CreatePayment handles creating a payment for a booking
func (h *BookingHandler) CreatePayment(c *fiber.Ctx) error {
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid booking ID",
			Code:        "INVALID_ID",
			Description: "The provided booking ID is not in a valid format",
		})
	}

	var req requests.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}
	userID := c.Locals("userID").(uuid.UUID)
	payment, err := h.bookingUseCase.CreatePayment(c.Context(), bookingID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Payment created successfully",
		Data:    payment,
	})
}

func (h *BookingHandler) UpdatePayment(c *fiber.Ctx) error {
	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid booking ID",
			Code:        "INVALID_ID",
			Description: "The provided booking ID is not in a valid format",
		})
	}

	var req requests.UpdatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}
	userID := c.Locals("userID").(uuid.UUID)
	payment, err := h.bookingUseCase.UpdatePayment(c.Context(), bookingID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Payment created successfully",
		Data:    payment,
	})
}

func (h *BookingHandler) ChangeCourtStatus(c *fiber.Ctx) error {

	err := h.bookingUseCase.ChangeCourtStatus(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Court status changed successfully",
	})
}

// handleError centralizes error handling
func (h *BookingHandler) handleError(c *fiber.Ctx, err error) error {
	// Add specific error types
	switch {
	case err == booking.ErrBookingNotFound:
		return c.Status(fiber.StatusNotFound).JSON(responses.ErrorResponse{
			Error: "Booking not found",
			Code:  "BOOKING_NOT_FOUND",
		})
	case err == booking.ErrUnauthorized:
		return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
			Error: "Unauthorized",
			Code:  "UNAUTHORIZED",
		})
	case err == booking.ErrValidation:
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error: "Validation error",
			Code:  "VALIDATION_ERROR",
		})
	case err == booking.ErrBookingConflict:
		return c.Status(fiber.StatusConflict).JSON(responses.ErrorResponse{
			Error: "Booking conflict",
			Code:  "BOOKING_CONFLICT",
		})
	case err == booking.ErrPaymentRequired:
		return c.Status(fiber.StatusPaymentRequired).JSON(responses.ErrorResponse{
			Error: "Payment required",
			Code:  "PAYMENT_REQUIRED",
		})
	default:
		// Log the error here
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}
}

// Additional helper methods for validation and response formatting

// validateTime validates time string format
func (h *BookingHandler) validateTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	return err
}

// validateDate validates date string format
func (h *BookingHandler) validateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	return err
}

// validateUUID validates UUID string format
func (h *BookingHandler) validateUUID(id string) error {
	_, err := uuid.Parse(id)
	return err
}

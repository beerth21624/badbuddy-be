package rest

import (
	"errors"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/usecase/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SessionHandler struct {
	sessionUseCase session.UseCase
}

func NewSessionHandler(sessionUseCase session.UseCase) *SessionHandler {
	return &SessionHandler{
		sessionUseCase: sessionUseCase,
	}
}

func (h *SessionHandler) SetupSessionRoutes(app *fiber.App) {
	sessions := app.Group("/api/sessions")

	// Public routes
	sessions.Get("/", h.ListSessions)
	sessions.Get("/search", h.SearchSessions)
	sessions.Get("/:id", h.GetSession)

	// Protected routes
	sessions.Use(middleware.AuthRequired())
	sessions.Get("/:id/status", h.GetSessionStatus)
	sessions.Get("/join/me", h.GetMyJoinedSessions)
	sessions.Get("/host/me", h.GetMyHostedSessions)
	sessions.Post("/", h.CreateSession)
	sessions.Put("/:id", h.UpdateSession)
	sessions.Post("/:id/join", h.JoinSession)
	sessions.Post("/:id/leave", h.LeaveSession)
	sessions.Post("/:id/cancel", h.CancelSession)
	sessions.Get("/user/me", h.GetUserSessions)
	sessions.Put("/:id/status", h.ChangeParticipantStatus)
	sessions.Get("/:id/participants", h.GetSessionParticipants)
}

func (h *SessionHandler) CreateSession(c *fiber.Ctx) error {
	var req requests.CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	hostID := c.Locals("userID").(uuid.UUID)

	session, err := h.sessionUseCase.CreateSession(c.Context(), hostID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Session created successfully",
		Data:    session,
	})
}

func (h *SessionHandler) GetSession(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}
	session, err := h.sessionUseCase.GetSession(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: session,
	})
}

func (h *SessionHandler) GetSessionStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	session, err := h.sessionUseCase.GetSessionStatus(c.Context(), id, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: session,
	})
}
func (h *SessionHandler) ListSessions(c *fiber.Ctx) error {
	// Parse and validate filters
	filters := make(map[string]interface{})

	if date := c.Query("date"); date != "" {
		filters["date"] = date
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if playerLevel := c.Query("player_level"); playerLevel != "" {
		filters["player_level"] = playerLevel
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Parse pagination params with defaults
	limit := c.QueryInt("limit", 10)

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	sessions, err := h.sessionUseCase.ListSessions(c.Context(), filters, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sessions)
}

func (h *SessionHandler) SearchSessions(c *fiber.Ctx) error {
	query := c.Query("q")
	filters := make(map[string]interface{})
	if date := c.Query("date"); date != "" {
		filters["date"] = date
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if playerLevel := c.Query("player_level"); playerLevel != "" {
		filters["player_level"] = playerLevel
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	limit := c.QueryInt("limit", 10)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	sessions, err := h.sessionUseCase.SearchSessions(c.Context(), query, filters, limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(sessions)
}

func (h *SessionHandler) UpdateSession(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	var req requests.UpdateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	hostID := c.Locals("userID").(uuid.UUID)

	if err := h.sessionUseCase.UpdateSession(c.Context(), sessionID, hostID, req); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Session updated successfully",
	})
}

func (h *SessionHandler) JoinSession(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	var req requests.JoinSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	if err := h.sessionUseCase.JoinSession(c.Context(), sessionID, userID, req); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Successfully joined session",
	})
}

func (h *SessionHandler) LeaveSession(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)

	if err := h.sessionUseCase.LeaveSession(c.Context(), sessionID, userID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Successfully left session",
	})
}

func (h *SessionHandler) CancelSession(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	hostID := c.Locals("userID").(uuid.UUID)

	if err := h.sessionUseCase.CancelSession(c.Context(), sessionID, hostID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Message: "Session cancelled successfully",
	})
}

func (h *SessionHandler) GetUserSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	includeHistory := c.QueryBool("include_history", false)

	sessions, err := h.sessionUseCase.GetUserSessions(c.Context(), userID, includeHistory)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: sessions,
	})
}

func (h *SessionHandler) ChangeParticipantStatus(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	var req requests.ChangeParticipantStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid request body",
			Code:        "INVALID_REQUEST",
			Description: err.Error(),
		})
	}

	hostID := c.Locals("userID").(uuid.UUID)

	if err := h.sessionUseCase.ChangeParticipantStatus(c.Context(), sessionID, hostID, req); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(responses.SuccessResponse{
		Message: "Participant status updated successfully",
	})
}

func (h *SessionHandler) GetSessionParticipants(c *fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Error:       "Invalid session ID",
			Code:        "INVALID_ID",
			Description: "The provided session ID is not in a valid format",
		})
	}

	participants, err := h.sessionUseCase.GetSessionParticipants(c.Context(), sessionID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: participants,
	})
}

func (h *SessionHandler) GetMyJoinedSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	includeHistory := c.QueryBool("include_history", false)

	sessions, err := h.sessionUseCase.GetMyJoinedSessions(c.Context(), userID, includeHistory)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: sessions,
	})
}

func (h *SessionHandler) GetMyHostedSessions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	includeHistory := c.QueryBool("include_history", false)

	sessions, err := h.sessionUseCase.GetMyHostedSessions(c.Context(), userID, includeHistory)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.SuccessResponse{
		Data: sessions,
	})
}

func (h *SessionHandler) handleError(c *fiber.Ctx, err error) error {
	var status int
	var errorResponse responses.ErrorResponse

	switch {
	case errors.Is(err, session.ErrSessionNotFound):
		status = fiber.StatusNotFound
		errorResponse = responses.ErrorResponse{
			Error: "Session not found",
			Code:  "SESSION_NOT_FOUND",
		}
	case errors.Is(err, session.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		errorResponse = responses.ErrorResponse{
			Error: "Unauthorized",
			Code:  "UNAUTHORIZED",
		}
	case errors.Is(err, session.ErrValidation):
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

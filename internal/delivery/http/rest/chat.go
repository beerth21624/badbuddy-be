package rest

import (
	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/delivery/http/middleware"
	"badbuddy/internal/delivery/http/ws"
	"badbuddy/internal/usecase/chat"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"errors"
)

type ChatHandler struct {
	chatUseCase chat.UseCase
	chatHub     *ws.ChatHub
}

func NewChatHandler(chatUseCase chat.UseCase, chatHub *ws.ChatHub) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUseCase,
		chatHub:     chatHub,
	}
}

func (h *ChatHandler) SetupChatRoutes(app *fiber.App) {
	chat := app.Group("/api/chats")

	// Public routes

	// Protected routes
	chat.Use(middleware.AuthRequired())
	chat.Get("/", h.GetChats)
	chat.Get("/:chatID/messages", h.GetChatMessage)
	chat.Post("/:chatID/messages", h.SendMessage)
	chat.Delete("/:chatID/messages/:messageID", h.DeleteMessage)
	chat.Put("/:chatID/messages/:messageID", h.UpdateMessage)

	chat.Get("/:chatID/users", h.GetUsersInChat)

	chat.Get("direct/:userID/messages", h.GetDirectChat)
	chat.Get("session/:sessionID/messages", h.GetChatMessageOfSession)
}

func (h *ChatHandler) GetChatMessage(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid limit format"))
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid offset format"))
	}

	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return h.handleError(c, errors.New("invalid chat ID format"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	chat, err := h.chatUseCase.GetChatMessageByID(c.Context(), chatUUID, limit, offset, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	message_bytes, _ := json.Marshal(responses.BoardCastMessageResponse{
		MessageaType: "read_all_message",
		Data:         map[string]interface{}{"user_id": userID},
	})
	h.chatHub.GetRoom(chatUUID.String()).Broadcast <- message_bytes

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Chat messages retrieved successfully",
		Data:    chat,
	})
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var req requests.SendAndUpdateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.New("invalid request body"))
	}

	if req.Message == "" {
		return h.handleError(c, errors.New("message cannot be empty"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	chatID := c.Params("chatID")
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return h.handleError(c, errors.New("invalid chat ID format"))
	}

	chatMessage, err := h.chatUseCase.SendMessage(c.Context(), userID, chatUUID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	messageBytes, _ := json.Marshal(responses.BoardCastMessageResponse{
		MessageaType: "send_message",
		Data:         chatMessage,
	})
	h.chatHub.GetRoom(chatUUID.String()).Broadcast <- messageBytes

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Message sent successfully",
		Data:    chatMessage,
	})
}

func (h *ChatHandler) handleError(c *fiber.Ctx, err error) error {
	var status int
	var errorResponse responses.ErrorResponse

	// Add specific error type checks here if needed
	switch {
	case errors.Is(err, chat.ErrChatNotFound):
		status = fiber.StatusNotFound
		errorResponse = responses.ErrorResponse{
			Error: "Chat not found",
			Code:  "CHAT_NOT_FOUND",
		}
	case errors.Is(err, chat.ErrUnauthorized):
		status = fiber.StatusUnauthorized
		errorResponse = responses.ErrorResponse{
			Error: "Unauthorized",
			Code:  "UNAUTHORIZED",
		}
	case errors.Is(err, chat.ErrValidation):
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

func (h *ChatHandler) DeleteMessage(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	messageID := c.Params("messageID")

	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return h.handleError(c, errors.New("invalid chat ID format"))
	}

	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		return h.handleError(c, errors.New("invalid message ID format"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	err = h.chatUseCase.DeleteMessage(c.Context(), chatUUID, messageUUID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	messageBytes, _ := json.Marshal(responses.BoardCastMessageResponse{
		MessageaType: "delete_message",
		Data:         map[string]interface{}{"message_id": messageID},
	})
	h.chatHub.GetRoom(chatUUID.String()).Broadcast <- messageBytes

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Message deleted successfully",
	})
}

func (h *ChatHandler) UpdateMessage(c *fiber.Ctx) error {
	var req requests.SendAndUpdateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.New("invalid request body"))
	}

	if req.Message == "" {
		return h.handleError(c, errors.New("message cannot be empty"))
	}

	chatID := c.Params("chatID")
	messageID := c.Params("messageID")

	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return h.handleError(c, errors.New("invalid chat ID format"))
	}

	messageUUID, err := uuid.Parse(messageID)
	if err != nil {
		return h.handleError(c, errors.New("invalid message ID format"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	err = h.chatUseCase.UpdateMessage(c.Context(), chatUUID, messageUUID, userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	messageBytes, _ := json.Marshal(responses.BoardCastMessageResponse{
		MessageaType: "update_message",
		Data: map[string]interface{}{
			"message_id": messageID,
			"message":    req.Message,
		},
	})
	h.chatHub.GetRoom(chatUUID.String()).Broadcast <- messageBytes

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Message updated successfully",
	})
}

func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	chats, err := h.chatUseCase.GetChats(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Chats retrieved successfully",
		Data:    chats,
	})
}

func (h *ChatHandler) GetUsersInChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return h.handleError(c, errors.New("invalid chat ID format"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	users, err := h.chatUseCase.GetUsersInChat(c.Context(), chatUUID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Chat users retrieved successfully",
		Data:    users,
	})
}

func (h *ChatHandler) GetDirectChat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	otherUserID := c.Params("userID")
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid limit format"))
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid offset format"))
	}

	otherUserUUID, err := uuid.Parse(otherUserID)
	if err != nil {
		return h.handleError(c, errors.New("invalid user ID format"))
	}

	chat, err := h.chatUseCase.GetDirectChat(c.Context(), userID, otherUserUUID, limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Direct chat retrieved successfully",
		Data:    chat,
	})
}

func (h *ChatHandler) GetChatMessageOfSession(c *fiber.Ctx) error {
	sessionID := c.Params("sessionID")
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid limit format"))
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return h.handleError(c, errors.New("invalid offset format"))
	}

	userID := c.Locals("userID").(uuid.UUID)

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return h.handleError(c, errors.New("invalid session ID format"))
	}

	chat, err := h.chatUseCase.GetChatMessageOfSession(c.Context(), sessionUUID, limit, offset, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{
		Message: "Chat messages retrieved successfully",
		Data:    chat,
	})
}
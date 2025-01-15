// middleware/auth.go
package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrNoAuthHeader  = errors.New("authorization header required")
	ErrInvalidFormat = errors.New("invalid token format")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid token claims")
	ErrInvalidUserID = errors.New("invalid user ID in token")
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrNoAuthHeader.Error(),
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrInvalidFormat.Error(),
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte("your-jwt-secret"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrInvalidToken.Error(),
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrInvalidClaims.Error(),
			})
		}

		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrInvalidUserID.Error(),
			})
		}

		// Set user ID in context for later use
		c.Locals("userID", userID)

		return c.Next()
	}
}

// GetUserID gets the user ID from the Fiber context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}

// MustGetUserID gets the user ID from context or panics
// Use this only when you're sure the auth middleware has run
func MustGetUserID(c *fiber.Ctx) uuid.UUID {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		panic("user ID not found in context")
	}
	return userID
}

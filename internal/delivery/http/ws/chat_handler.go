// internal/delivery/http/ws/chat_handler.go
package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func ChatWebSocketHandler(hub *ChatHub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		chatID := c.Params("chat_id")
		room := hub.GetRoom(chatID)

		room.Clients[c] = true
		defer func() {
			delete(room.Clients, c)
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	})
}

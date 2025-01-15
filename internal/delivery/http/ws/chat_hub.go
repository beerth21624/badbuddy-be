// internal/delivery/http/ws/chat_hub.go
package ws

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type ChatRoom struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
}

type ChatHub struct {
	Rooms map[string]*ChatRoom
	mu    sync.Mutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Rooms: make(map[string]*ChatRoom),
	}
}

func (h *ChatHub) GetRoom(chatID string) *ChatRoom {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Rooms[chatID]; !ok {
		h.Rooms[chatID] = &ChatRoom{
			Clients:   make(map[*websocket.Conn]bool),
			Broadcast: make(chan []byte),
		}
		go h.runRoom(h.Rooms[chatID])
	}
	return h.Rooms[chatID]
}

func (h *ChatHub) runRoom(room *ChatRoom) {
	for {
		msg := <-room.Broadcast
		for client := range room.Clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Close()
				delete(room.Clients, client)
			}
		}
	}
}

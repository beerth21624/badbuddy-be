// internal/domain/models/chat.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string
type MessageType string
type MessageStatus string

const (
	ChatTypeDirect  ChatType = "direct"
	ChatTypeGroup   ChatType = "group"
	ChatTypeSession ChatType = "session"

	MessageTypeText   MessageType = "text"
	MessageTypeImage  MessageType = "image"
	MessageTypeSystem MessageType = "system"

	MessageStatusSent MessageStatus = "sent"
	MessageStatusRead MessageStatus = "read"
)

// Chat represents a conversation between users
type Chat struct {
	ID   uuid.UUID `db:"id"`
	Type ChatType  `db:"type"`
	SessionID *uuid.UUID `db:"session_id"`
	LastMessage *Message `db:"last_message,omitempty"`
	Users []User `db:"users,omitempty"`
}

// ChatParticipant represents a user in a chat
type ChatParticipant struct {
	ID         uuid.UUID `db:"id"`
	ChatID     uuid.UUID `db:"chat_id"`
	UserID     uuid.UUID `db:"user_id"`
	IsAdmin    bool      `db:"is_admin"`
	LastReadAt time.Time `db:"last_read_at"`
	JoinedAt   time.Time `db:"joined_at"`
	LeftAt     time.Time `db:"left_at"`

	// Populated fields
	User *User `db:"user,omitempty"`
}

// Message represents a single message in a chat
type Message struct {
	ID           uuid.UUID     `db:"m_id"`
	ChatID       uuid.UUID     `db:"chat_id"`
	SenderID     uuid.UUID     `db:"sender_id"`
	Type         MessageType   `db:"type"`
	Content      string        `db:"content"`
	Status       MessageStatus `db:"status"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
	DeletedAt    *time.Time    `db:"delete_at"`
	UserID       uuid.UUID     `db:"u_id"`
	Email        string        `db:"email"`
	FirstName    string        `db:"first_name"`
	LastName     string        `db:"last_name"`
	Phone        string        `db:"phone"`
	PlayLevel    string        `db:"play_level"`
	AvatarURL    *string        `db:"avatar_url"`
	Gender       *string       `db:"gender"`
	Location     *string       `db:"location"`
	Bio          *string        `db:"bio"`
	LastActiveAt time.Time     `db:"last_active_at"`

	// Populated fields
	// Sender *User       `db:"sender,omitempty"`
	// ReadBy []uuid.UUID `db:"read_by,omitempty"`
}

// MessageReceipt tracks message delivery and read status
type MessageReceipt struct {
	ID        uuid.UUID     `db:"id"`
	MessageID uuid.UUID     `db:"message_id"`
	UserID    uuid.UUID     `db:"user_id"`
	Status    MessageStatus `db:"status"`
	ReadAt    time.Time     `db:"read_at"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

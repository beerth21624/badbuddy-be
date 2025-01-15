package postgres

import (
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chatRepository struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) interfaces.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) GetChatMessageByID(ctx context.Context, chatID uuid.UUID, limit int, offset int) (*[]models.Message, error) {
	// Get chat
	chat := models.Chat{}

	query := `SELECT * FROM chats WHERE id = $1`

	err := r.db.GetContext(ctx, &chat, query, chatID)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT 
			m.id AS m_id,
			m.chat_id,
			m.sender_id,
			m.type,
			m.content,
			m.created_at,
			m.updated_at,
			u.email,
			u.first_name,
			u.last_name,
			u.phone,
			u.play_level,
			u.avatar_url,
			u.play_level,
			u.gender,
			u.location,
			u.bio,
			u.last_active_at
		FROM 
			chat_messages m
		JOIN 
			users u ON m.sender_id = u.id
		WHERE 
			m.chat_id = $1
			AND m.delete_at IS NULL
		ORDER BY 
			m.created_at ASC
		LIMIT $2
		OFFSET $3`

	// Get messages
	messages := []models.Message{}
	err = r.db.SelectContext(ctx, &messages, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &messages, nil
}

func (r *chatRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error) {
	chat := models.Chat{}

	query := `SELECT * FROM chats WHERE id = $1`

	err := r.db.GetContext(ctx, &chat, query, chatID)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *chatRepository) IsUserPartOfChat(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM chat_participants WHERE user_id = $1 AND chat_id = $2`

	err := r.db.GetContext(ctx, &count, query, userID, chatID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error) {

	query := `INSERT INTO chat_messages (id, chat_id, sender_id, type, content, created_at, updated_at, status) VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6)`

	_, err := r.db.ExecContext(ctx, query, message.ID, message.ChatID, message.SenderID, message.Type, message.Content, message.Status)
	if err != nil {
		return nil, err
	}

	messageReturn, err := r.GetMessageByID(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	return messageReturn, nil

}

func (r *chatRepository) GetMessageByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	message := models.Message{}

	query := `
		SELECT 
			m.id AS m_id,
			m.chat_id,
			m.sender_id,
			m.type,
			m.content,
			m.created_at,
			m.updated_at,
			u.email,
			u.first_name,
			u.last_name,
			u.phone,
			u.play_level,
			u.avatar_url,
			u.play_level,
			u.gender,
			u.location,
			u.bio,
			u.last_active_at
		FROM 
			chat_messages m
		JOIN 
			users u ON m.sender_id = u.id
		WHERE 
			m.id = $1`

	err := r.db.GetContext(ctx, &message, query, messageID)
	if err != nil {
		return nil, err
	}

	return &message, nil

}

func (r *chatRepository) CreateChat(ctx context.Context, chat *models.Chat) error {

	query := `INSERT INTO chats (id, type, session_id) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, chat.ID, chat.Type, chat.SessionID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) AddUserToChat(ctx context.Context, userID, chatID uuid.UUID) error {

	query := `INSERT INTO chat_participants (id, chat_id, user_id) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, uuid.New(), chatID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) RemoveUserFromChat(ctx context.Context, userID, chatID uuid.UUID) error {

	query := `DELETE FROM chat_participants WHERE chat_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) UpdateChatMessage(ctx context.Context, message *models.Message) error {

	query := `UPDATE chat_messages SET content = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, message.Content, message.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) DeleteChatMessage(ctx context.Context, messageID uuid.UUID) error {

	query := `UPDATE chat_messages SET delete_at = NOW(), updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) UpdateChatMessageReadStatus(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {

	query := `UPDATE chat_messages SET status = 'read' WHERE chat_id = $1 AND sender_id != $2 AND status = 'sent'`

	_, err := r.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) IsUserIsSender(ctx context.Context, userID, messageID uuid.UUID) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM chat_messages WHERE sender_id = $1 AND id = $2`

	err := r.db.GetContext(ctx, &count, query, userID, messageID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *chatRepository) GetChats(ctx context.Context, userID uuid.UUID) (*[]models.Chat, error) {
	chats := []models.Chat{}

	query := `
		SELECT 
			id,
			type,
			session_id
		FROM
			chats
		WHERE
			id IN (SELECT chat_id FROM chat_participants WHERE user_id = $1)`

	err := r.db.SelectContext(ctx, &chats, query, userID)
	if err != nil {
		return nil, err
	}

	lastMessages := []models.Message{}
	for i, chat := range chats {
		query = `
			SELECT
				m.id AS m_id,
				m.chat_id,
				m.sender_id,
				m.type,
				m.content,
				m.created_at,
				m.updated_at,
				u.email,
				u.first_name,
				u.last_name,
				u.phone,
				u.play_level,
				u.avatar_url,
				u.play_level,
				u.avatar_url,
				u.gender,
				u.location,
				u.bio,
				u.last_active_at
			FROM
				chat_messages m
			JOIN
				users u ON m.sender_id = u.id
			WHERE
				m.chat_id = $1
			ORDER BY
				m.created_at DESC
			LIMIT 1`

		err = r.db.SelectContext(ctx, &lastMessages, query, chat.ID)
		if err != nil {
			return nil, err
		}

		if len(lastMessages) > 0 {
			chats[i].LastMessage = &lastMessages[0]
		} else {
			chats[i].LastMessage = nil
		}

	}

	chatUsers := []models.User{}
	for i, chat := range chats {
		query = `
			SELECT
				u.id,
				u.email,
				u.first_name,
				u.last_name,
				u.phone,
				u.play_level,
				u.location,
				u.bio,
				u.avatar_url,
				u.last_active_at,
				u.gender,
				u.play_hand,
				u.avatar_url
			FROM
				chat_participants cp
			JOIN
				users u ON cp.user_id = u.id
			WHERE
				cp.chat_id = $1`

		err = r.db.SelectContext(ctx, &chatUsers, query, chat.ID)
		if err != nil {
			return nil, err
		}

		chats[i].Users = chatUsers
	}

	return &chats, nil
}

func (r *chatRepository) GetUsersInChat(ctx context.Context, chatID uuid.UUID) (*[]models.User, error) {
	users := []models.User{}

	query := `
		SELECT
			u.id,
			u.email,
			u.first_name,
			u.last_name,
			u.phone,
			u.play_level,
			u.location,
			u.bio,
			u.avatar_url,
			u.last_active_at
		FROM
			chat_participants cp
		JOIN
			users u ON cp.user_id = u.id
		WHERE
			cp.chat_id = $1`

	err := r.db.SelectContext(ctx, &users, query, chatID)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *chatRepository) GetDirectChatID(ctx context.Context, userID, otherUserUUID uuid.UUID) (uuid.UUID, error) {
	var chatID uuid.UUID

	query := `
		SELECT 
			chat_id
		FROM 
			chat_participants
		WHERE 
			user_id = $1
			AND chat_id IN (SELECT chat_id FROM chat_participants WHERE user_id = $2)`

	err := r.db.GetContext(ctx, &chatID, query, userID, otherUserUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			if chatID == uuid.Nil {
				chatID = uuid.New()
				query = `INSERT INTO chats (id, type) VALUES ($1, 'direct')`
				_, err = r.db.ExecContext(ctx, query, chatID)
				if err != nil {
					return uuid.Nil, err
				}

				query = `INSERT INTO chat_participants (id, chat_id, user_id) VALUES ($1, $2, $3), ($4, $2, $5)`
				_, err = r.db.ExecContext(ctx, query, uuid.New(), chatID, userID, uuid.New(), otherUserUUID)
				if err != nil {
					return uuid.Nil, err
				}
			}
		} else {
			return uuid.Nil, err
		}
	}

	return chatID, nil
}

func (r *chatRepository) GetChatIDBySessionID(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error) {
	var chatID uuid.UUID

	query := `SELECT id FROM chats WHERE session_id = $1`

	err := r.db.GetContext(ctx, &chatID, query, sessionID)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func (r *chatRepository) IsUserPartOfSession(ctx context.Context, userID, sessionID uuid.UUID) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM chat_participants WHERE user_id = $1 AND chat_id IN (SELECT id FROM chats WHERE session_id = $2)`

	// query := `SELECT COUNT(*) FROM session_participants WHERE user_id = $1 AND session_id = $2`

	err := r.db.GetContext(ctx, &count, query, userID, sessionID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

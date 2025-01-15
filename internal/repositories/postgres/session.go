package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) interfaces.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO play_sessions (
			id, host_id, venue_id, title, description,
			session_date, start_time, end_time, player_level,
			max_participants, cost_per_person, allow_cancellation,
			cancellation_deadline_hours, is_public, status,
			created_at, updated_at
		) VALUES (
			:id, :host_id, :venue_id, :title, :description,
			:session_date, :start_time, :end_time, :player_level,
			:max_participants, :cost_per_person, :allow_cancellation,
			:cancellation_deadline_hours, :is_public, :status,
			:created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SessionDetail, error) {
	query := `
		SELECT 
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp.id) FILTER (WHERE sp.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		WHERE ps.id = $1
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level`

	session := &models.SessionDetail{}
	err := r.db.GetContext(ctx, session, query, id)
	if err != nil {
		return nil, err
	}

	// Get participants
	participantsQuery := `
		SELECT sp.*, u.first_name || ' ' || u.last_name as user_name
		FROM session_participants sp
		JOIN users u ON u.id = sp.user_id
		WHERE sp.session_id = $1
		ORDER BY sp.joined_at`

	err = r.db.SelectContext(ctx, &session.Participants, participantsQuery, id)
	if err != nil {
		return nil, err
	}

	// Get session rules
	rulesQuery := `
		SELECT *
		FROM session_rules
		WHERE session_id = $1`

	err = r.db.SelectContext(ctx, &session.Rules, rulesQuery, id)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE play_sessions SET
			title = :title,
			description = :description,
			session_date = :session_date,
			start_time = :start_time,
			end_time = :end_time,
			player_level = :player_level,
			max_participants = :max_participants,
			cost_per_person = :cost_per_person,
			allow_cancellation = :allow_cancellation,
			cancellation_deadline_hours = :cancellation_deadline_hours,
			is_public = :is_public,
			status = :status,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *sessionRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.SessionDetail, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filters {
		switch key {
		case "date":
			conditions = append(conditions, fmt.Sprintf("ps.session_date = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "location":
			conditions = append(conditions, fmt.Sprintf("v.location = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "player_level":
			conditions = append(conditions, fmt.Sprintf("ps.player_level = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "status":
			conditions = append(conditions, fmt.Sprintf("ps.status = $%d", argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT 
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp.id) FILTER (WHERE sp.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		WHERE %s
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level
		ORDER BY ps.session_date ASC, ps.start_time ASC
		LIMIT $%d OFFSET $%d`,
		strings.Join(conditions, " AND "),
		argIndex,
		argIndex+1,
	)

	var sessions []models.SessionDetail
	err := r.db.SelectContext(ctx, &sessions, query, args...)
	return sessions, err
}
func (r *sessionRepository) Search(ctx context.Context, searchQuery string, filters map[string]interface{}, limit, offset int) ([]models.SessionDetail, error) {
	conditions := []string{}
	args := []interface{}{searchQuery} // First argument ($1) is always the search query
	argIndex := 2                      // Start from $2 for filter conditions

	conditions = append(conditions, `(
		ps.search_vector @@ plainto_tsquery('english', $1)
		OR v.name ILIKE '%' || $1 || '%'
		OR v.location ILIKE '%' || $1 || '%'
		OR u.first_name ILIKE '%' || $1 || '%'
		OR u.last_name ILIKE '%' || $1 || '%'
	)`)

	// Add filter conditions
	for key, value := range filters {
		switch key {
		case "date":
			conditions = append(conditions, fmt.Sprintf("ps.session_date = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "location":
			conditions = append(conditions, fmt.Sprintf("v.location = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "player_level":
			conditions = append(conditions, fmt.Sprintf("ps.player_level = $%d", argIndex))
			args = append(args, value)
			argIndex++
		case "status":
			conditions = append(conditions, fmt.Sprintf("ps.status = $%d", argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	// Add limit and offset to args
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT 
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp.id) FILTER (WHERE sp.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		WHERE %s
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level
		ORDER BY 
			CASE 
				WHEN ps.search_vector @@ plainto_tsquery('english', $1) 
				THEN ts_rank(ps.search_vector, plainto_tsquery('english', $1))
				ELSE 0
			END DESC,
			ps.session_date ASC,
			ps.start_time ASC
		LIMIT $%d OFFSET $%d`,
		strings.Join(conditions, " AND "),
		argIndex,
		argIndex+1,
	)

	sessions := []models.SessionDetail{}
	err := r.db.SelectContext(ctx, &sessions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search sessions: %w", err)
	}

	return sessions, nil
}

func (r *sessionRepository) AddParticipant(ctx context.Context, participant *models.SessionParticipant) error {
	query := `
		INSERT INTO session_participants (
			id, session_id, user_id, status, joined_at
		) VALUES (
			:id, :session_id, :user_id, :status, :joined_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, participant)
	return err
}

func (r *sessionRepository) UpdateParticipantStatus(ctx context.Context, sessionID, userID uuid.UUID, status models.ParticipantStatus) error {

	query := `
		UPDATE session_participants SET
			status = :status,
			joined_at = CASE WHEN :status = 'confirmed' THEN NOW() ELSE joined_at END,
			cancelled_at = CASE WHEN :status = 'cancelled' THEN NOW() ELSE cancelled_at END
		WHERE session_id = :session_id AND user_id = :user_id`

	result, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
		"status":     status,
		"session_id": sessionID,
		"user_id":    userID,
	})
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("participant not found")
	}

	return nil

}

func (r *sessionRepository) GetParticipants(ctx context.Context, sessionID uuid.UUID) ([]models.SessionParticipant, error) {
	query := `
		SELECT sp.*, u.first_name || ' ' || u.last_name as user_name,u.play_level as player_level
		FROM session_participants sp
		JOIN users u ON u.id = sp.user_id
		WHERE sp.session_id = $1
		ORDER BY sp.joined_at`

	var participants []models.SessionParticipant
	err := r.db.SelectContext(ctx, &participants, query, sessionID)
	return participants, err
}

func (r *sessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error) {
	conditions := []string{
		"(ps.host_id = $1 OR sp.user_id = $1)",
	}

	if !includeHistory {
		conditions = append(conditions, "ps.session_date >= CURRENT_DATE")
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp2.id) FILTER (WHERE sp2.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		LEFT JOIN session_participants sp2 ON sp2.session_id = ps.id
		WHERE %s
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level
		ORDER BY ps.session_date DESC, ps.start_time DESC`,
		strings.Join(conditions, " AND "),
	)

	var sessions []models.SessionDetail
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	return sessions, err
}
func (r *sessionRepository) GetMyJoinedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error) {
	conditions := []string{
		"(sp.user_id = $1)",
	}

	if !includeHistory {
		conditions = append(conditions, "ps.session_date >= CURRENT_DATE")
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp2.id) FILTER (WHERE sp2.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		LEFT JOIN session_participants sp2 ON sp2.session_id = ps.id
		WHERE %s
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level
		ORDER BY ps.session_date DESC, ps.start_time DESC`,
		strings.Join(conditions, " AND "),
	)

	var sessions []models.SessionDetail
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	return sessions, err
}

func (r *sessionRepository) GetMyHostedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error) {
	conditions := []string{
		"(ps.host_id = $1)",
	}

	if !includeHistory {
		conditions = append(conditions, "ps.session_date >= CURRENT_DATE")
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT
			ps.*,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as host_name,
			u.play_level as host_level,
			COUNT(sp.id) FILTER (WHERE sp.status = 'confirmed') as confirmed_players
		FROM play_sessions ps
		JOIN venues v ON v.id = ps.venue_id
		JOIN users u ON u.id = ps.host_id
		LEFT JOIN session_participants sp ON sp.session_id = ps.id
		WHERE %s
		GROUP BY ps.id, v.name, v.location, u.first_name, u.last_name, u.play_level
		ORDER BY ps.session_date DESC, ps.start_time DESC`,
		strings.Join(conditions, " AND "),
	)

	var sessions []models.SessionDetail
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	return sessions, err
}
func (r *sessionRepository) GetJoinStatus(ctx context.Context, userID, venueID uuid.UUID) (models.JoinStatus, error) {
	query := `
		SELECT 
			CASE 
				WHEN ps.host_id = $1 THEN 'host'::text
				WHEN sp.id IS NOT NULL THEN sp.status::text
				ELSE NULL
			END as status
		FROM play_sessions ps
		LEFT JOIN session_participants sp ON sp.session_id = ps.id AND sp.user_id = $1
		WHERE ps.id = $2`

	var status sql.NullString
	err := r.db.GetContext(ctx, &status, query, userID, venueID)
	if err != nil {
		return "", fmt.Errorf("failed to query join status: %w", err)
	}

	if !status.Valid || status.String == "" {
		return "not_joined", nil
	}

	return models.JoinStatus(status.String), nil
}

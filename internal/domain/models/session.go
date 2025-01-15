package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string
type ParticipantStatus string

const (
	SessionStatusOpen      SessionStatus = "open"
	SessionStatusFull      SessionStatus = "full"
	SessionStatusCancelled SessionStatus = "cancelled"
	SessionStatusCompleted SessionStatus = "completed"

	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
	ParticipantStatusPending   ParticipantStatus = "pending"
	ParticipantStatusCancelled ParticipantStatus = "cancelled"
)

// Session represents a play session
type Session struct {
	ID                        uuid.UUID     `db:"id"`
	HostID                    uuid.UUID     `db:"host_id"`
	VenueID                   uuid.UUID     `db:"venue_id"`
	Title                     string        `db:"title"`
	Description               *string       `db:"description"`
	SessionDate               time.Time     `db:"session_date"`
	StartTime                 time.Time     `db:"start_time"`
	EndTime                   time.Time     `db:"end_time"`
	PlayerLevel               PlayerLevel   `db:"player_level"`
	MaxParticipants           int           `db:"max_participants"`
	CostPerPerson             float64       `db:"cost_per_person"`
	AllowCancellation         bool          `db:"allow_cancellation"`
	CancellationDeadlineHours *int          `db:"cancellation_deadline_hours"`
	IsPublic                  bool          `db:"is_public"`
	Status                    SessionStatus `db:"status"`
	CreatedAt                 time.Time     `db:"created_at"`
	UpdatedAt                 time.Time     `db:"updated_at"`
}

// SessionRule represents a rule for a session
type SessionRule struct {
	ID        uuid.UUID `db:"id"`
	SessionID uuid.UUID `db:"session_id"`
	RuleText  string    `db:"rule_text"`
	CreatedAt time.Time `db:"created_at"`
}

// SessionParticipant represents a participant in a session
type SessionParticipant struct {
	ID          uuid.UUID         `db:"id"`
	SessionID   uuid.UUID         `db:"session_id"`
	UserID      uuid.UUID         `db:"user_id"`
	Status      ParticipantStatus `db:"status"`
	JoinedAt    time.Time         `db:"joined_at"`
	CancelledAt *time.Time        `db:"cancelled_at"`
	UserName    string            `db:"user_name,omitempty"` // From JOIN with users table
	PlayerLevel PlayerLevel       `db:"player_level,"`       // From JOIN with users table
}

// Court represents a court at a venue

// SessionDetail represents a session with additional details
type SessionDetail struct {
	Session
	VenueName        string               `db:"venue_name"`
	VenueLocation    string               `db:"venue_location"`
	HostName         string               `db:"host_name"`
	HostLevel        PlayerLevel          `db:"host_level"`
	ConfirmedPlayers int                  `db:"confirmed_players"`
	Participants     []SessionParticipant `db:"participants,omitempty"`
	Rules            []SessionRule        `db:"rules,omitempty"`
	Search_vector    string               `db:"search_vector"`
	IsPublic         bool                 `db:"is_public"`
}

type JoinStatus string

const (
	JoinStatusHost JoinStatus = "host"

	JoinStatusConfirmed JoinStatus = "confirmed"

	JoinStatusPending JoinStatus = "pending"

	JoinStatusCancelled JoinStatus = "cancelled"

	JoinStatusNotJoined JoinStatus = "not_joined"
)

package interfaces

import (
	"context"

	"badbuddy/internal/domain/models"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SessionDetail, error)
	Update(ctx context.Context, session *models.Session) error
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.SessionDetail, error)
	Search(ctx context.Context, searchQuery string, filters map[string]interface{}, limit, offset int) ([]models.SessionDetail, error)
	AddParticipant(ctx context.Context, participant *models.SessionParticipant) error
	UpdateParticipantStatus(ctx context.Context, sessionID, userID uuid.UUID, status models.ParticipantStatus) error
	GetParticipants(ctx context.Context, sessionID uuid.UUID) ([]models.SessionParticipant, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error)
	GetMyJoinedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error)
	GetMyHostedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.SessionDetail, error)
	GetJoinStatus(ctx context.Context, userID, venueID uuid.UUID) (models.JoinStatus, error)
}

package session

import (
	"context"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"

	"github.com/google/uuid"
)

type UseCase interface {
	CreateSession(ctx context.Context, hostID uuid.UUID, req requests.CreateSessionRequest) (*responses.SessionResponse, error)
	UpdateSession(ctx context.Context, sessionID uuid.UUID, hostID uuid.UUID, req requests.UpdateSessionRequest) error
	GetSession(ctx context.Context, id uuid.UUID) (*responses.SessionResponse, error)
	GetSessionStatus(ctx context.Context, id uuid.UUID, user_id uuid.UUID) (*responses.SessionResponse, error)

	ListSessions(ctx context.Context, filters map[string]interface{}, limit, offset int) (*responses.SessionListResponse, error)
	SearchSessions(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) (*responses.SessionListResponse, error)
	JoinSession(ctx context.Context, sessionID, userID uuid.UUID, req requests.JoinSessionRequest) error
	LeaveSession(ctx context.Context, sessionID, userID uuid.UUID) error
	CancelSession(ctx context.Context, sessionID, hostID uuid.UUID) error
	GetUserSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]responses.SessionResponse, error)
	ChangeParticipantStatus(ctx context.Context, sessionID, hostID uuid.UUID, req requests.ChangeParticipantStatusRequest) error
	GetSessionParticipants(ctx context.Context, sessionID uuid.UUID) ([]responses.ParticipantResponse, error)
	GetMyJoinedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]responses.SessionResponse, error)
	GetMyHostedSessions(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]responses.SessionResponse, error)
}

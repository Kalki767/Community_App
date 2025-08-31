package usecaseinterfaces

import (
	"auth/internal/delivery/http/dto"

	"github.com/google/uuid"
)

type SessionUsecaseInterface interface {
    ListActiveSessions(userID uuid.UUID) ([]*dto.SessionDTO, error)
    GetSession(sessionID uuid.UUID) (*dto.SessionDTO, error)
    Logout(sessionID uuid.UUID) error
    LogoutAllExcept(userID uuid.UUID, keepSessionID uuid.UUID) error
    Refresh(refreshToken string) (string, error)
}
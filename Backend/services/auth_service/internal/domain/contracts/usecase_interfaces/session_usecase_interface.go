package usecaseinterfaces

import (
	"auth/internal/delivery/http/dto"

	"github.com/google/uuid"
)

type SessionUsecaseInterface interface {
	ListActiveSessions(userID uuid.UUID) ([]*dto.SessionResponseDTO, error)
	GetSession(sessionID uuid.UUID) (*dto.SessionResponseDTO, error)
	Logout(sessionID uuid.UUID) error
	LogoutAllExcept(userID uuid.UUID, keepSessionID uuid.UUID) error
	Refresh(refreshToken string) (string, error)
    IsSessionActive(sessionID uuid.UUID) (bool, error)
}

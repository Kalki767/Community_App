package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionResponseDTO struct {
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
    UserAgent  string `json:"user_agent"`
    IP         string `json:"ip"`
    LastUsedAt time.Time `json:"last_used_at"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type MessageResponse struct {
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}

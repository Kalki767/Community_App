package dto

import (
	"time"

	"github.com/google/uuid"
)

type SessionDTO struct {
	ID uuid.UUID
	UserID uuid.UUID
	TokenHash  string     `gorm:"not null"`
    ExpiresAt  time.Time  `gorm:"index;not null"`
    UserAgent  string
    IP         string
    LastUsedAt time.Time
}
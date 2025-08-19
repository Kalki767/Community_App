package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
    ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID     uuid.UUID  `gorm:"type:uuid;index;not null"`
    TokenHash  string     `gorm:"not null"`
    ExpiresAt  time.Time  `gorm:"index;not null"`
    UserAgent  string
    IP         string
    LastUsedAt time.Time
    CreatedAt  time.Time  `gorm:"autoCreateTime"`
    RevokedAt  *time.Time

    User       User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` 
}

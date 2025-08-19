package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName      string    `gorm:"not null"`
	Email         string    `gorm:"uniqueIndex;not null"`
	Username      string    `gorm:"uniqueIndex;not null"`
	PhoneNumber   string    `gorm:"uniqueIndex;not null"`
	Country       string    `gorm:"not null"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"not null;default:user"`
	AcceptedTerms bool      `gorm:"not null;default:false"`
	IsVerified    bool      `gorm:"not null;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	RefreshSessions []*Session `gorm:"foreignKey:UserID"`
}

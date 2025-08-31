package repointerfaces

import (
	"auth/internal/domain/entity"

	"github.com/google/uuid"
)

type SessionRepoInterface interface {
	AddSession(session *entity.Session) (*entity.Session, error)
	GetById(Id uuid.UUID) (*entity.Session, error)
	GetAll(userId uuid.UUID)([]*entity.Session, error)
	RevokeSession(Id uuid.UUID) error
	RevokeForAllUser(userId uuid.UUID) error
	RevokeAllExceptCurrent(userId uuid.UUID, keepsessionId uuid.UUID) error
	UpdateLastUsed(Id uuid.UUID) error
}
package repointerfaces

import (
	"auth/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepoInterface interface {
	Create(user *entity.User) (*entity.User, error)
	GetById(Id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)	
}
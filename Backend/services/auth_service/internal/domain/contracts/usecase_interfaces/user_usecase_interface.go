package usecaseinterfaces

import (
	"auth/internal/delivery/http/dto"

	"github.com/google/uuid"
)

type UserUsecaseInterface interface {
	Register(user *dto.RegisterUser) (*dto.UserDto, error)
	Login(identification string, password string) (*dto.UserDto, string, string, error)
	GetUserProfile(Id uuid.UUID) (*dto.UserDto, error)
	IsVerifiedUser(Id uuid.UUID) (bool, error)
}
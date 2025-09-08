package repository

import (
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	"auth/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repointerfaces.UserRepoInterface{
	return &UserRepo{db:db}
}

func (repo *UserRepo) Create(user *entity.User) (*entity.User, error){
	err := repo.db.Create(user).Error
	if err != nil{
		return nil, err
	}
	return user, nil
}
func (repo *UserRepo) GetById(Id uuid.UUID) (*entity.User, error){
	var user entity.User
	err := repo.db.Where("id = ?", Id).First(&user).Error
	if err != nil{
		return nil, err
	}
	return &user,nil
}
func (repo *UserRepo) GetByEmail(email string) (*entity.User, error){
	var user entity.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil{
		return nil, err
	}
	return &user,nil
}
func (repo *UserRepo) GetByUsername(username string) (*entity.User, error){
	var user entity.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil{
		return nil, err
	}
	return &user,nil
}
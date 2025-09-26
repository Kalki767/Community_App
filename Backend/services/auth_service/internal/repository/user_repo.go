package repository

import (
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	"auth/internal/domain/entity"
	"context"
	"encoding/json"
	"errors"

	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
	redis *redis.Client
	ttl time.Duration
}

func NewUserRepo(db *gorm.DB, redis *redis.Client, ttl time.Duration) repointerfaces.UserRepoInterface{
	return &UserRepo{db:db, redis: redis, ttl: ttl}
}

func (repo *UserRepo) Create(user *entity.User) (*entity.User, error){
	ctx := context.Background()
	err := repo.db.Create(user).Error
	if err != nil{
		return nil, err
	}

	data, _ := json.Marshal(&user)
	repo.redis.Set(ctx,"user: "+ user.ID.String(), data, repo.ttl)
	return user, nil
}
func (repo *UserRepo) GetById(Id uuid.UUID) (*entity.User, error){
	ctx := context.Background()
	key := "user: " + Id.String()

	val, err := repo.redis.Get(ctx,key).Result()
	if err == nil{
		var user entity.User
		if jsonErr := json.Unmarshal([]byte(val), &user); jsonErr == nil{
			return &user, nil
		} else{
			return nil, errors.New("unable to marshal json response from redis")
		}
	}
	var user entity.User
	err = repo.db.Where("id = ?", Id).First(&user).Error
	if err != nil{
		return nil, err
	}

	data, _ := json.Marshal(&user)
	repo.redis.Set(ctx, "user: " + Id.String(), data, repo.ttl)
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
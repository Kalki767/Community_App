package usecase

import (
	"auth/internal/delivery/http/dto"
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"
	"auth/internal/domain/entity"
	"auth/internal/services"
	"errors"
	// "fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"auth/helper"
)

type UserUsecase struct {
	user_repo repointerfaces.UserRepoInterface
	session_repo repointerfaces.SessionRepoInterface
	tokenservice services.TokenService
}

func NewUserUsecase(user_repo repointerfaces.UserRepoInterface, session_repo repointerfaces.SessionRepoInterface, tokenservice services.TokenService) usecaseinterfaces.UserUsecaseInterface{
	return &UserUsecase{user_repo:user_repo, session_repo: session_repo, tokenservice: tokenservice}
}

func (uc *UserUsecase) Register(userdto *dto.RegisterUser) (*dto.UserDto, error){
	user := &entity.User{
		FullName: userdto.FullName,
		ID: uuid.New(),
		Email: userdto.Email,
		Username: userdto.Username,
		Country: userdto.Country,
		PhoneNumber: userdto.PhoneNumber,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userdto.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user.PasswordHash = string(hashedPassword)
	created_user, err := uc.user_repo.Create(user)

	if err != nil {
		return nil, err
	}
	created_user_dto := &dto.UserDto{
		ID: created_user.ID,
		FullName: created_user.FullName,
		Email: created_user.Email,
		Username: created_user.Username,
		PhoneNumber: created_user.PhoneNumber,
		Country: created_user.Country,
		IsVerified: created_user.IsVerified,
	}

    return created_user_dto, nil
}

func (uc *UserUsecase)	Login(identification string, password string) (*dto.UserDto,string, string, error){
	user, err := uc.user_repo.GetByEmail(identification)
	if err != nil {
		user, err = uc.user_repo.GetByUsername(identification)
		if err != nil{
			return nil,"", "",errors.New("user not found")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil{
		return nil, "", "",errors.New("invalid credentials")
	}

	sessionId := uuid.New()
	refreshToken, err := uc.tokenservice.GenerateRefreshToken(user.ID,sessionId)
	if err != nil {
		return nil, "", "",err
	}
	
	
	HashedToken := helper.HashTokenSHA512(refreshToken)
	// if err != nil {
	// 	return nil, "", "",fmt.Errorf("cannot hash refresh token: %w", err)
	// }
	session := &entity.Session{
        ID:        sessionId,
        UserID:    user.ID,
        TokenHash: string(HashedToken), // store hash
        ExpiresAt: time.Now().UTC().Add(30 * 24 * time.Hour), // 30 days
        LastUsedAt: time.Now().UTC(),
        CreatedAt: time.Now().UTC(),
    }

	_, err = uc.session_repo.AddSession(session)
	if err != nil {
		return nil, "", "",err
	}

	var roles []string
	access_token, err := uc.tokenservice.GenerateAccessToken(user.ID,session.ID,roles)

	if err != nil {
		return nil, "", "", err
	}

	user_dto := &dto.UserDto{
		ID: user.ID,
		FullName: user.FullName,
		Email: user.Email,
		Username: user.Username,
		PhoneNumber: user.PhoneNumber,
		Country: user.Country,
		IsVerified: user.IsVerified,
	}
	
	return user_dto, refreshToken, access_token, nil
	
	
}
func (uc *UserUsecase)	GetUserProfile(Id uuid.UUID) (*dto.UserDto, error){
	user, err := uc.user_repo.GetById(Id)
	if err != nil {
		return nil, err
	}

	user_dto := &dto.UserDto{
		ID: user.ID,
		FullName: user.FullName,
		Email: user.Email,
		Username: user.Username,
		PhoneNumber: user.PhoneNumber,
		Country: user.Country,
		IsVerified: user.IsVerified,
	}
	return user_dto, nil
	
}
func (uc *UserUsecase)	IsVerifiedUser(Id uuid.UUID) (bool, error){
	user, err := uc.user_repo.GetById(Id)
	if err != nil {
		return false, err
	}
	
	return user.IsVerified, nil
}
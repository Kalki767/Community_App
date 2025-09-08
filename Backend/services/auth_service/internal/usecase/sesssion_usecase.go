package usecase

import (
	"auth/helper"
	"auth/internal/delivery/http/dto"
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	usecaseinterfaces "auth/internal/domain/contracts/usecase_interfaces"
	"auth/internal/services"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	// "golang.org/x/crypto/bcrypt"
)

type SessionUsecase struct {
	repo repointerfaces.SessionRepoInterface
	tokenService services.TokenService
}

func NewSessionUsecase(repo repointerfaces.SessionRepoInterface, tokenservice services.TokenService) usecaseinterfaces.SessionUsecaseInterface{
	return &SessionUsecase{repo:repo, tokenService: tokenservice}
}

func (uc *SessionUsecase) ListActiveSessions(userID uuid.UUID) ([]*dto.SessionResponseDTO, error){
	sessions, err := uc.repo.GetAll(userID)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		return []*dto.SessionResponseDTO{}, nil
	}

	var sessionsDto []*dto.SessionResponseDTO
	for _, session := range sessions {
		sessionDto := &dto.SessionResponseDTO{
			ID: session.ID,
			UserID: session.UserID,
			UserAgent: session.UserAgent,
			IP: session.IP,
			LastUsedAt: session.LastUsedAt,
		}

		sessionsDto = append(sessionsDto, sessionDto)
	}

	return sessionsDto, nil
}
func (uc *SessionUsecase) GetSession(sessionID uuid.UUID) (*dto.SessionResponseDTO, error){
	session, err := uc.repo.GetById(sessionID)
	if err != nil {
		return nil, err
	}
	sessionDto := &dto.SessionResponseDTO{
			ID: session.ID,
			UserID: session.UserID,
			UserAgent: session.UserAgent,
			IP: session.IP,
			LastUsedAt: session.LastUsedAt,
		}
	return sessionDto, nil
}
func (uc *SessionUsecase) Logout(sessionID uuid.UUID) error{
	return uc.repo.RevokeSession(sessionID)
}
func (uc *SessionUsecase) LogoutAllExcept(userID uuid.UUID, keepSessionID uuid.UUID) error{
	return uc.repo.RevokeAllExceptCurrent(userID,keepSessionID)
}
func (uc *SessionUsecase) Refresh(refreshToken string) (string, error){
	// 1. Parse refresh token using the injected service
    claims, err := uc.tokenService.ParseRefreshToken(refreshToken)
    if err != nil {
        return "",  err
    }

    // 2. Load session from DB
    session, err := uc.repo.GetById(claims.SessionID)
    if err != nil {
        return "",err
    }
    if session.RevokedAt != nil || time.Now().After(session.ExpiresAt) {
        return "",  errors.New("session expired or revoked")
    }

    // 3. Verify refresh token hash (if stored)
	
	if matched := helper.CompareTokenSHA512(refreshToken,session.TokenHash); !matched{
		return "",  errors.New("refresh token mismatch")
	}
    
	// 4. Generate new tokens
    accessToken, err := uc.tokenService.GenerateAccessToken(session.UserID, session.ID, []string{})
    if err != nil {
        return "", fmt.Errorf("failed to create access token: %w", err)
    }

    return accessToken, nil
}

func (uc *SessionUsecase) IsSessionActive(sessionID uuid.UUID) (bool, error){
	session, err := uc.repo.GetById(sessionID)
	if err != nil {
		return false, fmt.Errorf("something went wrong %w", err)
	}

	if session.RevokedAt != nil || session.ExpiresAt.Before(time.Now().UTC()){
		return false, nil
	}

	return true, nil
}
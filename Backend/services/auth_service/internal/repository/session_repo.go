package repository

import (
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	"auth/internal/domain/entity"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) repointerfaces.SessionRepoInterface {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) AddSession(session *entity.Session) (*entity.Session, error) {
	err := repo.db.Create(session).Error
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (repo *SessionRepository) GetAll(userId uuid.UUID) ([]*entity.Session, error) {
	var sessions []*entity.Session
	now := time.Now().UTC()
	fmt.Println("Now is:", now)
	err := repo.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userId, time.Now().UTC()).Order("last_used_at DESC").Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
func (repo *SessionRepository) GetById(Id uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	err := repo.db.Where("id = ?", Id).First(&session, Id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
func (repo *SessionRepository) RevokeSession(Id uuid.UUID) error {
	now := time.Now()
	err := repo.db.Model(&entity.Session{}).Where("id = ?", Id).Update("revoked_at", &now).Error
	if err != nil {
		return err
	}
	return nil
}
func (repo *SessionRepository) RevokeForAllUser(userId uuid.UUID) error {
	now := time.Now()
	return repo.db.Model(&entity.Session{}).
		Where("user_id = ?", userId).
		Update("revoked_at", &now).Error
}
func (repo *SessionRepository) RevokeAllExceptCurrent(userId uuid.UUID, keepSessionId uuid.UUID) error {
	now := time.Now()
	err := repo.db.Model(&entity.Session{}).
		Where("user_id = ? AND id != ?", userId, keepSessionId).
		Update("revoked_at", &now).Error

	if err != nil {
		return err
	}
	return nil
}
func (repo *SessionRepository) UpdateLastUsed(Id uuid.UUID) error {
	return repo.db.Model(&entity.Session{}).
		Where("id = ?", Id).
		Update("last_used_at", time.Now()).Error
}

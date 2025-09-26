package repository

import (
	repointerfaces "auth/internal/domain/contracts/repo_interfaces"
	"auth/internal/domain/entity"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
	redis *redis.Client
	ttl time.Duration
}

func NewSessionRepository(db *gorm.DB, redis *redis.Client, ttl time.Duration) repointerfaces.SessionRepoInterface {
	return &SessionRepository{db: db, redis: redis, ttl: ttl}
}

func (repo *SessionRepository) AddSession(session *entity.Session) (*entity.Session, error) {
    if err := repo.db.Create(session).Error; err != nil {
        return nil, err
    }

    ctx := context.Background()
    ttl := time.Until(session.ExpiresAt)
    if ttl <= 0 {
        // already expired — don't cache
        return session, nil
    }

    if err := repo.redis.Set(ctx, "refresh_token:"+session.TokenHash, session.ID.String(), ttl).Err(); err != nil {
        // log or return error depending on how critical redis caching is
        return nil, err
    }

    return session, nil
}

func (repo *SessionRepository) GetAll(userId uuid.UUID) ([]*entity.Session, error) {
	ctx := context.Background()
	var sessions []*entity.Session

	err := repo.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userId, time.Now().UTC()).Order("last_used_at DESC").Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	var ttl time.Duration
	// refresh cache
	for _, s := range sessions {
		ttl = time.Until(s.ExpiresAt)
		if err = repo.redis.Set(ctx, "refresh_token:" + s.TokenHash, s.ID.String(),ttl).Err(); err != nil{
			return nil, err
		}
	}


	return sessions, nil
}
func (repo *SessionRepository) GetById(id uuid.UUID) (*entity.Session, error) {
    var session entity.Session
    if err := repo.db.First(&session, id).Error; err != nil {
        return nil, err
    }
    if session.RevokedAt != nil || session.ExpiresAt.Before(time.Now().UTC()) {
        return nil, errors.New("session expired or revoked")
    }

    // refresh redis mapping for the refresh token
    ttl := time.Until(session.ExpiresAt)
    if ttl > 0 {
        ctx := context.Background()
        _ = repo.redis.Set(ctx, "refresh_token:"+session.TokenHash, session.ID.String(), ttl).Err()
    }
    return &session, nil
}

func (repo *SessionRepository) RevokeSession(id uuid.UUID) error {
    ctx := context.Background()
    now := time.Now().UTC()

    if err := repo.db.Model(&entity.Session{}).Where("id = ?", id).Update("revoked_at", &now).Error; err != nil {
        return err
    }

    // find token hash in DB (safe fallback)
    var s entity.Session
    if err := repo.db.First(&s, id).Error; err != nil {
        return err
    }

    // delete redis keys (safe even if they don't exist)
    keys := []string{}
    if s.TokenHash != "" {
        keys = append(keys, "refresh_token:"+s.TokenHash)
    }
    return repo.redis.Del(ctx, keys...).Err()
}


func (repo *SessionRepository) RevokeForAllUser(userId uuid.UUID) error {
	ctx := context.Background()
	now := time.Now().UTC()

	// 1) Load sessions that are currently active (not revoked yet)
	var sessions []entity.Session
	if err := repo.db.Where("user_id = ? AND revoked_at IS NULL", userId).Find(&sessions).Error; err != nil {
		return err
	}

	// If no sessions found, still update DB just to be safe (idempotent)
	if len(sessions) == 0 {
		// mark everything revoked (this is idempotent)
		if err := repo.db.Model(&entity.Session{}).
			Where("user_id = ?", userId).
			Update("revoked_at", &now).Error; err != nil {
			return err
		}
		return nil
	}

	// 2) Collect redis keys to delete (refresh_token:..., session_refresh:...)
	keys := make([]string, 0, len(sessions)*2)
	for _, s := range sessions {
		if s.TokenHash != "" {
			keys = append(keys, "refresh_token:"+s.TokenHash)
		}
		// if you maintain reverse mapping (session_refresh), remove that too
		keys = append(keys, "session_refresh:"+s.ID.String())
	}

	// 3) Mark sessions revoked in DB
	if err := repo.db.Model(&entity.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userId).
		Update("revoked_at", &now).Error; err != nil {
		return err
	}

	// 4) Delete Redis keys in batches
	if len(keys) > 0 {
		if err := delKeysInChunks(ctx, repo.redis, keys); err != nil {
			return err
		}
	}

	return nil
}

func (repo *SessionRepository) RevokeAllExceptCurrent(userId uuid.UUID, keepSessionId uuid.UUID) error {
	now := time.Now().UTC()
	err := repo.db.Model(&entity.Session{}).
		Where("user_id = ? AND id != ?", userId, keepSessionId).
		Update("revoked_at", &now).Error

	if err != nil {
		return err
	}
	var session entity.Session
	err = repo.db.Where("id = ?", keepSessionId).First(&session).Error
	if err != nil {
		return  err
	}

	ctx := context.Background()
	ttl := time.Until(session.ExpiresAt)
	repo.redis.Set(ctx, "refresh_token:" + session.TokenHash, session.ID.String(),ttl)
	
	return nil
}
func (repo *SessionRepository) UpdateLastUsed(Id uuid.UUID) error {
	err := repo.db.Model(&entity.Session{}).
		Where("id = ?", Id).
		Update("last_used_at", time.Now()).Error
	if err != nil {
		return err
	}

	var session entity.Session
	if err := repo.db.Where("id = ?", Id).First(&session).Error; err == nil {
		ctx := context.Background()
		ttl := time.Until(session.ExpiresAt)
		
		if ttl > 0 {
			repo.redis.Set(ctx, "refresh_token:" + session.TokenHash, session.ID.String(),ttl)
        } else {
            // Already expired → remove from Redis
			repo.redis.Del(ctx, "refresh_token:" + session.TokenHash)
        }
	}

	return nil
}

// helper: delete keys in chunks to avoid very large redis.Del() calls
func delKeysInChunks(ctx context.Context, r *redis.Client, keys []string) error {
	const chunkSize = 500
	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}
		if err := r.Del(ctx, keys[i:end]...).Err(); err != nil {
			return err
		}
	}
	return nil
}

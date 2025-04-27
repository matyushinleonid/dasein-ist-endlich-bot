package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/redis"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

var ErrSessionNotFound = fmt.Errorf("session not found")

type SessionRepository struct {
	redisClient redis.Client
}

func NewSessionRepository(client redis.Client) *SessionRepository {
	return &SessionRepository{redisClient: client}
}

func (r *SessionRepository) Get(ctx context.Context, userID int64) (*model.Session, error) {
	var sess model.Session
	err := r.redisClient.Load(ctx, userID, &sess)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("redis: could not load session %d: %w", userID, err)
	}
	return &sess, nil
}

func (r *SessionRepository) Save(ctx context.Context, userID int64, sess *model.Session) error {
	if err := r.redisClient.Save(ctx, userID, sess); err != nil {
		return fmt.Errorf("redis: could not save session %d: %w", userID, err)
	}
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, userID int64) error {
	if err := r.redisClient.Delete(ctx, userID); err != nil {
		return fmt.Errorf("redis: could not delete session %d: %w", userID, err)
	}
	return nil
}

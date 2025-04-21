package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
	"github.com/redis/go-redis/v9"
)

type RealClient struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRealClient(cfg config.RedisConfig) *RealClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RealClient{rdb: rdb, ttl: cfg.TTL}
}

func (c *RealClient) key(id int64) string {
	return fmt.Sprintf("sess:%d", id)
}

func (c *RealClient) Save(ctx context.Context, id int64, sess *model.Session) error {
	data, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("session.Marshal: %w", err)
	}
	if err := c.rdb.Set(ctx, c.key(id), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}
	return nil
}

func (c *RealClient) Load(ctx context.Context, id int64) (*model.Session, error) {
	data, err := c.rdb.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("session not found for id %d", id)
		}
		return nil, fmt.Errorf("redis.Get: %w", err)
	}
	var sess model.Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("session.Unmarshal: %w", err)
	}
	return &sess, nil
}

func (c *RealClient) Delete(ctx context.Context, id int64) error {
	if err := c.rdb.Del(ctx, c.key(id)).Err(); err != nil {
		return fmt.Errorf("redis.Del: %w", err)
	}
	return nil
}

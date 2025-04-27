package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/retry"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
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

func (c *RealClient) Save(ctx context.Context, key int64, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	err = retry.Do(ctx, retry.DefaultConfig(), func() error {
		return c.rdb.Set(ctx, c.key(key), data, c.ttl).Err()
	})
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}
	return nil
}

func (c *RealClient) Load(ctx context.Context, key int64, result interface{}) error {
	var data []byte
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		data, e = c.rdb.Get(ctx, c.key(key)).Bytes()
		return e
	})
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return fmt.Errorf("redis.Get: %w", err)
	}
	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	return nil
}

func (c *RealClient) Delete(ctx context.Context, key int64) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		return c.rdb.Del(ctx, c.key(key)).Err()
	})
	if err != nil {
		return fmt.Errorf("redis.Del: %w", err)
	}
	return nil
}

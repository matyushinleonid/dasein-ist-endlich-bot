package redis

import (
	"context"
	"errors"
)

var (
	ErrNotFound    = errors.New("data not found")
	ErrInvalidType = errors.New("invalid type: must be pointer to struct or slice")
)

type Client interface {
	Save(ctx context.Context, key int64, value interface{}) error
	Load(ctx context.Context, key int64, result interface{}) error
	Delete(ctx context.Context, key int64) error
}

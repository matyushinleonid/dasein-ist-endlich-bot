package redis

import (
	"context"
)

type Client interface {
	Save(ctx context.Context, key int64, value interface{}) error
	Load(ctx context.Context, key int64, result interface{}) error
	Delete(ctx context.Context, key int64) error
}

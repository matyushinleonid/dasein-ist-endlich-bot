package mongo

import (
	"context"
)

type Client interface {
	Get(ctx context.Context, key int64, result interface{}) error
	Create(ctx context.Context, doc interface{}) (int64, error)
	Update(ctx context.Context, key int64, update interface{}) (int64, error)
	Delete(ctx context.Context, key int64) (int64, error)
}

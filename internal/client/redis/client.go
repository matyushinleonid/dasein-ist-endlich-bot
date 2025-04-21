package redis

import (
	"context"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

type Client interface {
	Save(ctx context.Context, id int64, sess *model.Session) error
	Load(ctx context.Context, id int64) (*model.Session, error)
	Delete(ctx context.Context, id int64) error
}

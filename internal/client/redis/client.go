package redis

import (
	"context"
)

type Session struct {
	Stage   int      `json:"stage"`
	Answers []string `json:"answers"`
}

type Client interface {
	Save(ctx context.Context, id int64, sess *Session) error
	Load(ctx context.Context, id int64) (*Session, error)
	Delete(ctx context.Context, id int64) error
}

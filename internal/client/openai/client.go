package openai

import (
	"context"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Client interface {
	SendText(ctx context.Context, userID int64, userMessage string) (string, error)
}

func NewClient(cfg config.OpenAIConfig) Client {
	if cfg.Dummy {
		return NewDummyClient()
	}
	return NewRealClient(cfg)
}

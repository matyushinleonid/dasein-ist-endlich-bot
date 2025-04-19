package openai

import (
	"context"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Client interface {
	SendText(ctx context.Context, prompt string) (string, error)
}

func NewClient(cfg config.OpenAIConfig) Client {
	if cfg.Dummy {
		return NewDummyClient(cfg)
	}
	return NewRealClient(cfg)
}

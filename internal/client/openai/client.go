package openai

import (
	"context"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Client interface {
	SendText(ctx context.Context, userID int64, userMessage string) (string, error)
	SendJSON(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}) (string, error)
	SendJSONUnmarshal(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}, out any) error
}

func NewClient(cfg config.OpenAIConfig) Client {
	if cfg.Dummy {
		return NewDummyClient()
	}
	return NewRealClient(cfg)
}

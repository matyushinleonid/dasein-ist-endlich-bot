package openai

import (
	"context"
	"fmt"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DummyClient struct {
	model string
}

func NewDummyClient(cfg config.OpenAIConfig) *DummyClient {
	return &DummyClient{
		model: cfg.Model,
	}
}

func (d *DummyClient) SendText(ctx context.Context, prompt string) (string, error) {
	return fmt.Sprintf("Model: %s \n Dummy response to: \n %s", d.model, prompt), nil
}

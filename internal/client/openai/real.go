package openai

import (
	"context"
	"fmt"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type RealClient struct {
	client *openai.Client
	model  string
}

func NewRealClient(cfg config.OpenAIConfig) *RealClient {
	cli := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)
	return &RealClient{
		client: &cli,
		model:  cfg.Model,
	}
}

func (r *RealClient) SendText(ctx context.Context, prompt string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Model: r.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}
	resp, err := r.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned by openai")
	}
	return resp.Choices[0].Message.Content, nil
}

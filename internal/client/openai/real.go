package openai

import (
	"context"
	"fmt"
	"strconv"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

type RealClient struct {
	client           *openai.Client
	model            string
	developerMessage string
}

func NewRealClient(cfg config.OpenAIConfig) *RealClient {
	cli := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)
	return &RealClient{
		client:           &cli,
		model:            cfg.Model,
		developerMessage: cfg.DeveloperMessage,
	}
}

func (r *RealClient) SendText(ctx context.Context, userID int64, userMessage string) (string, error) {
	params := openai.ChatCompletionNewParams{
		Model: r.model,
		User:  param.NewOpt(strconv.FormatInt(userID, 10)),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage(r.developerMessage),
			openai.UserMessage(userMessage),
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

package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

func (r *RealClient) SendJSON(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}) (string, error) {
	jsonSchema := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   schemaName,
		Schema: schema,
		Strict: param.NewOpt(true),
	}

	params := openai.ChatCompletionNewParams{
		Model: r.model,
		User:  param.NewOpt(strconv.FormatInt(userID, 10)),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: jsonSchema,
			},
		},
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage(r.developerMessage),
			openai.UserMessage(userMessage),
			openai.DeveloperMessage("Current datetime (When user have answered the questions): " + time.Now().Format("2006-01-02 15:04:05")),
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

func (r *RealClient) SendJSONUnmarshal(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}, out any) error {
	raw, err := r.SendJSON(ctx, userID, userMessage, schemaName, schema)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(raw), out); err != nil {
		return fmt.Errorf("failed to unmarshal OpenAI JSON into %T: %w", out, err)
	}
	return nil
}

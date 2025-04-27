package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type DummyClient struct {
	Sleep                  time.Duration
	SendJSONOutput         string
	ErrOnSendJSON          error
	ErrOnSendJSONUnmarshal error
}

func NewDummyClient(opts ...time.Duration) *DummyClient {
	sleep := 2 * time.Second
	if len(opts) > 0 {
		sleep = opts[0]
	}
	return &DummyClient{Sleep: sleep}
}

func (d *DummyClient) SendText(ctx context.Context, userID int64, userMessage string) (string, error) {
	time.Sleep(d.Sleep)
	return fmt.Sprintf("Dummy response to: \n %s", userMessage), nil
}

func (d *DummyClient) SendJSON(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}) (string, error) {
	time.Sleep(d.Sleep)
	if d.ErrOnSendJSON != nil {
		return "", d.ErrOnSendJSON
	}
	fake := map[string]interface{}{
		"schemaName":    schemaName,
		"originalQuery": userMessage,
		"echoedSchema":  schema,
		"dummy":         "value",
	}
	b, err := json.MarshalIndent(fake, "", "  ")
	if err != nil {
		return "", err
	}
	if d.SendJSONOutput != "" {
		return d.SendJSONOutput, nil
	}
	return string(b), nil
}

func (d *DummyClient) SendJSONUnmarshal(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}, out any) error {
	raw, err := d.SendJSON(ctx, userID, userMessage, schemaName, schema)
	if err != nil {
		return err
	}
	if d.ErrOnSendJSONUnmarshal != nil {
		return d.ErrOnSendJSONUnmarshal
	}
	if err := json.Unmarshal([]byte(raw), out); err != nil {
		return fmt.Errorf("dummy unmarshal error into %T: %w", out, err)
	}
	return nil
}

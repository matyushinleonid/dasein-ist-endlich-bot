package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type DummyClient struct {
	sleep time.Duration
}

func NewDummyClient() *DummyClient {
	return &DummyClient{
		sleep: 2 * time.Second,
	}
}

func (d *DummyClient) SendText(ctx context.Context, userID int64, userMessage string) (string, error) {
	time.Sleep(d.sleep)
	return fmt.Sprintf("Dummy response to: \n %s", userMessage), nil
}

func (d *DummyClient) SendJSON(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}) (string, error) {
	time.Sleep(d.sleep)
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
	return string(b), nil
}

func (d *DummyClient) SendJSONUnmarshal(ctx context.Context, userID int64, userMessage, schemaName string, schema map[string]interface{}, out any) error {
	raw, err := d.SendJSON(ctx, userID, userMessage, schemaName, schema)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(raw), out); err != nil {
		return fmt.Errorf("dummy unmarshal error into %T: %w", out, err)
	}
	return nil
}

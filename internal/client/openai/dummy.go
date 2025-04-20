package openai

import (
	"context"
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

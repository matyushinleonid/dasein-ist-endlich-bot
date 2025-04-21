package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

type DummyClient struct {
	mu      sync.RWMutex
	storage map[int64]*model.Session
}

func NewDummyClient() *DummyClient {
	return &DummyClient{
		storage: make(map[int64]*model.Session),
	}
}

func (c *DummyClient) Save(ctx context.Context, id int64, sess *model.Session) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	copySess := *sess
	c.storage[id] = &copySess
	return nil
}

func (c *DummyClient) Load(ctx context.Context, id int64) (*model.Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	sess, ok := c.storage[id]
	if !ok {
		return nil, fmt.Errorf("session not found for id %d", id)
	}
	copySess := *sess
	return &copySess, nil
}

func (c *DummyClient) Delete(ctx context.Context, id int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, id)
	return nil
}

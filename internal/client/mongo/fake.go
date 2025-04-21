package mongo

import (
	"context"
	"errors"
	"sync"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/record"
)

type DummyClient struct {
	mu      sync.RWMutex
	storage map[int64]record.Record
}

func NewDummyClient() *DummyClient {
	return &DummyClient{storage: make(map[int64]record.Record)}
}

func (c *DummyClient) Get(ctx context.Context, key int64, result interface{}) error {
	c.mu.RLock()
	rec, ok := c.storage[key]
	c.mu.RUnlock()
	if !ok {
		return errors.New("document not found")
	}
	r, ok := result.(*record.Record)
	if !ok {
		return errors.New("result argument must be *record.Record")
	}
	*r = rec
	return nil
}

func (c *DummyClient) Create(ctx context.Context, doc interface{}) (int64, error) {
	r, ok := doc.(record.Record)
	if !ok {
		return 0, errors.New("doc must be record.Record")
	}
	c.mu.Lock()
	c.storage[r.ID] = r
	c.mu.Unlock()
	return r.ID, nil
}

func (c *DummyClient) Update(ctx context.Context, key int64, update interface{}) (int64, error) {
	upd, ok := update.(record.Record)
	if !ok {
		return 0, errors.New("update must be record.Record")
	}
	c.mu.Lock()
	orig, exists := c.storage[key]
	if !exists {
		c.mu.Unlock()
		return 0, errors.New("document not found")
	}
	orig.DaysLeft = upd.DaysLeft
	orig.Calculated = upd.Calculated
	c.storage[key] = orig
	c.mu.Unlock()
	return 1, nil
}

func (c *DummyClient) Delete(ctx context.Context, key int64) (int64, error) {
	c.mu.Lock()
	_, existed := c.storage[key]
	if existed {
		delete(c.storage, key)
		c.mu.Unlock()
		return 1, nil
	}
	c.mu.Unlock()
	return 0, nil
}

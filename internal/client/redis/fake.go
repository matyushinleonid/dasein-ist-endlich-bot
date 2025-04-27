package redis

import (
	"context"
	"reflect"
	"sync"
)

type DummyClient struct {
	mu          sync.RWMutex
	storage     map[int64]interface{}
	ErrOnSave   error
	ErrOnLoad   error
	ErrOnDelete error
}

func NewDummyClient() *DummyClient {
	return &DummyClient{storage: make(map[int64]interface{})}
}

func (c *DummyClient) Save(ctx context.Context, key int64, value interface{}) error {
	if c.ErrOnSave != nil {
		return c.ErrOnSave
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || (rv.Elem().Kind() != reflect.Struct && rv.Elem().Kind() != reflect.Slice) {
		return ErrInvalidType
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] = value
	return nil
}

func (c *DummyClient) Load(ctx context.Context, key int64, result interface{}) error {
	if c.ErrOnLoad != nil {
		return c.ErrOnLoad
	}
	c.mu.RLock()
	v, ok := c.storage[key]
	c.mu.RUnlock()
	if !ok {
		return ErrNotFound
	}
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || (rv.Elem().Kind() != reflect.Struct && rv.Elem().Kind() != reflect.Slice) {
		return ErrInvalidType
	}
	val := reflect.ValueOf(v)
	if val.Type() != rv.Type() {
		return ErrInvalidType
	}
	rv.Elem().Set(val.Elem())
	return nil
}

func (c *DummyClient) Delete(ctx context.Context, key int64) error {
	if c.ErrOnDelete != nil {
		return c.ErrOnDelete
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, key)
	return nil
}

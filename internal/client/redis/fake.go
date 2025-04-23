package redis

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type DummyClient struct {
	mu      sync.RWMutex
	storage map[int64]interface{}
}

func NewDummyClient() *DummyClient {
	return &DummyClient{storage: make(map[int64]interface{})}
}

func (c *DummyClient) Save(ctx context.Context, key int64, value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || (rv.Elem().Kind() != reflect.Struct && rv.Elem().Kind() != reflect.Slice) {
		return errors.New("value must be pointer to struct or slice")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] = value
	return nil
}

func (c *DummyClient) Load(ctx context.Context, key int64, result interface{}) error {
	c.mu.RLock()
	v, ok := c.storage[key]
	c.mu.RUnlock()
	if !ok {
		return errors.New("data not found for key")
	}
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || (rv.Elem().Kind() != reflect.Struct && rv.Elem().Kind() != reflect.Slice) {
		return errors.New("result argument must be pointer to struct or slice")
	}
	val := reflect.ValueOf(v)
	if val.Type() != rv.Type() {
		return errors.New("stored value type does not match result type")
	}
	rv.Elem().Set(val.Elem())
	return nil
}

func (c *DummyClient) Delete(ctx context.Context, key int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, key)
	return nil
}

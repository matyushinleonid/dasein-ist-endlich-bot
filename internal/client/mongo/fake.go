package mongo

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

func (c *DummyClient) Get(ctx context.Context, key int64, result interface{}) error {
	c.mu.RLock()
	v, ok := c.storage[key]
	c.mu.RUnlock()
	if !ok {
		return errors.New("document not found")
	}
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("result argument must be pointer to struct")
	}
	rv.Elem().Set(reflect.ValueOf(v))
	return nil
}

func (c *DummyClient) Create(ctx context.Context, doc interface{}) (int64, error) {
	rv := reflect.ValueOf(doc)
	if rv.Kind() != reflect.Struct {
		return 0, errors.New("doc must be struct")
	}
	idField := rv.FieldByName("ID")
	if !idField.IsValid() || idField.Kind() != reflect.Int64 {
		return 0, errors.New("doc must have ID int64 field")
	}
	id := idField.Int()
	c.mu.Lock()
	c.storage[id] = doc
	c.mu.Unlock()
	return id, nil
}

func (c *DummyClient) Update(ctx context.Context, key int64, update interface{}) (int64, error) {
	c.mu.RLock()
	_, exists := c.storage[key]
	c.mu.RUnlock()
	if !exists {
		return 0, errors.New("document not found")
	}
	c.mu.Lock()
	c.storage[key] = update
	c.mu.Unlock()
	return 1, nil
}

func (c *DummyClient) Delete(ctx context.Context, key int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, existed := c.storage[key]; existed {
		delete(c.storage, key)
		return 1, nil
	}
	return 0, nil
}

func (c *DummyClient) FindAll(ctx context.Context) (Cursor, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data := make([]interface{}, 0, len(c.storage))
	for _, v := range c.storage {
		data = append(data, v)
	}
	return &dummyCursor{data: data}, nil
}

var ErrBadDecodeTarget = errors.New("decode target must be pointer to struct")

type dummyCursor struct {
	data []interface{}
	idx  int
	err  error
}

func (c *dummyCursor) Next(ctx context.Context) bool {
	if c.err != nil {
		return false
	}
	if c.idx < len(c.data) {
		return true
	}
	return false
}

func (c *dummyCursor) Decode(v interface{}) error {
	if c.idx >= len(c.data) {
		return c.err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return ErrBadDecodeTarget
	}
	rv.Elem().Set(reflect.ValueOf(c.data[c.idx]))
	c.idx++
	return nil
}

func (c *dummyCursor) Err() error {
	return c.err
}

func (c *dummyCursor) Close(ctx context.Context) error {
	return nil
}

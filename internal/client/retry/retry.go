package retry

import (
	"context"
	"time"
)

type Config struct {
	Attempts  int
	BaseDelay time.Duration
	MaxDelay  time.Duration
}

func DefaultConfig() Config {
	return Config{
		Attempts:  3,
		BaseDelay: 100 * time.Millisecond,
		MaxDelay:  2 * time.Second,
	}
}

func Do(ctx context.Context, cfg Config, fn func() error) error {
	var err error
	delay := cfg.BaseDelay
	for i := 0; i < cfg.Attempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay *= 2
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
		err = fn()
		if err == nil {
			return nil
		}
	}
	return err
}

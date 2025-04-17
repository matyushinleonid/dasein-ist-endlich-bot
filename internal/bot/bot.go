package bot

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Bot struct {
	config *config.Config
}

func NewBot(cfg *config.Config) *Bot {
	return &Bot{
		config: cfg,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	logger := logr.FromContextOrDiscard(ctx)
	reqID := uuid.New().String()
	logger = logger.WithValues("request_id", reqID)
	ctx = logr.NewContext(ctx, logger)

	logger.Info("Bot is starting")
	b.doWork(ctx)
	logger.Info("Bot finished")
	return nil
}

func (b *Bot) doWork(ctx context.Context) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("performing work in doWork")
	time.Sleep(24 * time.Hour)
}

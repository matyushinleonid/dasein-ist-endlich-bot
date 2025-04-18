package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Bot struct {
	cfg *config.Config
}

func NewBot(cfg *config.Config) *Bot {
	return &Bot{cfg: cfg}
}

func (b *Bot) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	opts := []bot.Option{
		bot.WithDefaultHandler(b.handleUpdate),
	}
	if b.cfg.Bot.CheckIfUserAllowed {
		opts = append(opts, bot.WithMiddlewares(b.isUserAllowed))
	}

	tg, err := bot.New(b.cfg.Bot.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	logger.Info("starting bot")
	tg.Start(ctx)
	logger.Info("bot stopped")
	return nil
}

func (b *Bot) isUserAllowed(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, tg *bot.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("isUserAllowed")
		userID := update.Message.From.ID
		allowed := false
		for _, allowedUser := range b.cfg.Bot.AllowedUsers {
			if userID == int64(allowedUser) {
				allowed = true
				break
			}
		}
		if !allowed {
			logger.Info("user not allowed", "user_id", userID)
			return
		}
		next(ctx, tg, update)
	}
}

func (b *Bot) handleUpdate(ctx context.Context, tg *bot.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("handleUpdate")
	if update.Message == nil {
		return
	}
	logger.Info("received message",
		"chat_id", update.Message.Chat.ID,
		"user_id", update.Message.From.ID,
	)

	params := bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	}
	if _, err := tg.SendMessage(ctx, &params); err != nil {
		logger.Error(err,
			"unable to send message",
			"chat_id", update.Message.Chat.ID,
			"text", update.Message.Text,
		)
	}
}

package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/clients/telegram"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Bot struct {
	cfg *config.Config
	tg  *telegram.Client
}

func NewBot(cfg *config.Config) *Bot {
	return &Bot{
		cfg: cfg,
		tg:  telegram.NewClient(),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	opts := []gotelegram.Option{
		gotelegram.WithDefaultHandler(b.handleUpdate),
	}
	if b.cfg.Bot.CheckIfUserAllowed {
		opts = append(opts, gotelegram.WithMiddlewares(b.isUserAllowed))
	}

	tg, err := gotelegram.New(b.cfg.Bot.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	logger.Info("starting bot")
	tg.Start(ctx)
	logger.Info("bot stopped")
	return nil
}

func (b *Bot) handleUpdate(ctx context.Context, tg *gotelegram.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("handleUpdate")
	if update.Message == nil {
		return
	}
	logger.Info(
		"received message",
		"chat_id", update.Message.Chat.ID,
		"user_id", update.Message.From.ID,
	)

	if err := b.tg.SendMessage(ctx, tg, update.Message.Chat.ID, update.Message.Text); err != nil {
		logger.Error(err,
			"unable to send message",
			"chat_id", update.Message.Chat.ID,
			"text", update.Message.Text,
		)
	}
}

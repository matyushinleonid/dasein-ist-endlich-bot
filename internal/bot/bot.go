package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DaseinBot struct {
	cfg *config.Config
}

func NewBot(cfg *config.Config) *DaseinBot {
	return &DaseinBot{
		cfg: cfg,
	}
}

func (daseinBot *DaseinBot) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	opts := []gotelegram.Option{
		gotelegram.WithDefaultHandler(daseinBot.handleUpdate),
	}
	if daseinBot.cfg.CheckIfUserAllowed {
		opts = append(opts, gotelegram.WithMiddlewares(daseinBot.isUserAllowed))
	}

	bot, err := gotelegram.New(daseinBot.cfg.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	logger.Info("starting bot")
	bot.Start(ctx)
	logger.Info("bot stopped")
	return nil
}

func (daseinBot *DaseinBot) handleUpdate(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("handleUpdate")
	if update.Message == nil {
		return
	}
	logger.Info(
		"received message",
		"chat_id", update.Message.Chat.ID,
		"user_id", update.Message.From.ID,
	)

	if err := utils.SendMessage(ctx, bot, update.Message.Chat.ID, update.Message.Text); err != nil {
		logger.Error(err,
			"unable to send message",
			"chat_id", update.Message.Chat.ID,
			"text", update.Message.Text,
		)
	}
}

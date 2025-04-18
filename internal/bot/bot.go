package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
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
		gotelegram.WithDefaultHandler(daseinBot.answerHandler),
	}
	if daseinBot.cfg.CheckIfUserAllowed {
		opts = append(opts, gotelegram.WithMiddlewares(daseinBot.isUserAllowed))
	}
	bot, err := gotelegram.New(daseinBot.cfg.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %w", err)
	}
	bot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/askme", gotelegram.MatchTypeExact, daseinBot.askMeHandler)

	logger.Info("starting bot")
	bot.Start(ctx)
	logger.Info("bot stopped")
	return nil
}

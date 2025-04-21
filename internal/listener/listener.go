package listener

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/handler"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/middleware"
)

type Listener struct {
	*core.DaseinBot
}

func New(daseinBot *core.DaseinBot) *Listener {
	return &Listener{
		DaseinBot: daseinBot,
	}
}

func (b *Listener) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	opts := []gotelegram.Option{
		gotelegram.WithDefaultHandler(handler.AnswerHandler(b.DaseinBot)),
	}
	if b.Cfg.CheckIfUserAllowed {
		opts = append(opts, gotelegram.WithMiddlewares(middleware.IsUpdateLegal(b.DaseinBot), middleware.IsUserAllowed(b.DaseinBot)))
	}
	tgbot, err := gotelegram.New(b.Cfg.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram tgbot: %w", err)
	}

	tgbot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/start", gotelegram.MatchTypeExact, handler.StartHandler(b.DaseinBot))
	tgbot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/about", gotelegram.MatchTypeExact, handler.AboutHandler(b.DaseinBot))
	tgbot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/help", gotelegram.MatchTypeExact, handler.HelpHandler(b.DaseinBot))
	tgbot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/begin", gotelegram.MatchTypeExact, handler.BeginHandler(b.DaseinBot))

	logger.Info("starting tgbot")
	tgbot.Start(ctx)
	logger.Info("tgbot stopped")
	return nil
}

package notifier

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

type Notifier struct {
	*core.DaseinBot
}

func New(daseinBot *core.DaseinBot) *Notifier {
	return &Notifier{DaseinBot: daseinBot}
}

func (n *Notifier) NotifyAll(ctx context.Context, tgbot *gotelegram.Bot) error {
	users, err := n.UserRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("FindAll failed: %w", err)
	}

	now := time.Now()
	for _, u := range users {
		ok, err := ShouldNotify(u, now)
		if err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "cannot decide to notify", "user_id", u.ID)
			continue
		}
		if !ok {
			continue
		}

		msg := FormatNotificationMessage(n.Cfg.DaysLeftMessage, u, now)
		if err := n.TelegramClient.SendMessage(ctx, tgbot, u.ID, msg); err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "send failed", "user_id", u.ID)
			continue
		}

		u.LastNotification = now
		if err := n.UserRepository.Update(ctx, &u); err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "update user failed", "user_id", u.ID)
		}
	}

	return nil
}

func (n *Notifier) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	tgbot, err := gotelegram.New(n.Cfg.Token)
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	if err := n.NotifyAll(ctx, tgbot); err != nil {
		return err
	}
	logger.Info("notifier finished")
	return nil
}

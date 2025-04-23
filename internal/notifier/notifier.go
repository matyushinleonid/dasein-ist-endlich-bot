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
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type Notifier struct {
	*core.DaseinBot
}

func New(daseinBot *core.DaseinBot) *Notifier {
	return &Notifier{DaseinBot: daseinBot}
}

func (n *Notifier) NotifyAll(ctx context.Context, tgbot *gotelegram.Bot) error {
	cur, err := n.MongoClient.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("FindAll failed: %w", err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "close cursor failed")
		}
	}()

	now := time.Now()
	for cur.Next(ctx) {
		var u model.User
		if err := cur.Decode(&u); err != nil {
			return fmt.Errorf("decode user: %w", err)
		}

		ok, err := ShouldNotify(u, now)
		if err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "cannot decide to notify", "user_id", u.ID)
			continue
		}
		if !ok {
			continue
		}

		msg := FormatNotificationMessage(u, now)
		if err := n.TelegramClient.SendMessage(ctx, tgbot, u.ID, msg); err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "send failed", "user_id", u.ID)
			continue
		}

		if _, err := n.MongoClient.Update(ctx, u.ID, bson.M{"last_notification": now}); err != nil {
			logr.FromContextOrDiscard(ctx).
				Error(err, "update user failed", "user_id", u.ID)
		}
	}

	if err := cur.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
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

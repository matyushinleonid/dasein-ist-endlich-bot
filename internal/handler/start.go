package handler

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/record"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/utils"
)

func StartHandler(b *bot.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("startHandler").WithValues("chat_id", update.Message.Chat.ID)

		if update.Message == nil {
			return
		}

		chatID := update.Message.Chat.ID
		var rec record.Record
		err := b.MongoClient.Get(ctx, chatID, &rec)
		if err != nil {
			rec = record.Record{ID: chatID, DaysLeft: 0, Calculated: false}
			if _, err = b.MongoClient.Create(ctx, rec); err != nil {
				logger.Error(err, "failed to create record")
			}
		}

		if err = utils.SendMessage(ctx, tgbot, update.Message.Chat.ID, b.Cfg.Start); err != nil {
			logger.Error(err, "unable to send message")
		}
	}
}

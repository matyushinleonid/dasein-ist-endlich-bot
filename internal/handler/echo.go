package handler

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func EchoHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("defaultHandler").WithValues("chat_id", update.Message.Chat.ID)

		if update.Message == nil {
			return
		}
		if err := b.TelegramClient.SendMessage(ctx, tgbot, update.Message.Chat.ID, update.Message.Text); err != nil {
			logger.Error(err, "unable to send message")
		}
	}
}

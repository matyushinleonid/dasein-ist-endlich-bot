package middleware

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func IsUserAllowed(b *core.DaseinBot) gotelegram.Middleware {
	return func(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
		return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
			logger := logr.FromContextOrDiscard(ctx).WithName("isUserAllowed").WithValues("chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)
			userID := update.Message.From.ID
			allowed := false
			for _, allowedUser := range b.Cfg.AllowedUsers {
				if userID == allowedUser {
					allowed = true
					break
				}
			}
			if !allowed {
				logger.Info("user not allowed", "user_id", userID)
				if err := b.TelegramClient.SendMessage(ctx, tgbot, update.Message.Chat.ID, "Get lost!"); err != nil {
					logger.Error(err, "unable to send message")
				}
				return
			}
			next(ctx, tgbot, update)
		}
	}
}

package middleware

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/utils"
)

func IsUserAllowed(b *bot.DaseinBot) gotelegram.Middleware {
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
				if err := utils.SendMessage(ctx, tgbot, update.Message.Chat.ID, "Get lost!"); err != nil {
					logger.Error(err, "unable to send message")
				}
				return
			}
			next(ctx, tgbot, update)
		}
	}
}

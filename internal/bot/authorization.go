package bot

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) isUserAllowed(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
	return func(ctx context.Context, tg *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("isUserAllowed")
		userID := update.Message.From.ID
		allowed := false
		for _, allowedUser := range b.cfg.Bot.AllowedUsers {
			if userID == allowedUser {
				allowed = true
				break
			}
		}
		if !allowed {
			logger.Info("user not allowed", "user_id", userID)
			if err := b.tg.SendMessage(ctx, tg, update.Message.Chat.ID, "Get lost!"); err != nil {
				logger.Error(err,
					"unable to send message",
					"chat_id", update.Message.Chat.ID,
					"text", update.Message.Text,
				)
			}
			return
		}
		next(ctx, tg, update)
	}
}

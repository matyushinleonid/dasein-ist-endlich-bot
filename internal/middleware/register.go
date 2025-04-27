package middleware

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func RegisterUser(b *core.DaseinBot) gotelegram.Middleware {
	return func(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
		return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
			if update.CallbackQuery != nil {
				next(ctx, tgbot, update)
				return
			}
			logger := logr.FromContextOrDiscard(ctx).WithName("isUserAllowed").WithValues("chat_id", update.Message.Chat.ID, "user_id", update.Message.From.ID)
			chatId := update.Message.Chat.ID
			exists, err := b.UserRepository.UserExists(ctx, chatId)
			if err != nil {
				logger.Error(err, "failed to check if user exists")
				return
			}
			if exists {
				next(ctx, tgbot, update)
				return
			}
			if _, err = b.UserRepository.Create(ctx, chatId); err != nil {
				logger.Error(err, "failed to create user")
				return
			}
			if err != nil {
				logger.Error(err, "failed to create user")
			}
			next(ctx, tgbot, update)
		}
	}
}

package middleware

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func IsUpdateLegal(b *core.DaseinBot) gotelegram.Middleware {
	return func(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
		return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
			logger := logr.FromContextOrDiscard(ctx).WithName("isUpdateLegal")
			if update.Message == nil {
				logger.Info("message is empty")
				return
			}
			if update.Message.Text == "" {
				logger.Info("message text is empty")
				return
			}
			next(ctx, tgbot, update)
		}
	}
}

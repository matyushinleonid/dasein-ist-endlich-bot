package bot

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (daseinBot *DaseinBot) isUpdateLegal(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
	return func(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("isUpdateLegal")
		if update.Message == nil {
			logger.Info("message is empty")
			return
		}
		if update.Message.Text == "" {
			logger.Info("message text is empty")
			return
		}
		next(ctx, bot, update)
	}
}

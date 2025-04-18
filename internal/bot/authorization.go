package bot

import (
	"context"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
)

func (daseinBot *DaseinBot) isUserAllowed(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
	return func(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("isUserAllowed")
		userID := update.Message.From.ID
		allowed := false
		for _, allowedUser := range daseinBot.cfg.AllowedUsers {
			if userID == allowedUser {
				allowed = true
				break
			}
		}
		if !allowed {
			logger.Info("user not allowed", "user_id", userID)
			if err := utils.SendMessage(ctx, bot, update.Message.Chat.ID, "Get lost!"); err != nil {
				logger.Error(err,
					"unable to send message",
					"chat_id", update.Message.Chat.ID,
					"text", update.Message.Text,
				)
			}
			return
		}
		next(ctx, bot, update)
	}
}

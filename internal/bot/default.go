package bot

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
)

func (daseinBot *DaseinBot) defaultHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("defaultHandler").WithValues("chat_id", update.Message.Chat.ID)

	if update.Message == nil {
		return
	}
	if err := utils.SendMessage(ctx, bot, update.Message.Chat.ID, update.Message.Text); err != nil {
		logger.Error(err, "unable to send message")
	}
}

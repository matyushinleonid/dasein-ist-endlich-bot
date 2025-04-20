package bot

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
)

func (daseinBot *DaseinBot) aboutHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("aboutHandler")

	if update.Message == nil {
		return
	}
	if err := utils.SendMessage(ctx, bot, update.Message.Chat.ID, daseinBot.cfg.About); err != nil {
		logger.Error(err,
			"unable to send message",
			"chat_id", update.Message.Chat.ID,
			"text", daseinBot.cfg.About,
		)
	}
}

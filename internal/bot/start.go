package bot

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
)

func (daseinBot *DaseinBot) startHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("startHandler").WithValues("chat_id", update.Message.Chat.ID)

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	var rec Record
	err := daseinBot.mongoClient.Get(ctx, chatID, &rec)
	if err != nil {
		rec = Record{ID: chatID, DaysLeft: 0, Calculated: false}
		if _, err = daseinBot.mongoClient.Create(ctx, rec); err != nil {
			logger.Error(err, "failed to create record")
		}
	}

	if err = utils.SendMessage(ctx, bot, update.Message.Chat.ID, daseinBot.cfg.Start); err != nil {
		logger.Error(err, "unable to send message")
	}
}

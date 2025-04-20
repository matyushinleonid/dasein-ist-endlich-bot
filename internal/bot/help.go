package bot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
)

func (daseinBot *DaseinBot) helpHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("helpHandler").WithValues("chat_id", update.Message.Chat.ID)

	msg := daseinBot.cfg.Help
	if daseinBot.cfg.Debug {
		chatID := update.Message.Chat.ID
		var rec Record
		err := daseinBot.mongoClient.Get(ctx, chatID, &rec)
		if err != nil {
			logger.Error(err, "failed to get record")
		}
		prettyJSON, err := json.MarshalIndent(rec, "", "\t")
		if err != nil {
			logger.Error(err, "failed to marshal record")
		}
		msg += fmt.Sprintf("\n\nDebug info:\n\tRecord:\n%s\n", prettyJSON)
	}

	err := utils.SendMessage(ctx, bot, update.Message.Chat.ID, msg)
	if err != nil {
		logger.Error(err, "unable to send message")
	}
}

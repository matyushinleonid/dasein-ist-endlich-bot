package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

func HelpHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("helpHandler").WithValues("chat_id", update.Message.Chat.ID)

		msg := b.Cfg.Help
		if b.Cfg.Debug {
			chatID := update.Message.Chat.ID
			var rec model.User
			err := b.MongoClient.Get(ctx, chatID, &rec)
			if err != nil {
				logger.Error(err, "failed to get record")
			}
			prettyJSON, err := json.MarshalIndent(rec, "", "\t")
			if err != nil {
				logger.Error(err, "failed to marshal record")
			}
			msg += fmt.Sprintf("\n\nDebug info:\n\tRecord:\n%s\n", prettyJSON)
		}

		err := b.TelegramClient.SendMessage(ctx, tgbot, update.Message.Chat.ID, msg)
		if err != nil {
			logger.Error(err, "unable to send message")
		}
	}
}

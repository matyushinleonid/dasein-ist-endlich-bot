package handler

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

const NfCallbackPrefix = "nf_"

func NotificationFrequencyHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		kb := buildKeyboard(b, chatID)
		tgbot.SendMessage(ctx, &gotelegram.SendMessageParams{
			ChatID:      chatID,
			Text:        "Choose your notification frequency",
			ReplyMarkup: kb,
		})
	}
}

func NotificationFrequencyCallbackHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}
		data := update.CallbackQuery.Data
		if !strings.HasPrefix(data, NfCallbackPrefix) {
			return
		}

		err := b.TelegramClient.AnswerCallbackQuery(ctx, tgbot, update.CallbackQuery.ID)
		if err != nil {
			return
		}

		choice := strings.TrimPrefix(data, NfCallbackPrefix)
		msg := update.CallbackQuery.Message.Message
		if msg == nil {
			return
		}
		chatID := msg.Chat.ID

		var freq model.NotificationFrequency
		switch choice {
		case string(model.Daily):
			freq = model.Daily
		case string(model.Weekly):
			freq = model.Weekly
		case string(model.Monthly):
			freq = model.Monthly
		case string(model.Yearly):
			freq = model.Yearly
		case string(model.Never):
			freq = model.Never
		default:
			return
		}

		user, err := b.UserRepository.Get(ctx, chatID)
		if err != nil {
			logger := logr.FromContextOrDiscard(ctx)
			logger.Error(err, "unable to get user from DB")
			return
		}
		user.NotificationFrequency = freq
		if err := b.UserRepository.Update(ctx, user); err != nil {
			logger := logr.FromContextOrDiscard(ctx)
			logger.Error(err, "unable to update user in DB")
			return
		}

		kb := buildKeyboard(b, chatID)
		err = b.TelegramClient.EditMessageKeyboard(ctx, tgbot, chatID, msg.ID, kb)
		if err != nil {
			return
		}
	}
}

func buildKeyboard(b *core.DaseinBot, chatID int64) *models.InlineKeyboardMarkup {
	user, _ := b.UserRepository.Get(context.Background(), chatID)
	current := user.NotificationFrequency

	freqs := []model.NotificationFrequency{model.Daily, model.Weekly, model.Monthly, model.Yearly, model.Never}
	buttons := make([]models.InlineKeyboardButton, 0, len(freqs))
	for _, f := range freqs {
		text := string(f)
		if f == current {
			text = "✅ " + text
		}
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         text,
			CallbackData: NfCallbackPrefix + string(f),
		})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			buttons[0:2],
			buttons[2:4],
			{buttons[4]},
		},
	}
}

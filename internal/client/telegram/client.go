package telegram

import (
	"context"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Client interface {
	SendMessage(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string) error
	SendMessageKeyboard(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string, replyMarkup models.ReplyMarkup) error
	EditMessageKeyboard(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, messageID int, replyMarkup models.ReplyMarkup) error
	SendTypingAction(ctx context.Context, tgbot *gotelegram.Bot, chatID int64) error
	AnswerCallbackQuery(ctx context.Context, bot *gotelegram.Bot, callbackQueryID string) error
}

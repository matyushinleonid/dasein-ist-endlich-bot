package telegram

import (
	"context"

	gotelegram "github.com/go-telegram/bot"
)

type Client interface {
	SendMessage(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string) error
	SendTypingAction(ctx context.Context, tgbot *gotelegram.Bot, chatID int64) error
}

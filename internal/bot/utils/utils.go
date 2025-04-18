package utils

import (
	"context"
	"fmt"

	gotelegram "github.com/go-telegram/bot"
)

func SendMessage(ctx context.Context, bot *gotelegram.Bot, chatID int64, text string) error {
	params := gotelegram.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}
	if _, err := bot.SendMessage(ctx, &params); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

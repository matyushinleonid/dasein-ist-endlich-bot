package utils

import (
	"context"
	"fmt"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func SendTypingAction(ctx context.Context, bot *gotelegram.Bot, chatID int64) error {
	params := gotelegram.SendChatActionParams{
		ChatID: chatID,
		Action: models.ChatActionTyping,
	}
	if _, err := bot.SendChatAction(ctx, &params); err != nil {
		return fmt.Errorf("failed to send typing action: %w", err)
	}
	return nil
}

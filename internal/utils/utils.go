package utils

import (
	"context"
	"fmt"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendMessage(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string) error {
	params := gotelegram.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}
	if _, err := tgbot.SendMessage(ctx, &params); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func SendTypingAction(ctx context.Context, tgbot *gotelegram.Bot, chatID int64) error {
	params := gotelegram.SendChatActionParams{
		ChatID: chatID,
		Action: models.ChatActionTyping,
	}
	if _, err := tgbot.SendChatAction(ctx, &params); err != nil {
		return fmt.Errorf("failed to send typing action: %w", err)
	}
	return nil
}

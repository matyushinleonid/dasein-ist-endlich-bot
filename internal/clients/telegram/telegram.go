package telegram

import (
	"context"
	"fmt"

	gotelegram "github.com/go-telegram/bot"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (tg *Client) SendMessage(ctx context.Context, bot *gotelegram.Bot, chatID int64, text string) error {
	params := gotelegram.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}
	if _, err := bot.SendMessage(ctx, &params); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

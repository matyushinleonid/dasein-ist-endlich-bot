package telegram

import (
	"context"

	gotelegram "github.com/go-telegram/bot"
)

type DummyClient struct {
	SentMessages []struct {
		ChatID int64
		Text   string
	}
	SentTypingChats []int64
}

func NewDummyClient() *DummyClient {
	return &DummyClient{}
}

func (d *DummyClient) SendMessage(ctx context.Context, bot *gotelegram.Bot, chatID int64, text string) error {
	d.SentMessages = append(d.SentMessages, struct {
		ChatID int64
		Text   string
	}{chatID, text})
	return nil
}

func (d *DummyClient) SendTypingAction(ctx context.Context, bot *gotelegram.Bot, chatID int64) error {
	d.SentTypingChats = append(d.SentTypingChats, chatID)
	return nil
}

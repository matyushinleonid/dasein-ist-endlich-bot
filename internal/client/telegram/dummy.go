package telegram

import (
	"context"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DummyClient struct {
	SentMessages []struct {
		ChatID int64
		Text   string
	}
	SentTypingChats       []int64
	ErrOnSendMessage      error
	ErrOnSendTypingAction error
}

func NewDummyClient() *DummyClient {
	return &DummyClient{}
}

func (d *DummyClient) SendMessage(ctx context.Context, bot *gotelegram.Bot, chatID int64, text string) error {
	if d.ErrOnSendMessage != nil {
		return d.ErrOnSendMessage
	}
	d.SentMessages = append(d.SentMessages, struct {
		ChatID int64
		Text   string
	}{chatID, text})
	return nil
}

func (d *DummyClient) SendMessageKeyboard(ctx context.Context, bot *gotelegram.Bot, chatID int64, text string, replyMarkup models.ReplyMarkup) error {
	return d.SendMessage(ctx, bot, chatID, text)
}

func (d *DummyClient) EditMessageKeyboard(ctx context.Context, bot *gotelegram.Bot, chatID int64, messageID int, replyMarkup models.ReplyMarkup) error {
	if d.ErrOnSendMessage != nil {
		return d.ErrOnSendMessage
	}
	d.SentMessages = append(d.SentMessages, struct {
		ChatID int64
		Text   string
	}{chatID, ""})
	return nil
}

func (d *DummyClient) SendTypingAction(ctx context.Context, bot *gotelegram.Bot, chatID int64) error {
	if d.ErrOnSendTypingAction != nil {
		return d.ErrOnSendTypingAction
	}
	d.SentTypingChats = append(d.SentTypingChats, chatID)
	return nil
}

func (d *DummyClient) AnswerCallbackQuery(ctx context.Context, bot *gotelegram.Bot, callbackQueryID string) error {
	return nil
}

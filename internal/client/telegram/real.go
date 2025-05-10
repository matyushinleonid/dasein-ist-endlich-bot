package telegram

import (
	"context"
	"fmt"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/retry"
)

type RealClient struct{}

func NewRealClient() *RealClient {
	return &RealClient{}
}

func (c *RealClient) sendMessage(ctx context.Context, tgbot *gotelegram.Bot, parseMode models.ParseMode, chatID int64, text string) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		_, e := tgbot.SendMessage(ctx, &gotelegram.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: parseMode,
		})
		return e
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (c *RealClient) SendMessage(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string) error {
	return c.sendMessage(ctx, tgbot, "", chatID, text)
}

func (c *RealClient) SendMessageMarkdown(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string) error {
	return c.sendMessage(ctx, tgbot, models.ParseModeMarkdownV1, chatID, text)
}

func (c *RealClient) SendMessageKeyboard(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, text string, replyMarkup models.ReplyMarkup) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		_, e := tgbot.SendMessage(ctx, &gotelegram.SendMessageParams{
			ChatID:      chatID,
			Text:        text,
			ReplyMarkup: replyMarkup,
		})
		return e
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (c *RealClient) EditMessageKeyboard(ctx context.Context, tgbot *gotelegram.Bot, chatID int64, messageID int, replyMarkup models.ReplyMarkup) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		_, e := tgbot.EditMessageReplyMarkup(ctx, &gotelegram.EditMessageReplyMarkupParams{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: replyMarkup,
		})
		return e
	})
	if err != nil {
		return fmt.Errorf("failed to edit message keyboard: %w", err)
	}
	return nil
}

func (c *RealClient) SendTypingAction(ctx context.Context, tgbot *gotelegram.Bot, chatID int64) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		_, e := tgbot.SendChatAction(ctx, &gotelegram.SendChatActionParams{
			ChatID: chatID,
			Action: models.ChatActionTyping,
		})
		return e
	})
	if err != nil {
		return fmt.Errorf("failed to send typing action: %w", err)
	}
	return nil
}

func (c *RealClient) AnswerCallbackQuery(ctx context.Context, bot *gotelegram.Bot, callbackQueryID string) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		_, e := bot.AnswerCallbackQuery(ctx, &gotelegram.AnswerCallbackQueryParams{
			CallbackQueryID: callbackQueryID,
			ShowAlert:       false,
		})
		return e
	})
	if err != nil {
		return fmt.Errorf("failed to answer callback query: %w", err)
	}
	return nil
}

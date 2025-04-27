package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

func BeginHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).
			WithName("beginHandler").
			WithValues("chat_id", update.Message.Chat.ID)

		chatID := update.Message.Chat.ID

		sess := model.Session{
			Stage:   0,
			Answers: make([]string, len(b.Cfg.Questions)),
		}
		if err := b.SessionRepository.Save(ctx, chatID, &sess); err != nil {
			logger.Error(err, "failed to save session to Redis")
			err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "Redis is not available, please try again later.")
			if err != nil {
				logger.Error(err, "failed to send error message")
			}
			return
		}

		question := b.Cfg.Questions[0]
		if err := b.TelegramClient.SendMessage(ctx, tgbot, chatID, question); err != nil {
			logger.Error(err, "unable to send first question")
		}
	}
}

func AnswerHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).
			WithName("answerHandler").
			WithValues("chat_id", update.Message.Chat.ID)

		chatID := update.Message.Chat.ID

		sess, err := b.SessionRepository.Get(ctx, chatID)
		if err != nil {
			if errors.Is(err, repository.ErrSessionNotFound) {
				HelpHandler(b)(ctx, tgbot, update)
				return
			}
			logger.Error(err, "failed to get session from Redis")
			err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "Redis is not available, please try again later.")
			if err != nil {
				logger.Error(err, "failed to send error message")
			}
			return
		}

		sess.Answers[sess.Stage] = update.Message.Text
		sess.Stage++

		if sess.Stage < len(b.Cfg.Questions) {
			if err = b.SessionRepository.Save(ctx, chatID, sess); err != nil {
				logger.Error(err, "failed to update session in Redis")
				err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "Redis is not available, please try again later.")
				if err != nil {
					logger.Error(err, "failed to send error message")
				}
				return
			}
			nextQ := b.Cfg.Questions[sess.Stage]
			if err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, nextQ); err != nil {
				logger.Error(err, "unable to send next question")
			}
			return
		}

		if err = b.TelegramClient.SendTypingAction(ctx, tgbot, chatID); err != nil {
			logger.Error(err, "unable to send typing action")
		}

		summary := ""
		for i, ans := range sess.Answers {
			summary += fmt.Sprintf("%d) %s — %s\n", i+1, b.Cfg.Questions[i], ans)
		}

		var response model.OpenAIResponse
		if err = b.OpenAIClient.SendJSONUnmarshal(
			ctx, chatID, summary,
			"OpenAIResponseSchema", model.OpenAIResponseSchema,
			&response,
		); err != nil {
			logger.Error(err, "unable to query OpenAI")
			err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "OpenAI is not available, please try again later.")
			if err != nil {
				logger.Error(err, "failed to send error message")
			}
			return
		}

		user, err := b.UserRepository.Get(ctx, chatID)
		if err != nil {
			logger.Error(err, "unable to get user from DB")
			err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "Database is not available, please try again later.")
			if err != nil {
				logger.Error(err, "failed to send error message")
			}
			return
		}
		user.DeathTime = model.DeathTime(time.Now(), response.DaysLeft)
		user.LastNotification = time.Time{}
		user.Calculated = true
		user.NotificationFrequency = model.Daily
		if err = b.UserRepository.Update(ctx, user); err != nil {
			logger.Error(err, "unable to update user in DB")
			err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, "Database is not available, please try again later.")
			if err != nil {
				logger.Error(err, "failed to send error message")
			}
			return
		}

		respText := fmt.Sprintf(
			"Days left in this world: %d\n\n%s",
			response.DaysLeft,
			response.Description,
		)
		if err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, respText); err != nil {
			logger.Error(err, "unable to send final message")
			return
		}

		if err = b.SessionRepository.Delete(ctx, chatID); err != nil {
			logger.Error(err, "failed to delete session from Redis")
			return
		}
	}
}

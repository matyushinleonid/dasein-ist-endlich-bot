package handler

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

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
		if err := b.RedisClient.Save(ctx, chatID, &sess); err != nil {
			logger.Error(err, "failed to save session to Redis")
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

		var sess model.Session
		err := b.RedisClient.Load(ctx, chatID, &sess)
		if err != nil {

			EchoHandler(b)(ctx, tgbot, update)
			return
		}

		sess.Answers[sess.Stage] = update.Message.Text
		sess.Stage++

		if sess.Stage < len(b.Cfg.Questions) {
			if err = b.RedisClient.Save(ctx, chatID, &sess); err != nil {
				logger.Error(err, "failed to update session in Redis")
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
		}

		upd := model.User{DaysLeft: response.DaysLeft, Calculated: true}
		if _, err := b.MongoClient.Update(ctx, chatID, upd); err != nil {
			logger.Error(err, "failed to update record in MongoDB")
		}

		respText := fmt.Sprintf(
			"У вас осталось %d дней в этом мире.\n\n%s",
			response.DaysLeft,
			response.Description,
		)
		if err = b.TelegramClient.SendMessage(ctx, tgbot, chatID, respText); err != nil {
			logger.Error(err, "unable to send final message")
		}

		if err = b.RedisClient.Delete(ctx, chatID); err != nil {
			logger.Error(err, "failed to delete session from Redis")
		}
	}
}

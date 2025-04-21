package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/record"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/utils"
)

type Conversation struct {
	Stage   int
	Answers []string
}

type openAIResponse struct {
	DaysLeft    int64  `json:"days_left"`
	Description string `json:"description"`
}

var openAIResponseSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"days_left": map[string]interface{}{
			"type":        "integer",
			"description": "How many days the user have left in this world",
		},
		"description": map[string]interface{}{
			"type":        "string",
			"description": "How the amount of days left was calculated based on the user's answers",
		},
	},
	"required":             []interface{}{"days_left", "description"},
	"additionalProperties": false,
}

var (
	convs   = make(map[int64]*Conversation)
	convsMu sync.Mutex
)

func BeginHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("beginHandler").WithValues("chat_id", update.Message.Chat.ID)

		chatID := update.Message.Chat.ID
		convsMu.Lock()
		convs[chatID] = &Conversation{
			Stage:   0,
			Answers: make([]string, len(b.Cfg.Questions)),
		}
		convsMu.Unlock()

		err := utils.SendMessage(ctx, tgbot, chatID, b.Cfg.Questions[0])
		if err != nil {
			logger.Error(err, "unable to send message")
		}
	}
}

func AnswerHandler(b *core.DaseinBot) gotelegram.HandlerFunc {
	return func(ctx context.Context, tgbot *gotelegram.Bot, update *models.Update) {
		logger := logr.FromContextOrDiscard(ctx).WithName("answerHandler").WithValues("chat_id", update.Message.Chat.ID)

		chatID := update.Message.Chat.ID

		convsMu.Lock()
		conv, ok := convs[chatID]
		convsMu.Unlock()
		if !ok {
			EchoHandler(b)(ctx, tgbot, update)
			return
		}

		conv.Answers[conv.Stage] = update.Message.Text
		conv.Stage++

		if conv.Stage < len(b.Cfg.Questions) {
			err := utils.SendMessage(ctx, tgbot, chatID, b.Cfg.Questions[conv.Stage])
			if err != nil {
				logger.Error(err, "unable to send message")
			}
			return
		}

		if err := utils.SendTypingAction(ctx, tgbot, chatID); err != nil {
			logger.Error(err, "unable to send typing action")
		}

		summary := ""
		for i, ans := range conv.Answers {
			summary += fmt.Sprintf("%d) %s — %s\n", i+1, b.Cfg.Questions[i], ans)
		}

		var response openAIResponse
		err := b.OpenAIClient.SendJSONUnmarshal(ctx, chatID, summary, "openAIResponseSchema", openAIResponseSchema, &response)
		if err != nil {
			logger.Error(err, "unable to query OpenAI")
		}

		upd := record.Record{DaysLeft: response.DaysLeft, Calculated: true}
		if _, err = b.MongoClient.Update(ctx, chatID, upd); err != nil {
			logger.Error(err, "failed to update record")
		}

		responseText := fmt.Sprintf("У вас осталось %d дней в этом мире.\n\n%s", response.DaysLeft, response.Description)
		err = utils.SendMessage(ctx, tgbot, chatID, responseText)

		convsMu.Lock()
		delete(convs, chatID)
		convsMu.Unlock()
	}
}

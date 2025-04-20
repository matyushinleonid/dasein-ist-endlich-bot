package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/bot/utils"
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

func (daseinBot *DaseinBot) askMeHandler(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("askMeHandler")

	chatID := update.Message.Chat.ID
	convsMu.Lock()
	convs[chatID] = &Conversation{
		Stage:   0,
		Answers: make([]string, len(daseinBot.cfg.Questions)),
	}
	convsMu.Unlock()

	err := utils.SendMessage(ctx, bot, chatID, daseinBot.cfg.Questions[0])
	if err != nil {
		logger.Error(err,
			"unable to send message",
			"chat_id", chatID,
			"text", update.Message.Text,
		)
	}
}

func (daseinBot *DaseinBot) answerHandler(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
	logger := logr.FromContextOrDiscard(ctx).WithName("answerHandler")

	chatID := update.Message.Chat.ID

	convsMu.Lock()
	conv, ok := convs[chatID]
	convsMu.Unlock()
	if !ok {
		daseinBot.defaultHandler(ctx, bot, update)
		return
	}

	conv.Answers[conv.Stage] = update.Message.Text
	conv.Stage++

	if conv.Stage < len(daseinBot.cfg.Questions) {
		err := utils.SendMessage(ctx, bot, chatID, daseinBot.cfg.Questions[conv.Stage])
		if err != nil {
			logger.Error(err,
				"unable to send message",
				"chat_id", chatID,
				"text", update.Message.Text,
			)
		}
		return
	}

	if err := utils.SendTypingAction(ctx, bot, chatID); err != nil {
		logger.Error(err,
			"unable to send typing action",
			"chat_id", chatID,
		)
	}

	summary := ""
	for i, ans := range conv.Answers {
		summary += fmt.Sprintf("%d) %s — %s\n", i+1, daseinBot.cfg.Questions[i], ans)
	}

	var response openAIResponse
	err := daseinBot.openAIClient.SendJSONUnmarshal(ctx, chatID, summary, "openAIResponseSchema", openAIResponseSchema, &response)
	if err != nil {
		logger.Error(err,
			"unable to query OpenAI",
			"chat_id", chatID,
		)
	}

	responseText := fmt.Sprintf("У вас осталось %d дней в этом мире.\n\n%s", response.DaysLeft, response.Description)
	err = utils.SendMessage(ctx, bot, chatID, responseText)

	convsMu.Lock()
	delete(convs, chatID)
	convsMu.Unlock()
}

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
			"text", update.Message.Text,
		)
	}

	summary := ""
	for i, ans := range conv.Answers {
		summary += fmt.Sprintf("%d) %s — %s\n", i+1, daseinBot.cfg.Questions[i], ans)
	}
	response, err := daseinBot.openAIClient.SendText(ctx, chatID, summary)
	if err != nil {
		logger.Error(err,
			"unable to query OpenAI",
			"chat_id", chatID,
			"text", update.Message.Text,
		)
		return
	}
	err = utils.SendMessage(ctx, bot, chatID, response)

	convsMu.Lock()
	delete(convs, chatID)
	convsMu.Unlock()
}

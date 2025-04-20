package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/openai"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DaseinBot struct {
	cfg          *config.DaseinBotConfig
	openAIClient openai.Client
	mongoClient  mongo.Client
}

func NewBot(cfg *config.Config) *DaseinBot {
	mcli, err := mongo.NewClient(cfg.MongoDB)
	if err != nil {
		panic("failed to init mongo: " + err.Error())
	}
	return &DaseinBot{
		cfg:          &cfg.DaseinBot,
		openAIClient: openai.NewClient(cfg.OpenAI),
		mongoClient:  mcli,
	}
}

func (daseinBot *DaseinBot) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	logger := logr.FromContextOrDiscard(ctx)

	opts := []gotelegram.Option{
		gotelegram.WithDefaultHandler(daseinBot.answerHandler),
	}
	if daseinBot.cfg.CheckIfUserAllowed {
		opts = append(opts, gotelegram.WithMiddlewares(daseinBot.isUpdateLegal, daseinBot.isUserAllowed))
	}
	bot, err := gotelegram.New(daseinBot.cfg.Token, opts...)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	bot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/start", gotelegram.MatchTypeExact, daseinBot.startHandler)
	bot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/about", gotelegram.MatchTypeExact, daseinBot.aboutHandler)
	bot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/help", gotelegram.MatchTypeExact, daseinBot.helpHandler)
	bot.RegisterHandler(gotelegram.HandlerTypeMessageText, "/begin", gotelegram.MatchTypeExact, daseinBot.beginHandler)

	logger.Info("starting bot")
	bot.Start(ctx)
	logger.Info("bot stopped")
	return nil
}

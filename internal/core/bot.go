package core

import (
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/openai"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/redis"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/telegram"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DaseinBot struct {
	Cfg               *config.DaseinBotConfig
	OpenAIClient      openai.Client
	UserRepository    *repository.UserRepository
	TelegramClient    telegram.Client
	SessionRepository *repository.SessionRepository
}

func New(cfg *config.Config) *DaseinBot {
	mcli, err := mongo.NewRealClient(cfg.MongoDB)
	if err != nil {
		panic("failed to init mongo: " + err.Error())
	}
	return &DaseinBot{
		Cfg:               &cfg.DaseinBot,
		OpenAIClient:      openai.NewClient(cfg.OpenAI),
		UserRepository:    repository.NewUserRepository(mcli),
		TelegramClient:    telegram.NewRealClient(),
		SessionRepository: repository.NewSessionRepository(redis.NewRealClient(cfg.Redis)),
	}
}

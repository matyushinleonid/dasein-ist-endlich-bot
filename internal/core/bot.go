package core

import (
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/openai"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DaseinBot struct {
	Cfg          *config.DaseinBotConfig
	OpenAIClient openai.Client
	MongoClient  mongo.Client
}

func New(cfg *config.Config) *DaseinBot {
	mcli, err := mongo.NewClient(cfg.MongoDB)
	if err != nil {
		panic("failed to init mongo: " + err.Error())
	}
	return &DaseinBot{
		Cfg:          &cfg.DaseinBot,
		OpenAIClient: openai.NewClient(cfg.OpenAI),
		MongoClient:  mcli,
	}
}

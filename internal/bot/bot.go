package bot

import (
	"context"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/openai"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type DaseinBotInterface interface {
	Run(ctx context.Context) error
}

type DaseinBot struct {
	Cfg          *config.DaseinBotConfig
	OpenAIClient openai.Client
	MongoClient  mongo.Client
}

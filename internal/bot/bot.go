package bot

import (
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
)

type Bot struct {
	config *config.Config
}

func NewBot(config *config.Config) *Bot {
	return &Bot{
		config: config,
	}
}

func (b *Bot) Run() {
	fmt.Println("Bot is running...")
	time.Sleep(24 * time.Hour)
	fmt.Println("finished")
}

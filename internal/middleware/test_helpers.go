package middleware

import (
	"context"

	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func stubNext(called *bool) gotelegram.HandlerFunc {
	return func(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
		*called = true
	}
}

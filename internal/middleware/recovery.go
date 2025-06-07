package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/go-logr/logr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func Recovery(b *core.DaseinBot) gotelegram.Middleware {
	return func(next gotelegram.HandlerFunc) gotelegram.HandlerFunc {
		return func(ctx context.Context, bot *gotelegram.Bot, update *models.Update) {
			logger := logr.FromContextOrDiscard(ctx).WithName("recovery")
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					logger.Error(fmt.Errorf("panic: %v", r), "recovered panic in update handler", "stack", string(stack))
				}
			}()
			next(ctx, bot, update)
		}
	}
}

package middleware

import (
	"context"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/telegram"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func TestIsUserAllowed_Disallowed(t *testing.T) {
	botCore := &core.DaseinBot{
		Cfg:            &config.DaseinBotConfig{AllowedUsers: []int64{42}},
		TelegramClient: telegram.NewDummyClient(),
	}
	mw := IsUserAllowed(botCore)

	dummy := botCore.TelegramClient.(*telegram.DummyClient)
	called := false
	handler := mw(stubNext(&called))

	update := &models.Update{
		Message: &models.Message{
			From: &models.User{ID: 1},
			Chat: models.Chat{ID: 10},
			Text: "hello",
		},
	}
	handler(context.Background(), nil, update)

	if called {
		t.Error("IsUserAllowed: expected next NOT called for disallowed user")
	}
	if len(dummy.SentMessages) != 1 ||
		dummy.SentMessages[0].ChatID != 10 ||
		dummy.SentMessages[0].Text != "Get lost!" {
		t.Errorf("IsUserAllowed: expected dummy.SentMessages = [(10, \"Get lost!\")], got %v", dummy.SentMessages)
	}
}

func TestIsUserAllowed_Allowed(t *testing.T) {
	botCore := &core.DaseinBot{
		Cfg:            &config.DaseinBotConfig{AllowedUsers: []int64{42}},
		TelegramClient: telegram.NewDummyClient(),
		UserRepository: repository.NewUserRepository(mongo.NewDummyClient()),
	}
	mw := IsUserAllowed(botCore)

	dummy := botCore.TelegramClient.(*telegram.DummyClient)
	called := false
	handler := mw(stubNext(&called))

	update := &models.Update{
		Message: &models.Message{
			From: &models.User{ID: 42},
			Chat: models.Chat{ID: 20},
			Text: "hi",
		},
	}
	handler(context.Background(), nil, update)

	if !called {
		t.Error("IsUserAllowed: expected next called for allowed user")
	}
	if len(dummy.SentMessages) != 0 {
		t.Errorf("IsUserAllowed: expected no messages sent for allowed user, got %v", dummy.SentMessages)
	}
}

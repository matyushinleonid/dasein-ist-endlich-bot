package middleware

import (
	"context"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func TestRegisterUser_ExistingUser(t *testing.T) {
	ctx := context.Background()
	dc := mongo.NewDummyClient()
	repo := repository.NewUserRepository(dc)
	if _, err := repo.Create(ctx, 123, 10); err != nil {
		t.Fatalf("setup: Create error: %v", err)
	}

	bot := &core.DaseinBot{UserRepository: repo, Cfg: &config.DaseinBotConfig{OpenAIUserLimit: 30}}
	mw := RegisterUser(bot)

	called := false
	handler := mw(stubNext(&called))
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: 123},
			From: &models.User{ID: 456},
		},
	}

	handler(ctx, nil, update)

	if !called {
		t.Error("RegisterUser: expected next to be called for existing user")
	}
}

func TestRegisterUser_NewUser_CreatesAndCallsNext(t *testing.T) {
	ctx := context.Background()
	dc := mongo.NewDummyClient()
	repo := repository.NewUserRepository(dc)

	bot := &core.DaseinBot{UserRepository: repo, Cfg: &config.DaseinBotConfig{OpenAIUserLimit: 30}}
	mw := RegisterUser(bot)

	chatID := int64(789)
	called := false
	handler := mw(stubNext(&called))
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: chatID},
			From: &models.User{ID: chatID},
		},
	}

	handler(ctx, nil, update)

	if !called {
		t.Error("RegisterUser: expected next to be called after creating new user")
	}
	if _, err := repo.Get(ctx, chatID); err != nil {
		t.Errorf("RegisterUser: expected user %d to be created, got error: %v", chatID, err)
	}
}

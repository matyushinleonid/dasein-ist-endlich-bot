package middleware

import (
	"context"
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestIsUpdateLegal_MessageNil(t *testing.T) {
	mw := IsUpdateLegal(nil)
	called := false
	handler := mw(stubNext(&called))

	handler(context.Background(), nil, &models.Update{Message: nil})
	if called {
		t.Error("IsUpdateLegal: expected next NOT called when Message is nil")
	}
}

func TestIsUpdateLegal_EmptyText(t *testing.T) {
	mw := IsUpdateLegal(nil)
	called := false
	handler := mw(stubNext(&called))

	handler(context.Background(), nil, &models.Update{
		Message: &models.Message{Chat: models.Chat{ID: 1}, From: &models.User{ID: 1}, Text: ""},
	})
	if called {
		t.Error("IsUpdateLegal: expected next NOT called when Text is empty")
	}
}

func TestIsUpdateLegal_NonEmptyText(t *testing.T) {
	mw := IsUpdateLegal(nil)
	called := false
	handler := mw(stubNext(&called))

	handler(context.Background(), nil, &models.Update{
		Message: &models.Message{Chat: models.Chat{ID: 1}, From: &models.User{ID: 1}, Text: "hello"},
	})
	if !called {
		t.Error("IsUpdateLegal: expected next called when Text is non-empty")
	}
}

func TestIsUpdateLegal_ChatId_NotEqual_UserId(t *testing.T) {
	mw := IsUpdateLegal(nil)
	called := false
	handler := mw(stubNext(&called))

	handler(context.Background(), nil, &models.Update{
		Message: &models.Message{Chat: models.Chat{ID: 1}, From: &models.User{ID: 2}, Text: "hello"},
	})
	if called {
		t.Error("IsUpdateLegal: expected next NOT called when Chat ID is not equal to User ID")
	}
}

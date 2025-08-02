package handler

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/openai"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/redis"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/telegram"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
)

func makeBot() *core.DaseinBot {
	cfg := &config.DaseinBotConfig{
		Questions:       []string{"q1", "q2"},
		AnswerMaxLength: 5,
		DaysLeftMessage: "Days left in this world: %d",
	}
	ai := openai.NewDummyClient(0)
	ai.SendJSONOutput = "{\"days_left\":42,\"description\":\"some desc\"}"
	return &core.DaseinBot{
		Cfg:               cfg,
		SessionRepository: repository.NewSessionRepository(redis.NewDummyClient()),
		TelegramClient:    telegram.NewDummyClient(),
		OpenAIClient:      ai,
		UserRepository:    repository.NewUserRepository(mongo.NewDummyClient()),
	}
}

func TestBeginHandler(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))

	_, err := bot.UserRepository.Create(ctx, 10, 1)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	handler := BeginHandler(bot)
	upd := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 10}}}
	handler(ctx, &gotelegram.Bot{}, upd)

	d := bot.TelegramClient.(*telegram.DummyClient)
	if len(d.SentMessages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(d.SentMessages))
	}
	if got := d.SentMessages[0].Text; got != "q1" {
		t.Errorf("expected q1, got %q", got)
	}

	sess, err := bot.SessionRepository.Get(ctx, 10)
	if err != nil {
		t.Fatalf("expected session, got err %v", err)
	}
	if sess.Stage != 0 || len(sess.Answers) != 2 {
		t.Errorf("unexpected session: %+v", sess)
	}
}

func TestBeginHandler_LimitExceeded(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))

	// No OpenAI requests left
	_, err := bot.UserRepository.Create(ctx, 42, 0)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	handler := BeginHandler(bot)
	upd := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 42}}}
	handler(ctx, &gotelegram.Bot{}, upd)

	d := bot.TelegramClient.(*telegram.DummyClient)
	if len(d.SentMessages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(d.SentMessages))
	}
	want := "You have reached the limit of OpenAI requests. Contact the developer to increase your limit."
	if d.SentMessages[0].Text != want {
		t.Errorf("unexpected text:\n got %q\n want %q", d.SentMessages[0].Text, want)
	}
}

func TestAnswerHandler_FullFlow(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))

	_, err := bot.UserRepository.Create(ctx, 20, 1)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	err = bot.SessionRepository.Save(ctx, 20, &model.Session{Stage: 0, Answers: make([]string, 2)})
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	handler := AnswerHandler(bot)
	upd1 := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 20}, Text: "a1"}}
	handler(ctx, &gotelegram.Bot{}, upd1)

	d := bot.TelegramClient.(*telegram.DummyClient)
	if len(d.SentMessages) != 1 || d.SentMessages[0].Text != "q2" {
		t.Fatalf("expected next q2, got %v", d.SentMessages)
	}

	sess, err := bot.SessionRepository.Get(ctx, 20)
	if err != nil {
		t.Fatalf("got err %v", err)
	}
	if sess.Stage != 1 || sess.Answers[0] != "a1" {
		t.Errorf("session not updated: %+v", sess)
	}

	upd2 := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 20}, Text: "a2"}}
	handler(ctx, &gotelegram.Bot{}, upd2)

	if len(d.SentMessages) != 3 {
		t.Fatalf("expected final message count 3, got %d", len(d.SentMessages))
	}
	if d.SentMessages[1].Text != "some desc" {
		t.Errorf("unexpected final text: %q", d.SentMessages[1].Text)
	}
	if d.SentMessages[2].Text != "Days left in this world: 42" {
		t.Errorf("unexpected days left text: %q", d.SentMessages[2].Text)
	}

	if _, err = bot.SessionRepository.Get(ctx, 20); err == nil {
		t.Error("expected session deleted, still exists")
	}

	user, _ := bot.UserRepository.Get(ctx, 20)
	if user.OpenAIRequestsLeft != 0 {
		t.Errorf("expected OpenAI requests left to be 0, got %d", user.OpenAIRequestsLeft)
	}
}

func TestAnswerHandler_TooLongAnswer(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))

	err := bot.SessionRepository.Save(ctx, 30, &model.Session{Stage: 0, Answers: make([]string, 2)})
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	handler := AnswerHandler(bot)
	// 6 runes (AnswerMaxLength == 5)
	upd := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 30}, Text: "hello!"}}
	handler(ctx, &gotelegram.Bot{}, upd)

	d := bot.TelegramClient.(*telegram.DummyClient)
	if len(d.SentMessages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(d.SentMessages))
	}
	want := "Your message exceeds the maximum length of 5 characters. Please try again."
	if d.SentMessages[0].Text != want {
		t.Errorf("unexpected error text:\n got %q\n want %q", d.SentMessages[0].Text, want)
	}

	sess, err := bot.SessionRepository.Get(ctx, 30)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}
	if sess.Stage != 0 {
		t.Errorf("stage advanced on too-long answer: got %d, want 0", sess.Stage)
	}
}

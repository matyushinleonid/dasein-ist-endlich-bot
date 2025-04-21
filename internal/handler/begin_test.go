package handler

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	gotelegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		Questions: []string{"q1", "q2"},
	}
	ai := openai.NewDummyClient(0)
	ai.SendJSONOutput = "{\"days_left\":42,\"description\":\"some desc\"}"
	return &core.DaseinBot{
		Cfg:            cfg,
		RedisClient:    redis.NewDummyClient(),
		TelegramClient: telegram.NewDummyClient(),
		OpenAIClient:   ai,
		MongoClient:    mongo.NewDummyClient(),
	}
}

func TestBeginHandler(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))

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

	rs := bot.RedisClient.(*redis.DummyClient)
	sess, err := rs.Load(ctx, 10)
	if err != nil {
		t.Fatalf("expected session, got err %v", err)
	}
	if sess.Stage != 0 || len(sess.Answers) != 2 {
		t.Errorf("unexpected session: %+v", sess)
	}
}

func TestAnswerHandler_FullFlow(t *testing.T) {
	bot := makeBot()
	ctx := logr.NewContext(context.Background(), stdr.New(nil))
	rs := bot.RedisClient.(*redis.DummyClient)
	rs.Save(ctx, 20, &model.Session{Stage: 0, Answers: make([]string, 2)})

	handler := AnswerHandler(bot)
	upd1 := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 20}, Text: "a1"}}
	handler(ctx, &gotelegram.Bot{}, upd1)

	d := bot.TelegramClient.(*telegram.DummyClient)
	if len(d.SentMessages) != 1 || d.SentMessages[0].Text != "q2" {
		t.Fatalf("expected next q2, got %v", d.SentMessages)
	}

	sess, _ := rs.Load(ctx, 20)
	if sess.Stage != 1 || sess.Answers[0] != "a1" {
		t.Errorf("session not updated: %+v", sess)
	}

	upd2 := &models.Update{Message: &models.Message{Chat: models.Chat{ID: 20}, Text: "a2"}}
	handler(ctx, &gotelegram.Bot{}, upd2)

	if len(d.SentMessages) != 2 {
		t.Fatalf("expected final message count 2, got %d", len(d.SentMessages))
	}
	expected := "У вас осталось 42 дней в этом мире.\n\nsome desc"
	if d.SentMessages[1].Text != expected {
		t.Errorf("unexpected final text: %q", d.SentMessages[1].Text)
	}

	if _, err := rs.Load(ctx, 20); err == nil {
		t.Error("expected session deleted, still exists")
	}
}

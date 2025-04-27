package notifier

import (
	"context"
	"testing"
	"time"

	gotelegram "github.com/go-telegram/bot"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/adapter/repository"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/telegram"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/core"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

func TestNotifier_NotifyAll(t *testing.T) {
	ctx := context.Background()

	mc := mongo.NewDummyClient()
	now := time.Now()

	u1 := model.User{
		ID:                    1,
		DeathTime:             now.AddDate(0, 0, 5),
		NotificationFrequency: model.Daily,
		LastNotification:      time.Time{},
	}
	u2 := model.User{
		ID:                    2,
		DeathTime:             now,
		NotificationFrequency: model.Daily,
		LastNotification:      time.Time{},
	}

	if _, err := mc.Create(ctx, u1); err != nil {
		t.Fatalf("create u1: %v", err)
	}
	if _, err := mc.Create(ctx, u2); err != nil {
		t.Fatalf("create u2: %v", err)
	}

	tc := telegram.NewDummyClient()

	botCore := &core.DaseinBot{
		UserRepository: repository.NewUserRepository(mc),
		TelegramClient: tc,
	}
	notifier := New(botCore)

	if err := notifier.NotifyAll(ctx, (*gotelegram.Bot)(nil)); err != nil {
		t.Fatalf("NotifyAll failed: %v", err)
	}

	msgs := tc.SentMessages
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}

	want := map[int64]string{
		1: "У вас осталось 5 дней в этом мире.",
		2: "У вас осталось 0 дней в этом мире.",
	}

	for _, m := range msgs {
		expectedText, ok := want[m.ChatID]
		if !ok {
			t.Errorf("unexpected message for ChatID %d: %q", m.ChatID, m.Text)
			continue
		}
		if m.Text != expectedText {
			t.Errorf("message for ChatID %d = %q; want %q", m.ChatID, m.Text, expectedText)
		}
		delete(want, m.ChatID)
	}

	if len(want) > 0 {
		t.Errorf("missing messages for ChatIDs: %v", want)
	}
}

package notifier

import (
	"testing"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
	"github.com/robfig/cron/v3"
)

func TestFormatNotificationMessage(t *testing.T) {
	now := time.Date(2025, 4, 23, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		death    time.Time
		expected string
	}{
		{now.Add(5 * 24 * time.Hour), "Days left in this world: 5"},
		{now.Add(5*time.Hour + 1*time.Minute), "Days left in this world: 1"},
		{now.Add(-24 * time.Hour), "Days left in this world: 0"},
	}

	for _, tc := range tests {
		u := model.User{DeathTime: tc.death}
		got := FormatNotificationMessage(u, now)
		if got != tc.expected {
			t.Errorf("FormatNotificationMessage(%v): got %q, want %q",
				tc.death, got, tc.expected)
		}
	}
}

func mustParse(spec string) cron.Schedule {
	s, err := cron.ParseStandard(spec)
	if err != nil {
		panic(err)
	}
	return s
}

func TestShouldNotify(t *testing.T) {
	now := time.Date(2025, 4, 23, 9, 0, 00, 0, time.UTC)

	dailySpec := model.NotificationCronMap[model.Daily]
	if _, err := cron.ParseStandard(dailySpec); err != nil {
		t.Fatalf("invalid daily spec: %v", err)
	}

	tests := []struct {
		name    string
		user    model.User
		expect  bool
		wantErr bool
	}{
		{
			name:    "never",
			user:    model.User{NotificationFrequency: model.Never},
			expect:  false,
			wantErr: false,
		},
		{
			name: "first time",
			user: model.User{
				NotificationFrequency: model.Daily,
				LastNotification:      time.Time{},
			},
			expect: true,
		},
		{
			name: "just notified at 9:00",
			user: model.User{
				NotificationFrequency: model.Daily,
				LastNotification:      now,
			},
			expect: false,
		},
		{
			name: "missed 9:00 today",
			user: model.User{
				NotificationFrequency: model.Daily,
				LastNotification:      now.Add(-2 * time.Hour),
			},
			expect: true,
		},
		{
			name: "weekly not yet",
			user: model.User{
				NotificationFrequency: model.Weekly,
				LastNotification:      time.Date(2025, 4, 21, 9, 0, 0, 0, time.UTC),
			},
			expect: false,
		},
		{
			name: "weekly ready",
			user: model.User{
				NotificationFrequency: model.Weekly,
				LastNotification:      time.Date(2025, 4, 16, 9, 0, 0, 0, time.UTC),
			},
			expect: true,
		},
		{
			name:    "bad freq",
			user:    model.User{NotificationFrequency: "foo"},
			expect:  false,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		got, err := ShouldNotify(tc.user, now)
		if tc.wantErr {
			if err == nil {
				t.Errorf("%s: expected error, got nil", tc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
			continue
		}
		if got != tc.expect {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.expect)
		}
	}
}

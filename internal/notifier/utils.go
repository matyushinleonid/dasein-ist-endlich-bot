package notifier

import (
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
	"github.com/robfig/cron/v3"
)

func FormatNotificationMessage(u model.User, now time.Time) string {

	daysLeft := int((u.DeathTime.Sub(now)).Hours()/24 + 0.9999)
	if daysLeft < 0 {
		daysLeft = 0
	}
	return fmt.Sprintf("Days left in this world: %d", daysLeft)
}

func ShouldNotify(u model.User, now time.Time) (bool, error) {
	if u.NotificationFrequency == model.Never {
		return false, nil
	}

	spec, ok := model.NotificationCronMap[u.NotificationFrequency]
	if !ok || spec == "" {
		return false, fmt.Errorf("unknown notification frequency %q", u.NotificationFrequency)
	}

	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return false, fmt.Errorf("invalid cron spec %q: %w", spec, err)
	}

	if u.LastNotification.IsZero() {
		return true, nil
	}

	next := sched.Next(u.LastNotification)
	if !next.After(now) {
		return true, nil
	}
	return false, nil
}

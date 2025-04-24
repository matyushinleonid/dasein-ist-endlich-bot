package model

import "time"

type NotificationFrequency string

var (
	Daily   NotificationFrequency = "daily"
	Weekly  NotificationFrequency = "weekly"
	Monthly NotificationFrequency = "monthly"
	Yearly  NotificationFrequency = "yearly"
	Never   NotificationFrequency = "never"
)

var NotificationCronMap = map[NotificationFrequency]string{
	Daily:   "0 9 * * *",
	Weekly:  "0 9 * * 1",
	Monthly: "0 9 1 * *",
	Yearly:  "0 9 1 1 *",
	Never:   "",
}

type User struct {
	ID                    int64                 `bson:"_id"`
	Calculated            bool                  `bson:"calculated"`
	DeathTime             time.Time             `bson:"death_time"`
	LastNotification      time.Time             `bson:"last_notification"`
	NotificationFrequency NotificationFrequency `bson:"notification_frequency"`
}

func DeathTime(currentTime time.Time, daysLeft int) time.Time {
	return currentTime.AddDate(0, 0, daysLeft)
}

func NewUser(id int64) *User {
	return &User{
		ID:                    id,
		Calculated:            false,
		DeathTime:             time.Time{},
		LastNotification:      time.Time{},
		NotificationFrequency: Daily,
	}
}

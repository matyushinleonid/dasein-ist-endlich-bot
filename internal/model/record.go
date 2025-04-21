package model

type User struct {
	ID         int64 `bson:"_id"`
	DaysLeft   int64 `bson:"days_left"`
	Calculated bool  `bson:"calculated"`
}

package record

type Record struct {
	ID         int64 `bson:"_id"`
	DaysLeft   int64 `bson:"days_left"`
	Calculated bool  `bson:"calculated"`
}

package models

import (
	"time"
)

type TimeStamps struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func CreateTimeStamp(t *TimeStamps) {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now
}

func UpdateTimeStamp(t *TimeStamps) {
	now := time.Now().UTC()
	t.UpdatedAt = now
}

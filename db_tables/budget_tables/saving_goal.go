package db_tables

import (
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
)

type Goals struct {
	Goal_id         string    `bson:"goal_id" json:"goal_id"`
	Uid             string    `bson:"uid" json:"uid"`
	Goal_name       string    `bson:"goal_name" json:"goal_name"`
	Target_amount   float64   `bson:"target_amount" json:"target_amount"`
	Current_amount  float64   `bson:"current_amount" json:"current_amount"`
	Start_date      time.Time `bson:"start_date" json:"start_date"`
	Target_end_date time.Time `bson:"target_end_date" json:"target_end_date"`
	End_date        time.Time `bson:"end_date" json:"end_date"`
	Description     string    `bson:"description" json:"description"`
	db_tables.TimeStamps
}

package db_tables

import (
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
)

type Budget struct {
	Uid         string    `bson:"uid" json:"uid"`
	Budget_id   string    `bson:"budget_id" json:"budget_id"`
	Budget_name string    `bson:"budget_name" json:"budget_name"`
	Budget_desc string    `bson:"budget_desc" json:"budget_desc"`
	Amount      float64   `bson:"amount" json:"amount"`
	Start_date  time.Time `bson:"start_date" json:"start_date"`
	End_date    time.Time `bson:"end_date" json:"end_date"`
	db_tables.TimeStamps
}

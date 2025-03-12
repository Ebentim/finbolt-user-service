package db_tables

import (
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
)

type Rtrans_type string

const (
	Daily     Rtrans_type = "daily"
	Weekly    Rtrans_type = "weekly"
	Monthly   Rtrans_type = "monthly"
	Quarterly Rtrans_type = "quarterly"
	Yearly    Rtrans_type = "yearly"
)

type Recurring_transaction struct {
	Rtrans_id   string    `bson:"rtrans_id" json:"rtrans_id"`
	Uid         string    `bson:"uid" json:"uid"`
	Account     string    `bson:"account" json:"account"`
	Rtrans_desc string    `bson:"Rtrans_desc" json:"Rtrans_desc"`
	Amount      float64   `bson:"amount" json:"amount"`
	Start_date  time.Time `bson:"start_date" json:"start_date"`
	End_date    time.Time `bson:"end_date" json:"end_date"`
	db_tables.TimeStamps
}

package db_tables

import (
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
)

type CategoryType string

const (
	Income  CategoryType = "income"
	Expense CategoryType = "expense"
)

type Transaction struct {
	Uid               string       `bson:"uid" json:"uid"`
	Category_id       string       `bson:"category_id" json:"category_id"`
	Amount            float64      `bson:"amount" json:"amount"`
	Transaction_date  time.Time    `bson:"transaction_date" json:"transaction_date"`
	Trans_description string       `bson:"trans_description" json:"trans_description"`
	Category_name     string       `bson:"category_name" json:"category_name"`
	Category_type     CategoryType `bson:"category_type" json:"category_type"`
	db_tables.TimeStamps
}

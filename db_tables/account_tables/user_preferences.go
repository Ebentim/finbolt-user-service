package db_tables

import "github.com/Ebentim/finbolt-user-service/db_tables"

type User_Setting struct {
	Uid          string `bson:"uid" json:"uid"`
	Currency     string `bson:"currency" json:"currency"`
	Lang         string `bson:"lang" json:"lang"`
	Date_format  string `bson:"date_format" json:"date_format"`
	Trans_report string `bson:"trans_report" json:"trans_report"`
	db_tables.TimeStamps
}

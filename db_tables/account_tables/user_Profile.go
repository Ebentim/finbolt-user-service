package db_tables

import "github.com/Ebentim/finbolt-user-service/db_tables"

type User_Profile struct {
	Uid   string `bson:"uid" json:"uid"`
	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	db_tables.TimeStamps
}

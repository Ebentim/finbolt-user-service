package models

import "time"

type User_Subscription struct {
	Uid         string    `bson:"uid" json:"uid"`
	Last_Plan   string    `bson:"last_plan" json:"last_plan"`
	Active_Plan string    `bson:"active_plan" json:"active_plan"`
	Sub_status  string    `bson:"sub_status" json:"sub_status"`
	Start_Date  time.Time `bson:"start_date" json:"start_date"`
	End_Date    time.Time `bson:"end_date" json:"end_date"`
	TimeStamps
}

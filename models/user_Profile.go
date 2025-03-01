package models

type User_Profile struct {
	Uid   string `bson:"uid" json:"uid"`
	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	TimeStamps
}

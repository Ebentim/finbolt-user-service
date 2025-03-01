package models

type User_Account struct {
	Uid             string   `bson:"uid" json:"uid"`
	Email           string   `bson:"email" json:"email"`
	Email_verfified bool     `bson:"email_verified" json:"email_verified"`
	Interests       []string `bson:"interests" json:"interests"`
	Onboarded       bool     `bson:"onboarded" json:"onboarded"`
	Login_provider  string   `bson:"login_provider" json:"login_provider"`
	User_role       []string `bson:"user_role" json:"user_role"`
	TimeStamps
}

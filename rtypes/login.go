package rtypes

import "github.com/Ebentim/finbolt-user-service/models"

type LoginResponseFormat struct {
	Us_acc models.User_Account      `json:"us_acc"`
	Us_P   models.User_Profile      `json:"us_p"`
	Us_Sub models.User_Subscription `json:"us_sub"`
}

type Result struct {
			Data any
			Err  error
		}

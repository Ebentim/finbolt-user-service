package rtypes

import (
	acc_table "github.com/Ebentim/finbolt-user-service/db_tables/account_tables"
	db_tables "github.com/Ebentim/finbolt-user-service/db_tables/budget_tables"
)

type UserResponse struct {
	Us_acc acc_table.User_Account      `json:"us_acc"`
	Us_Pro acc_table.User_Profile      `json:"us_pro"`
	Us_Sub acc_table.User_Subscription `json:"us_sub"`
}

type Result struct {
	Data any
	Err  error
}

type BudgetResponse struct {
	Bud    []db_tables.Budget                `json:"bud"`
	Rtrans []db_tables.Recurring_transaction `json:"rtrans"`
	Goals  []db_tables.Goals                 `json:"goals"`
	T_list []db_tables.Transaction           `json:"t_list"`
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
	acc_table "github.com/Ebentim/finbolt-user-service/db_tables/account_tables"
	"github.com/Ebentim/finbolt-user-service/validate"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUserAccount(db *mongo.Database, r *http.Request) error {
	ua_col := db.Collection("user_accounts")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var us_acc acc_table.User_Account

	if err := json.NewDecoder(r.Body).Decode(&us_acc); err != nil {
		return fmt.Errorf("user with uid %s account data could not be created", us_acc.Uid)
	}
	defer r.Body.Close()

	if us_acc.Uid == "" {
		return fmt.Errorf("user with uid %s account data could not be created", us_acc.Uid)
	}
	// Check for the existence of uid
	filter := bson.M{"uid": us_acc.Uid}
	result := ua_col.FindOne(ctx, filter)

	if result.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("user with uid %s already exists", us_acc.Uid)
	}
	//Create time stamps before validation to avoid time stamp errors
	db_tables.CreateTimeStamp(&us_acc.TimeStamps)

	if err := validate.Validate_User_Account(&us_acc); err != nil {
		return err
	}

	_, err := ua_col.InsertOne(ctx, &us_acc)
	if err != nil {
		return fmt.Errorf("user with uid %s account data could not be created", us_acc.Uid)
	}

	return nil
}

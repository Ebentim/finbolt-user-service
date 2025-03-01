package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ebentim/finbolt-user-service/models"
	"github.com/Ebentim/finbolt-user-service/validate"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUserAccount(db *mongo.Database, r *http.Request) error {
	ua_col := db.Collection("user_accounts")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var us_acc models.User_Account

	if err := json.NewDecoder(r.Body).Decode(&us_acc); err != nil {
		return fmt.Errorf("user with uid %s account data could not be created", us_acc.Uid)
	}

	if err := validate.Validate_User_Account(&us_acc); err != nil {
		return err
	}

	p := &us_acc
	p.Uid = r.FormValue("uid")
	p.Email = r.FormValue("email")

	_, err := ua_col.InsertOne(ctx, &us_acc)
	if err != nil {
		return fmt.Errorf("user with uid %s account data could not be created", us_acc.Uid)
	}

	return nil
}

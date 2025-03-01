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

func CreateSubscriptionProfile(db *mongo.Database, r *http.Request) error {

	us_col := db.Collection("user_subscriptions")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var us_sub models.User_Subscription

	if err := json.NewDecoder(r.Body).Decode(&us_sub); err != nil {
		return fmt.Errorf("user with uid %s subscription data creation failed", us_sub.Uid)
	}

	if err := validate.Validate_User_Subscription(&us_sub); err != nil {
		return err
	}

	p := &us_sub
	p.Uid = r.FormValue("uid")
	p.Sub_status = "notsubscribed"
	models.CreateTimeStamp(&us_sub.TimeStamps)

	_, err := us_col.InsertOne(ctx, &us_sub)
	if err != nil {
		return err
	}
	return nil

}

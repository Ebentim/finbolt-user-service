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

func CreateSubscriptionProfile(db *mongo.Database, r *http.Request) error {
	// Create collection reference
	us_col := db.Collection("user_subscriptions")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel() // This ensures context resources are cleaned up

	// Initialize subscription struct
	var us_sub acc_table.User_Subscription

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&us_sub); err != nil {
		return fmt.Errorf("user with uid %s subscription data creation failed", us_sub.Uid)
	}
	defer r.Body.Close() // Ensure request body is closed

	if us_sub.Uid == "" {
		return fmt.Errorf("user with uid %s subscription data creation failed", us_sub.Uid)
	}

	filter := bson.M{"uid": us_sub.Uid}
	result := us_col.FindOne(ctx, filter)

	if result.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("subscription data for %s axists", us_sub.Uid)
	}
	// Validate subscription
	if err := validate.Validate_User_Subscription(&us_sub); err != nil {
		return err
	}

	// Set subscription fields
	p := &us_sub
	p.Sub_status = "notsubscribed"
	db_tables.CreateTimeStamp(&us_sub.TimeStamps)

	// Insert into database
	_, err := us_col.InsertOne(ctx, &us_sub)
	if err != nil {
		return err
	}
	return nil
}

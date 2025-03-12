package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
	acc_tables "github.com/Ebentim/finbolt-user-service/db_tables/account_tables"
	"github.com/Ebentim/finbolt-user-service/validate"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUserProfile(db *mongo.Database, r *http.Request) error {
	ua_col := db.Collection("user_profile")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var user acc_tables.User_Profile

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return fmt.Errorf("invalid request body")
	}
	defer r.Body.Close()

	filter := bson.M{"uid": user.Uid}
	if result := ua_col.FindOne(ctx, filter); result.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("user already exists %s", user.Uid)
	}

	db_tables.CreateTimeStamp(&user.TimeStamps)

	if err := validate.Validate_User_Profile(&user); err != nil {
		return err
	}

	_, err := ua_col.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

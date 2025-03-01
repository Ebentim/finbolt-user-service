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

func CreateUserProfile(db *mongo.Database, r *http.Request) error {
	ua_col := db.Collection("user_profile")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var user models.User_Profile

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return fmt.Errorf("invalid request body")

	}

	if err := validate.Validate_User_Profile(&user); err != nil {
		return err

	}

	p := &user
	p.Uid = r.FormValue("uid")
	p.Name = r.FormValue("name")
	models.CreateTimeStamp(&user.TimeStamps)

	_, err := ua_col.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

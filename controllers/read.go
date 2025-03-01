// package controllers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"time"
// 	"github.com/Ebentim/finbolt-user-service/models"
// 	"github.com/Ebentim/finbolt-user-service/rtypes"
// 	"github.com/Ebentim/finbolt-user-service/services"

// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// func ListAllUsers(db *mongo.Database) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		collections := map[string]*mongo.Collection{
// 			"accounts":     db.Collection("user_accounts"),
// 			"profiles":     db.Collection("user_profiles"),
// 			"subscription": db.Collection("user_subscription"),
// 		}

// 		ctx := r.Context()

// 		// limit for pagination
// 		limit := 10 //Default limit

// 		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
// 			if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
// 				limit = parsedLimit
// 			}
// 		}

// 		skip := 0

// 		if skipParam := r.URL.Query().Get("skip"); skipParam != "" {
// 			if parsedSkip, err := strconv.Atoi(skipParam); err == nil && parsedSkip > 0 {
// 				skip = parsedSkip
// 			}
// 		}

// 		findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(bson.D{{"created_at", -1}})
// 		pro_ch := make(chan []rtypes.Result)

// 		go func() {
// 			var account [limit]models.User_Account
// 			cursor, err := collections["accountsd"].Find(ctx, bson.M{}, findOptions)

// 			if err != nil {
//             log.Printf("Error fetching users: %v", err)
//             http.Error(w, "Error fetching data", http.StatusInternalServerError)
//             return
// 		}()

// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "Error fetching data", http.StatusInternalServerError)
// 			return
// 		}

// 		defer allDoc.Close(ctx)

// 		var u_pro models.User_Profile
// 		var u_acc models.User_Account
// 		var u_sub models.User_Subscription

// 		w.Header().Set("Content-type", "application/json")

// 		if err := allDoc.All(ctx, &u_acc); err != nil {
// 			log.Println(err)
// 			http.Error(w, "Error decoding data", http.StatusNotFound)
// 			return
// 		}

// 		if err := json.NewEncoder(w).Encode(u_acc); err != nil {
// 			log.Println(err)
// 			http.Error(w, "Error encoding data", http.StatusInternalServerError)
// 			return

// 		}
// 	}
// }

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ebentim/finbolt-user-service/models"
	"github.com/Ebentim/finbolt-user-service/rtypes"
	"github.com/Ebentim/finbolt-user-service/services"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// UserResponse represents the complete user data with profile, account and subscription
type UserResponse struct {
	Account      models.User_Account      `json:"account"`
	Profile      models.User_Profile      `json:"profile"`
	Subscription models.User_Subscription `json:"subscription"`
}

// PaginatedResponse wraps the user data with pagination metadata
type PaginatedResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"totalCount"`
	Limit      int            `json:"limit"`
	Skip       int            `json:"skip"`
	HasMore    bool           `json:"hasMore"`
}

func ListAllUsers(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collections := map[string]*mongo.Collection{
			"accounts":     db.Collection("user_accounts"),
			"profiles":     db.Collection("user_profiles"),
			"subscription": db.Collection("user_subscription"),
		}

		ctx := r.Context()

		// Parse pagination parameters
		limit := 10 // Default limit
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		skip := 0 // Default skip
		if skipParam := r.URL.Query().Get("skip"); skipParam != "" {
			if parsedSkip, err := strconv.Atoi(skipParam); err == nil && parsedSkip >= 0 {
				skip = parsedSkip
			}
		}

		// Set up options for pagination and sorting
		findOptions := options.Find().
			SetLimit(int64(limit)).
			SetSkip(int64(skip)).
			SetSort(bson.D{{"created_at", -1}})

		// First, get profiles with pagination
		profileCursor, err := collections["profiles"].Find(ctx, bson.M{}, findOptions)
		if err != nil {
			log.Printf("Error fetching profiles: %v", err)
			http.Error(w, "Error fetching user profiles", http.StatusInternalServerError)
			return
		}
		defer profileCursor.Close(ctx)

		// Get total count for pagination metadata
		totalCount, err := collections["profiles"].CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Printf("Error counting profiles: %v", err)
			http.Error(w, "Error counting user profiles", http.StatusInternalServerError)
			return
		}

		// Decode profiles
		var profiles []models.User_Profile
		if err := profileCursor.All(ctx, &profiles); err != nil {
			log.Printf("Error decoding profiles: %v", err)
			http.Error(w, "Error processing user profiles", http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := PaginatedResponse{
			Users:      make([]UserResponse, 0, len(profiles)),
			TotalCount: totalCount,
			Limit:      limit,
			Skip:       skip,
			HasMore:    (int64(skip) + int64(len(profiles))) < totalCount,
		}

		// For each profile, fetch the corresponding account and subscription
		for _, profile := range profiles {
			userResponse := UserResponse{
				Profile: profile,
			}

			// Find the corresponding account using the profile's UID
			var account models.User_Account
			err := collections["accounts"].FindOne(ctx, bson.M{"uid": profile.Uid}).Decode(&account)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					log.Printf("Error fetching account for UID %s: %v", profile.Uid, err)
				}
				// Continue even if no account is found
			} else {
				userResponse.Account = account
			}

			// Find the corresponding subscription using the profile's UID
			var subscription models.User_Subscription
			err = collections["subscription"].FindOne(ctx, bson.M{"uid": profile.Uid}).Decode(&subscription)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					log.Printf("Error fetching subscription for uid %s: %v", profile.Uid, err)
				}
				// Continue even if no subscription is found
			} else {
				userResponse.Subscription = subscription
			}

			response.Users = append(response.Users, userResponse)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Encode and return the response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			return
		}
	}
}

func LoginUser(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.FormValue("uid")
		fmt.Println("uid", uid)
		if uid == "" {
			http.Error(w, "missing user ID", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Create a context with timeout for database operations
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get collection references
		collections := map[string]*mongo.Collection{
			"accounts":     db.Collection("user_accounts"),
			"profiles":     db.Collection("user_profiles"),
			"subscription": db.Collection("user_subscription"),
		}

		// Create a struct to hold results and errors
		type result struct {
			data any
			err  error
		}

		accountCh := make(chan result, 1)
		profileCh := make(chan result, 1)
		subscriptionCh := make(chan result, 1)

		// Goroutines for parallel fetching
		go func() {
			var account models.User_Account
			err := services.FetchDocument(ctx, collections["accounts"], uid, &account)
			accountCh <- result{data: account, err: err}
		}()

		go func() {
			var profile models.User_Profile
			err := services.FetchDocument(ctx, collections["profiles"], uid, &profile)
			profileCh <- result{data: profile, err: err}
		}()

		go func() {
			var subscription models.User_Subscription
			err := services.FetchDocument(ctx, collections["subscription"], uid, &subscription)
			subscriptionCh <- result{data: subscription, err: err}
		}()

		// Collect results
		accountResult := <-accountCh
		profileResult := <-profileCh
		subscriptionResult := <-subscriptionCh

		// Check for errors
		for _, res := range []result{accountResult, profileResult, subscriptionResult} {
			if res.err != nil {
				http.Error(w, "User data not found: "+res.err.Error(), http.StatusNotFound)
				return
			}
		}

		// Prepare response
		responseData := rtypes.LoginResponseFormat{
			Us_acc: accountResult.data.(models.User_Account),
			Us_P:   profileResult.data.(models.User_Profile),
			Us_Sub: subscriptionResult.data.(models.User_Subscription),
		}

		// Marshal and send response
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "Error creating response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}

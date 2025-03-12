package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ebentim/finbolt-user-service/rtypes"
	"github.com/Ebentim/finbolt-user-service/services"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// PaginatedResponse wraps the user data with pagination metadata
type PaginatedResponse struct {
	Users      []rtypes.UserResponse `json:"users"`
	TotalCount int64                 `json:"totalCount"`
	Limit      int                   `json:"limit"`
	Skip       int                   `json:"skip"`
	HasMore    bool                  `json:"hasMore"`
}

func ListAllUsers(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add timeout to prevent hanging requests
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		collections := map[string]*mongo.Collection{
			"us_acc": db.Collection("user_accounts"),
			"us_pro": db.Collection("user_profiles"),
			"us_sub": db.Collection("user_subscription"),
		}

		// Verify collections are properly initialized
		if collections["us_pro"] == nil {
			log.Printf("Error: user_profiles collection is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

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

		pipeline := services.CreateUserAggregationPipeline(services.UserAggregationOptions{
			Skip:  skip,
			Limit: limit,
		})

		cursor, err := db.Collection("user_profiles").Aggregate(ctx, pipeline)
		if err != nil {
			log.Printf("Error in aggregation: %v", err)
			http.Error(w, "Error fetching user data", http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		var results []rtypes.UserResponse
		if err := cursor.All(ctx, &results); err != nil {
			log.Printf("Error decoding results: %v", err)
			http.Error(w, "Error processing user data", http.StatusInternalServerError)
			return
		}

		// Get total count for pagination metadata
		totalCount, err := collections["us_pro"].CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Printf("Error counting profiles: %v", err)
			http.Error(w, "Error counting user profiles", http.StatusInternalServerError)
			return
		}

		response := PaginatedResponse{
			Users:      results,
			TotalCount: totalCount,
			Limit:      limit,
			Skip:       skip,
			HasMore:    (int64(skip) + int64(len(results))) < totalCount,
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

		if uid == "" {
			http.Error(w, "missing user ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		pipeline := services.CreateUserAggregationPipeline(services.UserAggregationOptions{
			UID:   uid,
			Limit: 1,
		})

		var result rtypes.UserResponse
		cursor, err := db.Collection("user_profiles").Aggregate(ctx, pipeline)
		if err != nil {
			http.Error(w, "Error fetching user data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if !cursor.Next(ctx) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err := cursor.Decode(&result); err != nil {
			http.Error(w, "Error decoding user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

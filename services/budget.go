package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ebentim/finbolt-user-service/db_tables"
	budget_tables "github.com/Ebentim/finbolt-user-service/db_tables/budget_tables"
	"github.com/Ebentim/finbolt-user-service/validate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateBudget(db *mongo.Database, r *http.Request) error {
	budget_col := db.Collection("budgets")
	us_Pro := db.Collection("user_profiles")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var budget budget_tables.Budget

	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		return fmt.Errorf("failed to decode budget request: %w", err)
	}
	// create timestamps
	db_tables.CreateTimeStamp(&budget.TimeStamps)

	// Validate the budget
	if err := validate.Validate_Budget(&budget); err != nil {
		return fmt.Errorf("budget validation failed: %w", err)
	}

	// Validate date range
	if !budget.Start_date.IsZero() && !budget.End_date.IsZero() {
		if budget.Start_date.After(budget.End_date) {
			return fmt.Errorf("start date cannot be after end date")
		}
	}

	// Ensure dates are in UTC
	budget.Start_date = budget.Start_date.UTC()
	budget.End_date = budget.End_date.UTC()

	// Find User by ID
	filter := bson.M{"uid": budget.Uid}
	result := us_Pro.FindOne(ctx, filter)

	if result.Err() != nil {
		return fmt.Errorf("failed to find user profile: %w", result.Err())
	}

	_, err := budget_col.InsertOne(ctx, budget)
	if err != nil {
		return fmt.Errorf("failed to insert budget '%s': %w", budget.Budget_name, err)
	}

	return nil
}

func EditBudget(db *mongo.Database, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	budget_col := db.Collection("budgets")

	var budget budget_tables.Budget

	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		return fmt.Errorf("failed to decode budget data: %v", err)
	}

	if err := validate.Validate_Budget(&budget); err != nil {
		return err
	}

	budget_id := r.FormValue("budget_id")
	uid := r.FormValue("uid")

	if budget_id == "" || uid == "" {
		return fmt.Errorf("user not found is required")
	}

	filter := bson.M{"uid": uid, "budget_id": budget_id}

	// Parse dates
	startDate, err := time.Parse(time.RFC3339, r.FormValue("start_date"))
	if err != nil {
		return fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse(time.RFC3339, r.FormValue("end_date"))
	if err != nil {
		return fmt.Errorf("invalid end date format: %v", err)
	} else if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	// Parse amount
	amount, err := strconv.ParseFloat(r.FormValue("budget_amount"), 64)
	if err != nil {
		return fmt.Errorf("invalid amount format: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"budget_name": r.FormValue("budget_name"),
			"amount":      amount,
			"budget_desc": r.FormValue("budget_description"),
			"start_date":  startDate,
			"end_date":    endDate,
			"updatedAt":   time.Now().Format(time.RFC3339),
		},
	}

	result, err := budget_col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update budget: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("budget with id %s not found", budget_id)
	}

	return nil
}

func DeleteBudget(db *mongo.Database, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()
	budget_col := db.Collection("budgets")

	uid := r.FormValue("uid")
	budget_id := r.FormValue("budget_id")

	// Validate required fields
	if uid == "" || budget_id == "" {
		return fmt.Errorf("uid and budget_id are required")
	}

	filter := bson.M{"uid": uid, "budget_id": budget_id}

	result := budget_col.FindOneAndDelete(ctx, filter)

	if result.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("budget with id %s not found", budget_id)
	} else if result.Err() != nil {
		return fmt.Errorf("failed to delete budget: %v", result.Err())
	}

	return nil
}

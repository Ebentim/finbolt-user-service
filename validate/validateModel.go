package validate

import (
	"fmt"
	"strings"
	"time"

	acc_table "github.com/Ebentim/finbolt-user-service/db_tables/account_tables"
	budget_tables "github.com/Ebentim/finbolt-user-service/db_tables/budget_tables"
)

// fieldError represents a validation error for a specific field
type fieldError struct {
	field   string
	message string
}

// validateField is a helper function to validate a field and collect errors
func validateField(field any, fieldName, message string) *fieldError {
	switch v := field.(type) {
	case string:
		if v == "" {
			return &fieldError{fieldName, message}
		}
	case float64:
		if v == 0 {
			return &fieldError{fieldName, message}
		}
	case time.Time:
		if v.IsZero() {
			return &fieldError{fieldName, message}
		}
	case []string:
		if len(v) == 0 {
			return &fieldError{fieldName, message}
		}
	}
	return nil
}

// Validate_Budget validates the budget struct and returns any validation errors
func Validate_Budget(budget *budget_tables.Budget) error {
	var errors []string

	fields := []struct {
		value   any
		field   string
		message string
	}{
		{budget.Uid, "uid", "user id is required"},
		{budget.Budget_id, "budget_id", "budget id is required"},
		{budget.Budget_name, "budget_name", "budget name is required"},
		{budget.Budget_desc, "budget_desc", "budget description is required"},
		{budget.Amount, "amount", "budget amount is required"},
		{budget.Start_date, "start_date", "budget start date is required"},
		{budget.End_date, "end_date", "budget end date is required"},
	}

	for _, field := range fields {
		if err := validateField(field.value, field.field, field.message); err != nil {
			errors = append(errors, err.message)
		}
	}

	// Additional date validations
	if !budget.Start_date.IsZero() && !budget.End_date.IsZero() {
		if budget.Start_date.After(budget.End_date) {
			errors = append(errors, "start date must be before end date")
		}
	}

	// Validate amount is positive
	if budget.Amount < 0 {
		errors = append(errors, "budget amount cannot be negative")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

func Validate_User_Profile(user *acc_table.User_Profile) error {
	var errors []string

	fields := []struct {
		value   any
		field   string
		message string
	}{
		{user.Uid, "uid", "uid is required"},
		{user.Name, "name", "name is required"},
		{user.Image, "image", "image is required"},
	}

	for _, field := range fields {
		if err := validateField(field.value, field.field, field.message); err != nil {
			errors = append(errors, err.message)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

func Validate_User_Account(user *acc_table.User_Account) error {
	var errors []string

	fields := []struct {
		value   any
		field   string
		message string
	}{
		{user.Uid, "uid", "uid is required"},
		{user.Email, "email", "email is required"},
		{user.User_role, "user_role", "user role is required"},
	}

	for _, field := range fields {
		if err := validateField(field.value, field.field, field.message); err != nil {
			errors = append(errors, err.message)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	return nil
}

func Validate_User_Subscription(user *acc_table.User_Subscription) error {
	if user.Uid == "" {
		return fmt.Errorf("user not found")
	}
	return nil
}

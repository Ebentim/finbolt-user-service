package validate

import (
	"fmt"

	"github.com/Ebentim/finbolt-user-service/models"
)

func Validate_User_Profile(user *models.User_Profile) error {
	if user.Uid == "" {
		return fmt.Errorf("uid is required")

	}

	if user.Name == "" {
		return fmt.Errorf("name is required")

	}

	if user.Image == "" {
		return fmt.Errorf("image is required")

	}

	return nil
}

func Validate_User_Account(user *models.User_Account) error {
	if user.Uid == "" {
		return fmt.Errorf("user id is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	if len(user.User_role) == 0 {
		return fmt.Errorf("user role is required")
	}

	return nil
}

func Validate_User_Subscription(user *models.User_Subscription) error {
	if user.Uid == "" {
		return fmt.Errorf("user not found")
	}
	return nil
}

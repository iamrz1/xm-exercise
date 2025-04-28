package models

import (
	"errors"

	"xm-exercise/internal/utils"
)

// UserRegistration represents register credentials
// @Description User credentials for registration
type UserRegistration struct {
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"securepassword123"`
	Email    string `json:"email" example:"john@example.com"`
}

// Validate validates user registration credentials
func (c *UserRegistration) Validate() error {
	if len(c.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if !utils.IsValidEmail(c.Email) {
		return errors.New("invalid email address")
	}

	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

// UserLogin represents login credentials
// @Description User credentials for registration
type UserLogin struct {
	Password string `json:"password" example:"securepassword123"`
	Email    string `json:"email,omitempty" example:"john@example.com"`
}

// Validate validates user login credentials
func (c *UserLogin) Validate() error {
	if !utils.IsValidEmail(c.Email) {
		return errors.New("invalid email address")
	}

	return nil
}

// TokenResponse defines the response format for auth endpoints
type TokenResponse struct {
	// JWT token for authentication
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

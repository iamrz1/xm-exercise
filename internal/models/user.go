package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCredentials represents login/register credentials
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

// Validate validates user credentials
func (c UserCredentials) Validate() error {
	if len(c.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

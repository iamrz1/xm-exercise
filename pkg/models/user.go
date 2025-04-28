package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
// @Description User model with authentication details
type User struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"size:50;uniqueIndex;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate is hook for validation and mutation before creating an object
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate is hook for validation and mutation before updating an object
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.ID = ""
	return nil
}

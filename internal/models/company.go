package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CompanyType represents the type of a company
type CompanyType string

const (
	TypeCorporation    CompanyType = "Corporations"
	TypeNonProfit      CompanyType = "NonProfit"
	TypeCooperative    CompanyType = "Cooperative"
	TypeSoleProprietor CompanyType = "Sole Proprietorship"
)

// Company represents a company entity
type Company struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	Description   *string     `json:"description,omitempty"`
	EmployeeCount int         `json:"employee_count"`
	Registered    bool        `json:"registered"`
	Type          CompanyType `json:"type"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// Validate validates company fields
func (c Company) Validate() error {
	if c.ID == uuid.Nil {
		return errors.New("id is required")
	}

	if len(c.Name) == 0 {
		return errors.New("name is required")
	}

	if len(c.Name) > 15 {
		return errors.New("name must be 15 characters or less")
	}

	if c.Description != nil && len(*c.Description) > 3000 {
		return errors.New("description must be 3000 characters or less")
	}

	if c.EmployeeCount <= 0 {
		return errors.New("employee count must be positive")
	}

	switch c.Type {
	case TypeCorporation, TypeNonProfit, TypeCooperative, TypeSoleProprietor:
		// valid
	default:
		return errors.New("invalid company type")
	}

	return nil
}

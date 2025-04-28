package models

import (
	"errors"
	"time"
)

// CompanyCreateRequest represents a company to be created
// @Description Contains details to successfully create a Company
type CompanyCreateRequest struct {
	Name          string      `json:"name" example:"Acme Corp"`
	Description   *string     `json:"description,omitempty" example:"Leading provider of widgets"`
	EmployeeCount int         `json:"employee_count" example:"42"`
	Registered    *bool       `json:"registered" example:"true"`
	Type          CompanyType `json:"type" example:"Corporations"`
}

// Validate validates company fields
func (c *CompanyCreateRequest) Validate() error {
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

	if c.Registered == nil {
		return errors.New("registered field is required")
	}

	switch c.Type {
	case TypeCorporation, TypeNonProfit, TypeCooperative, TypeSoleProprietor:
	default:
		return errors.New("invalid company type")
	}

	return nil
}

type CompanyUpdateRequest struct {
	Name          *string     `json:"name" example:"Acme Corp"`
	Description   *string     `json:"description,omitempty" example:"Leading provider of widgets"`
	EmployeeCount *int        `json:"employee_count" example:"42"`
	Registered    *bool       `json:"registered" example:"true"`
	Type          CompanyType `json:"type" example:"Corporations"`
}

// Validate validates company fields
func (c *CompanyUpdateRequest) Validate() error {
	if c.Name != nil && len(*c.Name) == 0 {
		return errors.New("name is required")
	}

	if c.Name != nil && len(*c.Name) > 15 {
		return errors.New("name must be 15 characters or less")
	}

	if c.Description != nil && len(*c.Description) > 3000 {
		return errors.New("description must be 3000 characters or less")
	}

	if c.Name != nil && *c.EmployeeCount <= 0 {
		return errors.New("employee count must be positive")
	}

	switch c.Type {
	case TypeCorporation, TypeNonProfit, TypeCooperative, TypeSoleProprietor:
	default:
		return errors.New("invalid company type")
	}

	return nil
}

type CompanyResponse struct {
	ID            string      `json:"id" example:"df45-adf32.....e-358dc"`
	Name          string      `json:"name" example:"Acme Corp"`
	Description   *string     `json:"description" example:"Leading provider of widgets"`
	EmployeeCount int         `json:"employee_count" example:"42"`
	Registered    *bool       `json:"registered" example:"true"`
	Type          CompanyType `json:"type" example:"Corporations"`
	CreatedAt     time.Time   `json:"created_at" example:"05-04-2013"`
	UpdatedAt     time.Time   `json:"updated_at" example:"05-04-2013"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// CompanyType represents the type of a company
// @Description Type of company
type CompanyType string

const (
	TypeCorporation    CompanyType = "Corporations"
	TypeNonProfit      CompanyType = "NonProfit"
	TypeCooperative    CompanyType = "Cooperative"
	TypeSoleProprietor CompanyType = "Sole Proprietorship"
)

// Company represents a company entity
// @Description Company model with all details
type Company struct {
	ID            string      `gorm:"type:uuid;primaryKey"`
	Name          string      `gorm:"size:15;uniqueIndex;not null"`
	Description   *string     `gorm:"size:3000"`
	EmployeeCount int         `gorm:"not null"`
	Registered    *bool       `gorm:"not null"`
	Type          CompanyType `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime"`
}

// BeforeCreate is hook for validation and mutation before creating the object
func (c *Company) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (c *Company) ToResponse() *CompanyResponse {
	return &CompanyResponse{
		ID:            c.ID,
		Name:          c.Name,
		Description:   c.Description,
		EmployeeCount: c.EmployeeCount,
		Registered:    c.Registered,
		Type:          c.Type,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

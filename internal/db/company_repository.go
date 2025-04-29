package db

import (
	"errors"

	"gorm.io/gorm"

	"xm-exercise/pkg/models"
)

// CompanyRepository handles database operations for companies
type CompanyRepository struct {
	db *Database
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *Database) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create inserts a new company into the database
func (r *CompanyRepository) Create(company *models.Company) error {
	result := r.db.Create(company)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves a company by its ID
func (r *CompanyRepository) GetByID(id string) (*models.Company, error) {
	var company models.Company
	result := r.db.First(&company, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("company not found")
		}
		return nil, result.Error
	}
	return &company, nil
}

// Update updates an existing company
func (r *CompanyRepository) Update(company *models.Company) error {
	result := r.db.Model(&company).Updates(company)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("company not found")
	}

	return nil
}

// Delete removes a company by its ID
func (r *CompanyRepository) Delete(id string) error {
	result := r.db.Delete(&models.Company{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("company not found")
	}

	return nil
}

// ExistsByName checks if a company with the given name exists
func (r *CompanyRepository) ExistsByName(name string) (bool, error) {
	var count int64
	query := r.db.Model(&models.Company{}).Where("name = ?", name)

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamrz1/xm-exercise/internal/models"
)

// CompanyRepository handles database operations for companies
type CompanyRepository struct {
	db *PostgresDB
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *PostgresDB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create inserts a new company into the database
func (r *CompanyRepository) Create(company models.Company) error {
	query := `
		INSERT INTO companies (id, name, description, employee_count, registered, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		company.ID,
		company.Name,
		company.Description,
		company.EmployeeCount,
		company.Registered,
		company.Type,
		company.CreatedAt,
		company.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating company: %w", err)
	}

	return nil
}

// GetByID retrieves a company by its ID
func (r *CompanyRepository) GetByID(id uuid.UUID) (*models.Company, error) {
	query := `
		SELECT id, name, description, employee_count, registered, type, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	var company models.Company
	var description sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&description,
		&company.EmployeeCount,
		&company.Registered,
		&company.Type,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, fmt.Errorf("error getting company: %w", err)
	}

	if description.Valid {
		company.Description = &description.String
	}

	return &company, nil
}

// Update updates an existing company
func (r *CompanyRepository) Update(company models.Company) error {
	query := `
		UPDATE companies
		SET name = $2, description = $3, employee_count = $4, registered = $5, type = $6, updated_at = $7
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		company.ID,
		company.Name,
		company.Description,
		company.EmployeeCount,
		company.Registered,
		company.Type,
		company.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("company not found")
	}

	return nil
}

// Delete removes a company by its ID
func (r *CompanyRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM companies WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("company not found")
	}

	return nil
}

// ExistsByName checks if a company with the given name exists
func (r *CompanyRepository) ExistsByName(name string, excludeID *uuid.UUID) (bool, error) {
	query := "SELECT COUNT(*) FROM companies WHERE name = $1"
	args := []interface{}{name}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if company exists: %w", err)
	}

	return count > 0, nil
}

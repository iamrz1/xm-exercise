package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"xm-exercise/pkg/models"
)

// Database wraps a gorm.DB connection
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection based on dialect
func NewDatabase(dialect, connectionString string) (*Database, error) {
	var db *gorm.DB
	var err error

	switch dialect {
	case "postgres":
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dialect)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("could not get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Initialize models
	if err := db.AutoMigrate(&models.User{}, &models.Company{}); err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("could not get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

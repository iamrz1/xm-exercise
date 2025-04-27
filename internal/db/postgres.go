package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresDB wraps a sql.DB connection
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(databaseURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("could not initialize schema: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

// Initialize database schema
func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS companies (
		id UUID PRIMARY KEY,
		name VARCHAR(15) UNIQUE NOT NULL,
		description TEXT,
		employee_count INTEGER NOT NULL,
		registered BOOLEAN NOT NULL,
		type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	_, err := db.Exec(schema)
	return err
}

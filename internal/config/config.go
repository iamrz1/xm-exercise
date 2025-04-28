package config

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"xm-exercise/internal/utils"
)

const (
	DefaultDatabaseURL = "local.sql"
	DefaultJWTSecret   = "super-secret-x-api-key"
)

// Config holds application configuration
type Config struct {
	Port            string
	DatabaseURL     string
	DatabaseDialect string
	JWTSecret       string
	JWTExpiration   time.Duration
	KafkaBrokers    []string
	APITimeout      time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port := utils.GetEnv("PORT", "8080")

	dbURL := utils.GetEnv("DATABASE_URL", DefaultDatabaseURL)
	dbDialect := utils.GetEnv("DATABASE_DIALECT", "sqlite")
	jwtSecret := utils.GetEnv("JWT_SECRET", DefaultJWTSecret)
	jwtExpirationStr := utils.GetEnv("JWT_EXPIRATION_HOURS", "24")
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, errors.New("JWT_EXPIRATION_HOURS must be a valid integer")
	}
	if jwtExpiration == 0 {

	}

	kafkaBrokersStr := utils.GetEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	apiTimeoutStr := utils.GetEnv("API_TIMEOUT_SECONDS", "30")
	apiTimeout, err := strconv.Atoi(apiTimeoutStr)
	if err != nil {
		return nil, errors.New("API_TIMEOUT_SECONDS must be a valid integer")
	}

	return &Config{
		Port:            port,
		DatabaseURL:     dbURL,
		DatabaseDialect: dbDialect,
		JWTSecret:       jwtSecret,
		JWTExpiration:   time.Duration(jwtExpiration) * time.Hour,
		KafkaBrokers:    kafkaBrokers,
		APITimeout:      time.Duration(apiTimeout) * time.Second,
	}, nil
}

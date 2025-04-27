package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration time.Duration
	KafkaBrokers  []string
	APITimeout    time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port := getEnv("PORT", "8080")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	jwtExpirationStr := getEnv("JWT_EXPIRATION_HOURS", "24")
	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		return nil, errors.New("JWT_EXPIRATION_HOURS must be a valid integer")
	}

	kafkaBrokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	apiTimeoutStr := getEnv("API_TIMEOUT_SECONDS", "30")
	apiTimeout, err := strconv.Atoi(apiTimeoutStr)
	if err != nil {
		return nil, errors.New("API_TIMEOUT_SECONDS must be a valid integer")
	}

	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		JWTExpiration: time.Duration(jwtExpiration) * time.Hour,
		KafkaBrokers:  kafkaBrokers,
		APITimeout:    time.Duration(apiTimeout) * time.Second,
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

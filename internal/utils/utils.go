package utils

import (
	"os"
	"regexp"
)

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// IsValidEmail returns true if the input is a valid email address, returns false otherwise
func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regex, err := regexp.Compile(emailRegex)
	if err != nil {
		return false
	}

	return regex.MatchString(email)
}

package utils

import (
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	regex := regexp.MustCompile(emailRegex)
	return regex.MatchString(email)
}

// ExtractIDFromPath returns the resource ID from URL path
func ExtractIDFromPath(r *http.Request) string {
	re := regexp.MustCompile(`([^/]+)$`)
	match := re.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString returns a string of n random characters from the alphanumeric set.
func GenerateRandomString(n int) string {
	if n <= 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.Grow(n)

	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(alphanumeric))
		sb.WriteByte(alphanumeric[randomIndex])
	}

	return sb.String()
}

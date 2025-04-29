package utils

import (
	"net/http"
	"testing"
)

func TestExtractIDFromPath(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Valid path with ID",
			path:     "/api/users/123",
			expected: "123",
		},
		{
			name:     "Valid path with ID and dashes",
			path:     "/products/abc-456",
			expected: "abc-456",
		},
		{
			name:     "Valid path with ID and extra slashes",
			path:     "/items/details/789",
			expected: "789",
		},
		{
			name:     "Empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "Root path",
			path:     "/",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new http.Request for each test case.  Crucially, use
			// a non-nil URL, and set the Path.
			req, err := http.NewRequest("GET", "http://example.com"+tc.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.URL.Path = tc.path // Set the path

			actual := ExtractIDFromPath(req)
			if actual != tc.expected {
				t.Errorf("ExtractIDFromPath(%q) = %q; expected %q", tc.path, actual, tc.expected)
			}
		})
	}
}

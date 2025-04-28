package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(w http.ResponseWriter, data interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := Response{
		Success: true,
		Data:    data,
	}
	
	if len(message) > 0 {
		response.Message = message[0]
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondError sends an error JSON response
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondCreated sends a 201 Created response
func RespondCreated(w http.ResponseWriter, data interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	response := Response{
		Success: true,
		Data:    data,
	}
	
	if len(message) > 0 {
		response.Message = message[0]
	}
	
	json.NewEncoder(w).Encode(response)
}

// EnableCORS sets CORS headers for the response
func EnableCORS(w http.ResponseWriter, r *http.Request) {
	// In production, restrict to specific origins
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict in production
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")
	
	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

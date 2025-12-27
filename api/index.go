package handler

import (
	"encoding/json"
	"net/http"
)

// Handler for the root API endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "LMS Backend API is running",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"auth": []string{
				"/api/auth/register",
				"/api/auth/profile",
				"/api/auth/update",
				"/api/auth/set-role",
			},
			"courses": []string{
				"/api/courses/create",
				"/api/courses/list",
				"/api/courses/get",
				"/api/courses/update",
				"/api/courses/delete",
				"/api/courses/enroll",
				"/api/courses/my-enrollments",
			},
			"quizzes": []string{
				"/api/quizzes/create",
				"/api/quizzes/list",
				"/api/quizzes/get",
				"/api/quizzes/add-question",
				"/api/quizzes/start",
				"/api/quizzes/submit",
				"/api/quizzes/results",
				"/api/quizzes/resume",
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

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
			"auth": map[string]string{
				"register":   "/api/auth/register",
				"profile":    "/api/auth/profile",
				"update":     "/api/auth/update",
				"set-role":   "/api/auth/set-role",
			},
			"courses": map[string]string{
				"create":         "/api/courses/create",
				"list":           "/api/courses/list",
				"get":            "/api/courses/get",
				"update":         "/api/courses/update",
				"delete":         "/api/courses/delete",
				"enroll":         "/api/courses/enroll",
				"my-enrollments": "/api/courses/my-enrollments",
			},
			"quizzes": map[string]string{
				"create":       "/api/quizzes/create",
				"list":         "/api/quizzes/list",
				"get":          "/api/quizzes/get",
				"add-question": "/api/quizzes/add-question",
				"start":        "/api/quizzes/start",
				"submit":       "/api/quizzes/submit",
				"results":      "/api/quizzes/results",
				"resume":       "/api/quizzes/resume",
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

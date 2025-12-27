package handler

import (
	"net/http"
	"strings"

	authHandlers "github.com/Ravikiran27/GOLANG_SmartEdu-LMS/api/auth"
)

// Handler routes all auth-related requests
func AuthRouter(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /api/auth/
	path := strings.TrimPrefix(r.URL.Path, "/api/auth/")
	path = strings.TrimPrefix(path, "auth/") // Handle both /api/auth and /api/auth/auth
	
	// Route to appropriate handler based on path
	switch path {
	case "register":
		authHandlers.Register(w, r)
	case "profile":
		authHandlers.GetProfile(w, r)
	case "update":
		authHandlers.UpdateProfile(w, r)
	case "set-role":
		authHandlers.SetRole(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

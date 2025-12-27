package handler

import (
	"context"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"cloud.google.com/go/firestore"
)

// SetRole assigns or updates a user's role (Admin only)
func SetRole(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Only admins can set roles
	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()

		// Parse request
		var req struct {
			UID  string `json:"uid"`
			Role string `json:"role"`
		}
		if err := utils.ParseJSONBody(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.UID == "" || req.Role == "" {
			utils.RespondError(w, http.StatusBadRequest, "UID and role are required")
			return
		}

		// Validate role
		if req.Role != "admin" && req.Role != "teacher" && req.Role != "student" {
			utils.RespondError(w, http.StatusBadRequest, "Invalid role")
			return
		}

		// Get Auth client
		authClient, err := utils.GetAuthClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize auth client")
			return
		}

		// Set custom claims
		claims := map[string]interface{}{
			"role": req.Role,
		}
		if err := authClient.SetCustomUserClaims(ctx, req.UID, claims); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to set role")
			return
		}

		// Update Firestore
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		_, err = firestoreClient.Collection("users").Doc(req.UID).Update(ctx, []firestore.Update{
			{Path: "role", Value: req.Role},
			{Path: "updatedAt", Value: utils.GetCurrentTimestamp()},
		})
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to update user role in database")
			return
		}

		utils.RespondSuccess(w, map[string]string{
			"uid":  req.UID,
			"role": req.Role,
		}, "Role updated successfully")
	}, "admin")(w, r)
}

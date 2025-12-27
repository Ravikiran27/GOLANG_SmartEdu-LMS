package handler

import (
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"cloud.google.com/go/firestore"
)

// UpdateProfile updates the authenticated user's profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Apply auth middleware
	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, _ := utils.GetUserFromContext(ctx)

		// Parse request body
		var req models.UpdateUserRequest
		if err := utils.ParseJSONBody(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		// Build update fields
		updates := []firestore.Update{
			{Path: "updatedAt", Value: utils.GetCurrentTimestamp()},
		}

		if req.DisplayName != "" {
			updates = append(updates, firestore.Update{Path: "displayName", Value: req.DisplayName})
		}
		if req.PhotoURL != "" {
			updates = append(updates, firestore.Update{Path: "photoURL", Value: req.PhotoURL})
		}
		if req.Department != "" {
			updates = append(updates, firestore.Update{Path: "metadata.department", Value: req.Department})
		}
		if req.RollNumber != "" {
			updates = append(updates, firestore.Update{Path: "metadata.rollNumber", Value: req.RollNumber})
		}
		if req.EmployeeID != "" {
			updates = append(updates, firestore.Update{Path: "metadata.employeeId", Value: req.EmployeeID})
		}

		// Update user document
		_, err = firestoreClient.Collection("users").Doc(uid).Update(ctx, updates)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to update profile")
			return
		}

		// Fetch updated user
		doc, err := firestoreClient.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch updated profile")
			return
		}

		var user models.User
		doc.DataTo(&user)

		utils.RespondSuccess(w, user, "Profile updated successfully")
	})(w, r)
}

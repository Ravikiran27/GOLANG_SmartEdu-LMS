package handler

import (
	"context"
	"lms-platform/models"
	"lms-platform/utils"
	"net/http"

	"cloud.google.com/go/firestore"
)

// GetProfile retrieves the authenticated user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Apply auth middleware (no role restriction)
	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, _ := utils.GetUserFromContext(ctx)

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		// Get user document
		doc, err := firestoreClient.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "User not found")
			return
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse user data")
			return
		}

		// Update last login
		firestoreClient.Collection("users").Doc(uid).Update(ctx, []firestore.Update{
			{Path: "metadata.lastLogin", Value: utils.GetCurrentTimestamp()},
		})

		utils.RespondSuccess(w, user)
	})(w, r)
}

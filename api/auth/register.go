package handler

import (
	"context"
	"lms-platform/models"
	"lms-platform/utils"
	"net/http"

	"firebase.google.com/go/v4/auth"
)

// Register creates a new user with Firebase Auth and Firestore
func Register(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := context.Background()

	// Parse request body
	var req models.CreateUserRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.DisplayName == "" || req.Role == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email, password, displayName, and role are required")
		return
	}

	// Validate role
	if req.Role != "admin" && req.Role != "teacher" && req.Role != "student" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid role. Must be admin, teacher, or student")
		return
	}

	// Get Firebase Auth client
	authClient, err := utils.GetAuthClient(ctx)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize auth client")
		return
	}

	// Create user in Firebase Auth
	userParams := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.DisplayName).
		EmailVerified(false).
		Disabled(false)

	firebaseUser, err := authClient.CreateUser(ctx, userParams)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Failed to create user: "+err.Error())
		return
	}

	// Set custom claims for role-based access
	claims := map[string]interface{}{
		"role": req.Role,
	}
	if err := authClient.SetCustomUserClaims(ctx, firebaseUser.UID, claims); err != nil {
		// Rollback: delete user if claim setting fails
		authClient.DeleteUser(ctx, firebaseUser.UID)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to set user role")
		return
	}

	// Get Firestore client
	firestoreClient, err := utils.GetFirestoreClient(ctx)
	if err != nil {
		authClient.DeleteUser(ctx, firebaseUser.UID)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
		return
	}

	// Create user document in Firestore
	now := utils.GetCurrentTimestamp()
	user := models.User{
		UID:         firebaseUser.UID,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Role:        req.Role,
		PhotoURL:    "",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata: models.UserMetadata{
			LastLogin:  now,
			Department: req.Department,
			RollNumber: req.RollNumber,
			EmployeeID: req.EmployeeID,
		},
	}

	_, err = firestoreClient.Collection("users").Doc(firebaseUser.UID).Set(ctx, user)
	if err != nil {
		// Rollback: delete Firebase Auth user
		authClient.DeleteUser(ctx, firebaseUser.UID)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create user document")
		return
	}

	utils.RespondCreated(w, map[string]interface{}{
		"uid":   firebaseUser.UID,
		"email": req.Email,
		"role":  req.Role,
	}, "User registered successfully")
}

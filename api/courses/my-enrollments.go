package handler

import (
	"context"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"google.golang.org/api/iterator"
)

// GetMyEnrollments retrieves student's enrollments
func GetMyEnrollments(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, _ := utils.GetUserFromContext(ctx)

		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		// Get enrollments
		iter := firestoreClient.Collection("enrollments").
			Where("studentId", "==", uid).
			OrderBy("enrolledAt", firestore.Desc).
			Documents(ctx)
		defer iter.Stop()

		var enrollments []models.Enrollment
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch enrollments")
				return
			}

			var enrollment models.Enrollment
			doc.DataTo(&enrollment)
			enrollments = append(enrollments, enrollment)
		}

		utils.RespondSuccess(w, enrollments)
	}, "student")(w, r)
}

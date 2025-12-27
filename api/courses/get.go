package handler

import (
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"
)

// GetCourse retrieves a single course by ID
func GetCourse(w http.ResponseWriter, r *http.Request) {
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
		_, _, role := utils.GetUserFromContext(ctx)

		// Get course ID from query params
		courseID := r.URL.Query().Get("id")
		if courseID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Course ID is required")
			return
		}

		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		// Get course document
		doc, err := firestoreClient.Collection("courses").Doc(courseID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Course not found")
			return
		}

		var course models.Course
		doc.DataTo(&course)

		// Check if course is deleted
		if course.IsDeleted {
			utils.RespondError(w, http.StatusNotFound, "Course not found")
			return
		}

		// Students can only view published courses
		if role == "student" && !course.IsPublished {
			utils.RespondError(w, http.StatusForbidden, "Course not accessible")
			return
		}

		utils.RespondSuccess(w, course)
	})(w, r)
}

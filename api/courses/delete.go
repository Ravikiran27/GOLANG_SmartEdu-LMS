package handler

import (
	"context"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"cloud.google.com/go/firestore"
)

// DeleteCourse soft deletes a course (Teacher/Admin only)
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, role := utils.GetUserFromContext(ctx)

		// Get course ID
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

		// Get existing course
		doc, err := firestoreClient.Collection("courses").Doc(courseID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Course not found")
			return
		}

		var course models.Course
		doc.DataTo(&course)

		// Authorization: teacher can only delete own courses
		if role == "teacher" && course.TeacherID != uid {
			utils.RespondError(w, http.StatusForbidden, "You can only delete your own courses")
			return
		}

		// Soft delete
		_, err = firestoreClient.Collection("courses").Doc(courseID).Update(ctx, []firestore.Update{
			{Path: "isDeleted", Value: true},
			{Path: "updatedAt", Value: utils.GetCurrentTimestamp()},
		})
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to delete course")
			return
		}

		utils.RespondSuccess(w, map[string]string{"courseId": courseID}, "Course deleted successfully")
	}, "teacher", "admin")(w, r)
}

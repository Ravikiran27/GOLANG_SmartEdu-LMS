package handler

import (
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"cloud.google.com/go/firestore"
)

// UpdateCourse updates an existing course (Teacher/Admin only)
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, role := utils.GetUserFromContext(ctx)

		// Get course ID from query
		courseID := r.URL.Query().Get("id")
		if courseID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Course ID is required")
			return
		}

		// Parse request
		var req models.UpdateCourseRequest
		if err := utils.ParseJSONBody(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
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

		// Authorization check: teacher can only update own courses
		if role == "teacher" && course.TeacherID != uid {
			utils.RespondError(w, http.StatusForbidden, "You can only update your own courses")
			return
		}

		// Build updates
		updates := []firestore.Update{
			{Path: "updatedAt", Value: utils.GetCurrentTimestamp()},
		}

		if req.Title != "" {
			updates = append(updates, firestore.Update{Path: "title", Value: req.Title})
		}
		if req.Description != "" {
			updates = append(updates, firestore.Update{Path: "description", Value: req.Description})
		}
		if req.Syllabus != "" {
			updates = append(updates, firestore.Update{Path: "syllabus", Value: req.Syllabus})
		}
		if req.Category != "" {
			updates = append(updates, firestore.Update{Path: "category", Value: req.Category})
		}
		if req.Difficulty != "" {
			updates = append(updates, firestore.Update{Path: "difficulty", Value: req.Difficulty})
		}
		if req.Thumbnail != "" {
			updates = append(updates, firestore.Update{Path: "thumbnail", Value: req.Thumbnail})
		}
		if req.Materials != nil {
			updates = append(updates, firestore.Update{Path: "materials", Value: req.Materials})
		}
		if req.IsPublished {
			updates = append(updates, firestore.Update{Path: "isPublished", Value: req.IsPublished})
		}

		// Update document
		_, err = firestoreClient.Collection("courses").Doc(courseID).Update(ctx, updates)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to update course")
			return
		}

		// Fetch updated course
		updatedDoc, _ := firestoreClient.Collection("courses").Doc(courseID).Get(ctx)
		var updatedCourse models.Course
		updatedDoc.DataTo(&updatedCourse)

		utils.RespondSuccess(w, updatedCourse, "Course updated successfully")
	}, "teacher", "admin")(w, r)
}

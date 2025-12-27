package handler

import (
	"context"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
	"net/http"

	"github.com/google/uuid"
)

// CreateCourse creates a new course (Teacher/Admin only)
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, role := utils.GetUserFromContext(ctx)

		// Parse request
		var req models.CreateCourseRequest
		if err := utils.ParseJSONBody(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate required fields
		if req.Title == "" || req.Description == "" || req.Category == "" || req.Difficulty == "" {
			utils.RespondError(w, http.StatusBadRequest, "Title, description, category, and difficulty are required")
			return
		}

		// Get user info for teacher name
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		userDoc, err := firestoreClient.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user info")
			return
		}

		var user models.User
		userDoc.DataTo(&user)

		// Create course document
		now := utils.GetCurrentTimestamp()
		course := models.Course{
			CourseID:        uuid.New().String(),
			Title:           req.Title,
			Description:     req.Description,
			Syllabus:        req.Syllabus,
			TeacherID:       uid,
			TeacherName:     user.DisplayName,
			Category:        req.Category,
			Difficulty:      req.Difficulty,
			Thumbnail:       req.Thumbnail,
			Materials:       []models.CourseMaterial{},
			EnrollmentCount: 0,
			IsPublished:     false,
			CreatedAt:       now,
			UpdatedAt:       now,
			IsDeleted:       false,
		}

		_, err = firestoreClient.Collection("courses").Doc(course.CourseID).Set(ctx, course)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create course")
			return
		}

		utils.RespondCreated(w, course, "Course created successfully")
	}, "teacher", "admin")(w, r)
}

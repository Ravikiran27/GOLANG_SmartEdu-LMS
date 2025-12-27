package handler

import (
	"context"
	"lms-platform/models"
	"lms-platform/utils"
	"net/http"

	"github.com/google/uuid"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// EnrollCourse enrolls a student in a course
func EnrollCourse(w http.ResponseWriter, r *http.Request) {
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
		uid, _, _ := utils.GetUserFromContext(ctx)

		// Parse request
		var req models.EnrollmentRequest
		if err := utils.ParseJSONBody(r, &req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.CourseID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Course ID is required")
			return
		}

		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		// Check if course exists and is published
		courseDoc, err := firestoreClient.Collection("courses").Doc(req.CourseID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Course not found")
			return
		}

		var course models.Course
		courseDoc.DataTo(&course)

		if !course.IsPublished {
			utils.RespondError(w, http.StatusBadRequest, "Course is not published")
			return
		}

		// Check if already enrolled
		existingEnrollment := firestoreClient.Collection("enrollments").
			Where("studentId", "==", uid).
			Where("courseId", "==", req.CourseID).
			Limit(1).
			Documents(ctx)
		
		if doc, err := existingEnrollment.Next(); err == nil && doc != nil {
			utils.RespondError(w, http.StatusConflict, "Already enrolled in this course")
			return
		}

		// Get student info
		userDoc, err := firestoreClient.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user info")
			return
		}

		var user models.User
		userDoc.DataTo(&user)

		// Create enrollment
		now := utils.GetCurrentTimestamp()
		enrollment := models.Enrollment{
			EnrollmentID:       uuid.New().String(),
			StudentID:          uid,
			StudentName:        user.DisplayName,
			CourseID:           req.CourseID,
			CourseTitle:        course.Title,
			EnrolledAt:         now,
			Progress:           0,
			CompletedMaterials: []string{},
			Status:             "active",
			LastAccessedAt:     now,
		}

		_, err = firestoreClient.Collection("enrollments").Doc(enrollment.EnrollmentID).Set(ctx, enrollment)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create enrollment")
			return
		}

		// Increment course enrollment count
		firestoreClient.Collection("courses").Doc(req.CourseID).Update(ctx, []firestore.Update{
			{Path: "enrollmentCount", Value: firestore.Increment(1)},
		})

		utils.RespondCreated(w, enrollment, "Enrolled successfully")
	}, "student")(w, r)
}

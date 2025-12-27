package handler

import (
	"context"
	"net/http"

	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
)

// Handler gets a single quiz by ID
func GetQuiz(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Authenticate
	utils.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value("uid").(string)
		role := ctx.Value("role").(string)

		// Get quiz ID from query
		quizID := r.URL.Query().Get("id")
		if quizID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Quiz ID is required")
			return
		}

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Get quiz
		quizDoc, err := firestoreClient.Collection("quizzes").Doc(quizID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Quiz not found")
			return
		}

		var quiz models.Quiz
		if err := quizDoc.DataTo(&quiz); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse quiz data")
			return
		}
		quiz.ID = quizDoc.Ref.ID

		// Authorization checks
		if role == "student" {
			// Students can only see published quizzes
			if !quiz.IsPublished {
				utils.RespondError(w, http.StatusForbidden, "This quiz is not published")
				return
			}

			// Verify student is enrolled in the course
			enrollmentQuery := firestoreClient.Collection("enrollments").
				Where("studentId", "==", userID).
				Where("courseId", "==", quiz.CourseID).
				Where("status", "==", "active")

			enrollDocs, err := enrollmentQuery.Documents(ctx).GetAll()
			if err != nil || len(enrollDocs) == 0 {
				utils.RespondError(w, http.StatusForbidden, "You must be enrolled in this course")
				return
			}
		} else if role == "teacher" {
			// Teachers can only see their own quizzes
			if quiz.TeacherID != userID {
				utils.RespondError(w, http.StatusForbidden, "You can only view your own quizzes")
				return
			}
		}
		// Admins can see all quizzes

		utils.RespondSuccess(w, "Quiz fetched successfully", quiz)
	})).ServeHTTP(w, r)
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"cloud.google.com/go/firestore"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
)

// Handler creates a new quiz (teacher/admin only)
func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Authenticate and authorize (teacher/admin only)
	utils.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value("uid").(string)
		role := ctx.Value("role").(string)

		// Only teachers and admins can create quizzes
		if role != "teacher" && role != "admin" {
			utils.RespondError(w, http.StatusForbidden, "Only teachers and admins can create quizzes")
			return
		}

		// Parse request body
		var req models.CreateQuizRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate required fields
		if req.Title == "" {
			utils.RespondError(w, http.StatusBadRequest, "Title is required")
			return
		}
		if req.CourseID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Course ID is required")
			return
		}
		if req.TotalMarks <= 0 {
			utils.RespondError(w, http.StatusBadRequest, "Total marks must be greater than 0")
			return
		}
		if req.PassingMarks < 0 || req.PassingMarks > req.TotalMarks {
			utils.RespondError(w, http.StatusBadRequest, "Invalid passing marks")
			return
		}
		if req.Duration <= 0 {
			utils.RespondError(w, http.StatusBadRequest, "Duration must be greater than 0")
			return
		}

		// Validate deadline is in future (if provided)
		if !req.Deadline.IsZero() && req.Deadline.Before(time.Now()) {
			utils.RespondError(w, http.StatusBadRequest, "Deadline must be in the future")
			return
		}

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Verify course exists and user is the teacher or admin
		courseRef := firestoreClient.Collection("courses").Doc(req.CourseID)
		courseDoc, err := courseRef.Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Course not found")
			return
		}

		var course models.Course
		if err := courseDoc.DataTo(&course); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse course data")
			return
		}

		// Check if user is the course teacher or admin
		if role != "admin" && course.TeacherID != userID {
			utils.RespondError(w, http.StatusForbidden, "You can only create quizzes for your own courses")
			return
		}

		// Create quiz object
		now := time.Now()
		quiz := models.Quiz{
			Title:                req.Title,
			Description:          req.Description,
			CourseID:             req.CourseID,
			TeacherID:            userID,
			TotalMarks:           req.TotalMarks,
			PassingMarks:         req.PassingMarks,
			Duration:             req.Duration,
			Instructions:         req.Instructions,
			Deadline:             req.Deadline,
			ShowResultsAfterSubmit: req.ShowResultsAfterSubmit,
			ShuffleQuestions:     req.ShuffleQuestions,
			ShuffleOptions:       req.ShuffleOptions,
			MaxAttempts:          req.MaxAttempts,
			AllowReview:          req.AllowReview,
			
			// Cheating prevention features
			PreventTabSwitch:     req.PreventTabSwitch,
			MaxTabSwitches:       req.MaxTabSwitches,
			RequireFullscreen:    req.RequireFullscreen,
			DisableCopyPaste:     req.DisableCopyPaste,
			EnableProctoring:     req.EnableProctoring,
			RandomizeQuestionOrder: req.RandomizeQuestionOrder,
			TimePerQuestion:      req.TimePerQuestion,
			LockAfterSubmit:      true, // Always lock after submit
			
			// Teacher permissions
			AllowTeacherResume:   req.AllowTeacherResume,
			AllowTeacherExtendTime: req.AllowTeacherExtendTime,
			
			QuestionCount:        0,
			IsPublished:          req.IsPublished,
			CreatedAt:            now,
			UpdatedAt:            now,
		}

		// Set default values
		if quiz.MaxAttempts == 0 {
			quiz.MaxAttempts = 1 // Default to 1 attempt
		}
		if quiz.PreventTabSwitch && quiz.MaxTabSwitches == 0 {
			quiz.MaxTabSwitches = 3 // Default to 3 tab switches
		}

		// Save quiz to Firestore
		quizRef := firestoreClient.Collection("quizzes").NewDoc()
		quiz.ID = quizRef.ID

		if _, err := quizRef.Set(ctx, quiz); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create quiz")
			return
		}

		utils.RespondSuccess(w, quiz, "Quiz created successfully")
	})).ServeHTTP(w, r)
}

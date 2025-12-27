package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/api/iterator"
	"github.com/yourusername/lms/models"
	"github.com/yourusername/lms/utils"
)

// Handler starts a quiz attempt for a student
func Handler(w http.ResponseWriter, r *http.Request) {
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

	// Authenticate (students only)
	utils.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value("uid").(string)
		role := ctx.Value("role").(string)

		// Only students can start quiz attempts
		if role != "student" {
			utils.RespondError(w, http.StatusForbidden, "Only students can start quiz attempts")
			return
		}

		// Parse request body
		var req struct {
			QuizID string `json:"quizId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.QuizID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Quiz ID is required")
			return
		}

		// Get Firestore client
		firestoreClient := utils.GetFirestoreClient()
		if firestoreClient == nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Get quiz
		quizDoc, err := firestoreClient.Collection("quizzes").Doc(req.QuizID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Quiz not found")
			return
		}

		var quiz models.Quiz
		if err := quizDoc.DataTo(&quiz); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse quiz data")
			return
		}

		// Verify quiz is published
		if !quiz.IsPublished {
			utils.RespondError(w, http.StatusForbidden, "This quiz is not published")
			return
		}

		// Check deadline
		if !quiz.Deadline.IsZero() && time.Now().After(quiz.Deadline) {
			utils.RespondError(w, http.StatusForbidden, "Quiz deadline has passed")
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

		// Check previous attempts
		submissionsQuery := firestoreClient.Collection("quiz_submissions").
			Where("quizId", "==", req.QuizID).
			Where("studentId", "==", userID)

		submissionDocs, err := submissionsQuery.Documents(ctx).GetAll()
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to check previous attempts")
			return
		}

		// Check for in-progress submission
		for _, doc := range submissionDocs {
			var sub models.QuizSubmission
			if err := doc.DataTo(&sub); err == nil {
				if sub.Status == "in_progress" {
					// Resume existing attempt
					utils.RespondSuccess(w, "Resuming existing attempt", map[string]interface{}{
						"submission": sub,
						"resumed":    true,
					})
					return
				}
			}
		}

		// Count completed attempts
		completedAttempts := 0
		for _, doc := range submissionDocs {
			var sub models.QuizSubmission
			if err := doc.DataTo(&sub); err == nil {
				if sub.Status == "submitted" || sub.Status == "evaluated" {
					completedAttempts++
				}
			}
		}

		// Check max attempts
		if quiz.MaxAttempts > 0 && completedAttempts >= quiz.MaxAttempts {
			utils.RespondError(w, http.StatusForbidden, "Maximum attempts reached")
			return
		}

		// Get all questions for this quiz
		questionsQuery := firestoreClient.Collection("questions").
			Where("quizId", "==", req.QuizID).
			OrderBy("order", "asc")

		iter := questionsQuery.Documents(ctx)
		defer iter.Stop()

		questions := make([]models.Question, 0)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch questions")
				return
			}

			var question models.Question
			if err := doc.DataTo(&question); err != nil {
				continue
			}
			question.ID = doc.Ref.ID
			questions = append(questions, question)
		}

		if len(questions) == 0 {
			utils.RespondError(w, http.StatusBadRequest, "Quiz has no questions")
			return
		}

		// Shuffle questions if enabled
		if quiz.ShuffleQuestions || quiz.RandomizeQuestionOrder {
			utils.ShuffleQuestions(&questions)
		}

		// Shuffle options if enabled
		if quiz.ShuffleOptions {
			for i := range questions {
				utils.ShuffleOptions(&questions[i].Options)
			}
		}

		// Create submission
		now := time.Now()
		submission := models.QuizSubmission{
			QuizID:           req.QuizID,
			StudentID:        userID,
			CourseID:         quiz.CourseID,
			AttemptNumber:    completedAttempts + 1,
			Status:           "in_progress",
			StartedAt:        now,
			TimeLimit:        quiz.Duration,
			Questions:        questions,
			TabSwitchCount:   0,
			FullscreenExits:  0,
			SuspiciousActivity: make([]string, 0),
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		// Save submission
		submissionRef := firestoreClient.Collection("quiz_submissions").NewDoc()
		submission.ID = submissionRef.ID

		if _, err := submissionRef.Set(ctx, submission); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to start quiz")
			return
		}

		// Remove correct answers from questions before sending to client
		clientQuestions := make([]models.Question, len(questions))
		for i, q := range questions {
			clientQ := q
			// Remove correct answer info for auto-graded questions
			if q.Type == "mcq" || q.Type == "true_false" {
				clientOptions := make([]models.QuestionOption, len(q.Options))
				for j, opt := range q.Options {
					clientOptions[j] = models.QuestionOption{
						ID:   opt.ID,
						Text: opt.Text,
						// IsCorrect field is intentionally omitted
					}
				}
				clientQ.Options = clientOptions
			}
			clientQuestions[i] = clientQ
		}

		submission.Questions = clientQuestions

		utils.RespondSuccess(w, "Quiz started successfully", map[string]interface{}{
			"submission": submission,
			"resumed":    false,
			"cheatingPrevention": map[string]interface{}{
				"preventTabSwitch":  quiz.PreventTabSwitch,
				"maxTabSwitches":    quiz.MaxTabSwitches,
				"requireFullscreen": quiz.RequireFullscreen,
				"disableCopyPaste":  quiz.DisableCopyPaste,
				"enableProctoring":  quiz.EnableProctoring,
			},
		})
	})).ServeHTTP(w, r)
}

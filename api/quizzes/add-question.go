package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/yourusername/lms/models"
	"github.com/yourusername/lms/utils"
)

// Handler adds a question to a quiz (teacher/admin only)
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

	// Authenticate and authorize (teacher/admin only)
	utils.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value("uid").(string)
		role := ctx.Value("role").(string)

		// Only teachers and admins can add questions
		if role != "teacher" && role != "admin" {
			utils.RespondError(w, http.StatusForbidden, "Only teachers and admins can add questions")
			return
		}

		// Parse request body
		var req struct {
			QuizID string         `json:"quizId"`
			Question models.Question `json:"question"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate required fields
		if req.QuizID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Quiz ID is required")
			return
		}
		if req.Question.Text == "" {
			utils.RespondError(w, http.StatusBadRequest, "Question text is required")
			return
		}
		if req.Question.Type == "" {
			utils.RespondError(w, http.StatusBadRequest, "Question type is required")
			return
		}
		if req.Question.Points <= 0 {
			utils.RespondError(w, http.StatusBadRequest, "Points must be greater than 0")
			return
		}

		// Validate question type
		validTypes := map[string]bool{
			"mcq":        true,
			"true_false": true,
			"short_answer": true,
			"long_answer": true,
		}
		if !validTypes[req.Question.Type] {
			utils.RespondError(w, http.StatusBadRequest, "Invalid question type")
			return
		}

		// Validate MCQ and True/False questions
		if req.Question.Type == "mcq" || req.Question.Type == "true_false" {
			if len(req.Question.Options) == 0 {
				utils.RespondError(w, http.StatusBadRequest, "Options are required for MCQ and True/False questions")
				return
			}

			// Check if at least one option is marked as correct
			hasCorrect := false
			for _, opt := range req.Question.Options {
				if opt.IsCorrect {
					hasCorrect = true
					break
				}
			}
			if !hasCorrect {
				utils.RespondError(w, http.StatusBadRequest, "At least one option must be marked as correct")
				return
			}
		}

		// Get Firestore client
		firestoreClient := utils.GetFirestoreClient()
		if firestoreClient == nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Verify quiz exists and user is the teacher or admin
		quizRef := firestoreClient.Collection("quizzes").Doc(req.QuizID)
		quizDoc, err := quizRef.Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Quiz not found")
			return
		}

		var quiz models.Quiz
		if err := quizDoc.DataTo(&quiz); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse quiz data")
			return
		}

		// Check if user is the quiz teacher or admin
		if role != "admin" && quiz.TeacherID != userID {
			utils.RespondError(w, http.StatusForbidden, "You can only add questions to your own quizzes")
			return
		}

		// Create question object
		now := time.Now()
		question := models.Question{
			QuizID:     req.QuizID,
			Text:       req.Question.Text,
			Type:       req.Question.Type,
			Options:    req.Question.Options,
			Points:     req.Question.Points,
			Order:      quiz.QuestionCount + 1, // Auto-increment order
			Explanation: req.Question.Explanation,
			ImageURL:   req.Question.ImageURL,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Save question to Firestore
		questionRef := firestoreClient.Collection("questions").NewDoc()
		question.ID = questionRef.ID

		if _, err := questionRef.Set(ctx, question); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to add question")
			return
		}

		// Update quiz question count and total marks
		_, err = quizRef.Update(ctx, []firestore.Update{
			{Path: "questionCount", Value: firestore.Increment(1)},
			{Path: "totalMarks", Value: firestore.Increment(req.Question.Points)},
			{Path: "updatedAt", Value: now},
		})
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to update quiz")
			return
		}

		utils.RespondSuccess(w, "Question added successfully", question)
	})).ServeHTTP(w, r)
}

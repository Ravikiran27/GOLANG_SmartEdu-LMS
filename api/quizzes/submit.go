package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
)

// Handler submits quiz answers and auto-evaluates
func SubmitQuiz(w http.ResponseWriter, r *http.Request) {
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

		// Only students can submit quizzes
		if role != "student" {
			utils.RespondError(w, http.StatusForbidden, "Only students can submit quizzes")
			return
		}

		// Parse request body
		var req models.SubmitQuizRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.SubmissionID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Submission ID is required")
			return
		}

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Get submission
		submissionRef := firestoreClient.Collection("quiz_submissions").Doc(req.SubmissionID)
		submissionDoc, err := submissionRef.Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Submission not found")
			return
		}

		var submission models.QuizSubmission
		if err := submissionDoc.DataTo(&submission); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse submission data")
			return
		}

		// Verify ownership
		if submission.StudentID != userID {
			utils.RespondError(w, http.StatusForbidden, "You can only submit your own quiz")
			return
		}

		// Check if already submitted
		if submission.Status == "submitted" || submission.Status == "evaluated" {
			utils.RespondError(w, http.StatusBadRequest, "Quiz already submitted")
			return
		}

		// Get quiz details
		quizDoc, err := firestoreClient.Collection("quizzes").Doc(submission.QuizID).Get(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "Quiz not found")
			return
		}

		var quiz models.Quiz
		if err := quizDoc.DataTo(&quiz); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to parse quiz data")
			return
		}

		// Check time limit
		now := time.Now()
		timeElapsed := now.Sub(submission.StartedAt)
		if timeElapsed > time.Duration(quiz.Duration)*time.Minute {
			// Auto-submit if time expired
			req.TimedOut = true
		}

		// Get original questions with correct answers
		originalQuestions := make(map[string]models.Question)
		questionsQuery := firestoreClient.Collection("questions").
			Where("quizId", "==", submission.QuizID)

		questionDocs, err := questionsQuery.Documents(ctx).GetAll()
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch questions")
			return
		}

		for _, doc := range questionDocs {
			var q models.Question
			if err := doc.DataTo(&q); err == nil {
				q.ID = doc.Ref.ID
				originalQuestions[q.ID] = q
			}
		}

		// Auto-evaluate answers
		totalScore := 0.0
		evaluatedAnswers := make([]models.Answer, 0)

		for _, answer := range req.Answers {
			originalQ, exists := originalQuestions[answer.QuestionID]
			if !exists {
				continue
			}

			evaluatedAnswer := answer
			evaluatedAnswer.PointsAwarded = 0

			// Auto-grade MCQ and True/False
			if originalQ.Type == "mcq" || originalQ.Type == "true_false" {
				// Find correct option(s)
				correctIDs := make(map[string]bool)
				for _, opt := range originalQ.Options {
					if opt.IsCorrect {
						correctIDs[opt.ID] = true
					}
				}

				// Check if student's answer is correct
				allCorrect := true
				for _, selectedID := range answer.SelectedOptions {
					if !correctIDs[selectedID] {
						allCorrect = false
						break
					}
				}

				// Also check if all correct options were selected
				if allCorrect && len(answer.SelectedOptions) == len(correctIDs) {
					evaluatedAnswer.IsCorrect = true
					evaluatedAnswer.PointsAwarded = originalQ.Points
					totalScore += originalQ.Points
				} else {
					evaluatedAnswer.IsCorrect = false
				}
			}
			// Short and long answers need manual grading
			// Leave them at 0 points for now

			evaluatedAnswers = append(evaluatedAnswers, evaluatedAnswer)
		}

		// Calculate percentage
		percentage := 0.0
		if quiz.TotalMarks > 0 {
			percentage = (totalScore / float64(quiz.TotalMarks)) * 100
		}

		// Determine pass/fail
		passed := totalScore >= float64(quiz.PassingMarks)

		// Check for suspicious activity
		suspiciousFlags := make([]string, 0)
		if quiz.PreventTabSwitch && req.TabSwitches > quiz.MaxTabSwitches {
			suspiciousFlags = append(suspiciousFlags, "Excessive tab switching detected")
		}
		if quiz.RequireFullscreen && req.FullscreenExits > 0 {
			suspiciousFlags = append(suspiciousFlags, "Exited fullscreen mode")
		}
		if req.TimedOut {
			suspiciousFlags = append(suspiciousFlags, "Time limit exceeded")
		}

		// Update submission
		updates := []firestore.Update{
			{Path: "answers", Value: evaluatedAnswers},
			{Path: "status", Value: "evaluated"},
			{Path: "submittedAt", Value: now},
			{Path: "score", Value: totalScore},
			{Path: "percentage", Value: percentage},
			{Path: "passed", Value: passed},
			{Path: "timeTaken", Value: int(timeElapsed.Minutes())},
			{Path: "tabSwitchCount", Value: req.TabSwitches},
			{Path: "fullscreenExits", Value: req.FullscreenExits},
			{Path: "suspiciousActivity", Value: suspiciousFlags},
			{Path: "updatedAt", Value: now},
		}

		if _, err := submissionRef.Update(ctx, updates); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to submit quiz")
			return
		}

		// Update student's analytics
		analyticsQuery := firestoreClient.Collection("analytics").
			Where("studentId", "==", userID)

		analyticsDocs, err := analyticsQuery.Documents(ctx).GetAll()
		
		var analyticsRef *firestore.DocumentRef
		if err == nil && len(analyticsDocs) > 0 {
			// Update existing analytics
			analyticsRef = analyticsDocs[0].Ref
			_, err = analyticsRef.Update(ctx, []firestore.Update{
				{Path: "quizzesCompleted", Value: firestore.Increment(1)},
				{Path: "totalQuizScore", Value: firestore.Increment(int(totalScore))},
				{Path: "updatedAt", Value: now},
			})
		} else {
			// Create new analytics
			analytics := models.Analytics{
				StudentID:        userID,
				QuizzesCompleted: 1,
				TotalQuizScore:   int(totalScore),
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			analyticsRef = firestoreClient.Collection("analytics").NewDoc()
			_, err = analyticsRef.Set(ctx, analytics)
		}

		// Prepare response
		response := map[string]interface{}{
			"submissionId": req.SubmissionID,
			"score":        totalScore,
			"totalMarks":   quiz.TotalMarks,
			"percentage":   percentage,
			"passed":       passed,
			"timeTaken":    int(timeElapsed.Minutes()),
		}

		// Include results if enabled
		if quiz.ShowResultsAfterSubmit {
			response["answers"] = evaluatedAnswers
			response["suspiciousActivity"] = suspiciousFlags
		}

		utils.RespondSuccess(w, response, "Quiz submitted successfully")
	})).ServeHTTP(w, r)
}

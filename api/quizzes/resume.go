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

// Handler allows teacher to resume a student's locked quiz (teacher only)
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

	// Authenticate (teacher/admin only)
	utils.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value("uid").(string)
		role := ctx.Value("role").(string)

		// Only teachers and admins can resume quizzes
		if role != "teacher" && role != "admin" {
			utils.RespondError(w, http.StatusForbidden, "Only teachers and admins can resume quizzes")
			return
		}

		// Parse request body
		var req struct {
			SubmissionID string `json:"submissionId"`
			ExtendTime   int    `json:"extendTime"` // Additional minutes to add
			Reason       string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.SubmissionID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Submission ID is required")
			return
		}

		// Get Firestore client
		firestoreClient := utils.GetFirestoreClient()
		if firestoreClient == nil {
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

		// Get quiz
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

		// Verify teacher owns the quiz (or is admin)
		if role != "admin" && quiz.TeacherID != userID {
			utils.RespondError(w, http.StatusForbidden, "You can only resume quizzes for your own courses")
			return
		}

		// Check if quiz allows teacher resume
		if !quiz.AllowTeacherResume {
			utils.RespondError(w, http.StatusForbidden, "This quiz does not allow teacher resume")
			return
		}

		// Check if submission is in a resumable state
		if submission.Status != "submitted" && submission.Status != "evaluated" {
			utils.RespondError(w, http.StatusBadRequest, "Only completed submissions can be resumed")
			return
		}

		// Create resume record
		now := time.Now()
		updates := []firestore.Update{
			{Path: "status", Value: "in_progress"},
			{Path: "resumedBy", Value: userID},
			{Path: "resumedAt", Value: now},
			{Path: "resumeReason", Value: req.Reason},
			{Path: "updatedAt", Value: now},
		}

		// Extend time if requested and allowed
		if req.ExtendTime > 0 {
			if !quiz.AllowTeacherExtendTime {
				utils.RespondError(w, http.StatusForbidden, "This quiz does not allow time extension")
				return
			}
			updates = append(updates, firestore.Update{
				Path:  "timeLimit",
				Value: submission.TimeLimit + req.ExtendTime,
			})
		}

		// Update submission
		if _, err := submissionRef.Update(ctx, updates); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to resume quiz")
			return
		}

		// Log the resume action
		logEntry := map[string]interface{}{
			"action":       "quiz_resumed",
			"submissionId": req.SubmissionID,
			"quizId":       submission.QuizID,
			"studentId":    submission.StudentID,
			"teacherId":    userID,
			"reason":       req.Reason,
			"extendTime":   req.ExtendTime,
			"timestamp":    now,
		}
		firestoreClient.Collection("audit_logs").NewDoc().Set(ctx, logEntry)

		// Get student details for notification
		studentDoc, err := firestoreClient.Collection("users").Doc(submission.StudentID).Get(ctx)
		if err == nil {
			var student models.User
			if err := studentDoc.DataTo(&student); err == nil {
				// Create notification for student
				notification := models.Notification{
					UserID:    submission.StudentID,
					Type:      "quiz_resumed",
					Title:     "Quiz Resumed",
					Message:   "Your teacher has resumed your quiz: " + quiz.Title,
					Link:      "/quizzes/" + submission.QuizID,
					IsRead:    false,
					CreatedAt: now,
				}
				firestoreClient.Collection("notifications").NewDoc().Set(ctx, notification)
			}
		}

		utils.RespondSuccess(w, "Quiz resumed successfully", map[string]interface{}{
			"submissionId": req.SubmissionID,
			"extendedTime": req.ExtendTime,
			"newTimeLimit": submission.TimeLimit + req.ExtendTime,
		})
	})).ServeHTTP(w, r)
}

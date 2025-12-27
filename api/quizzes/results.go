package handler

import (
	"context"
	"net/http"

	"google.golang.org/api/iterator"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/models"
	"github.com/Ravikiran27/GOLANG_SmartEdu-LMS/utils"
)

// Handler gets quiz results for a student or all students (for teachers)
func GetResults(w http.ResponseWriter, r *http.Request) {
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

		// Get query parameters
		quizID := r.URL.Query().Get("quizId")
		submissionID := r.URL.Query().Get("submissionId")

		if quizID == "" && submissionID == "" {
			utils.RespondError(w, http.StatusBadRequest, "Either quizId or submissionId is required")
			return
		}

		// Get Firestore client
		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// If specific submission requested
		if submissionID != "" {
			submissionDoc, err := firestoreClient.Collection("quiz_submissions").Doc(submissionID).Get(ctx)
			if err != nil {
				utils.RespondError(w, http.StatusNotFound, "Submission not found")
				return
			}

			var submission models.QuizSubmission
			if err := submissionDoc.DataTo(&submission); err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to parse submission data")
				return
			}
			submission.ID = submissionDoc.Ref.ID

			// Authorization check
			if role == "student" && submission.StudentID != userID {
				utils.RespondError(w, http.StatusForbidden, "You can only view your own results")
				return
			}

			if role == "teacher" {
				// Verify teacher owns the quiz
				quizDoc, err := firestoreClient.Collection("quizzes").Doc(submission.QuizID).Get(ctx)
				if err != nil {
					utils.RespondError(w, http.StatusForbidden, "Access denied")
					return
				}

				var quiz models.Quiz
				if err := quizDoc.DataTo(&quiz); err != nil {
					utils.RespondError(w, http.StatusInternalServerError, "Failed to parse quiz data")
					return
				}

				if quiz.TeacherID != userID {
					utils.RespondError(w, http.StatusForbidden, "You can only view results for your own quizzes")
					return
				}
			}

			// Get student details
			studentDoc, err := firestoreClient.Collection("users").Doc(submission.StudentID).Get(ctx)
			if err == nil {
				var student models.User
				if err := studentDoc.DataTo(&student); err == nil {
					submission.StudentName = student.Name
					submission.StudentEmail = student.Email
				}
			}

			utils.RespondSuccess(w, "Results fetched successfully", submission)
			return
		}

		// If quiz results requested
		if quizID != "" {
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

			// Build query based on role
			query := firestoreClient.Collection("quiz_submissions").
				Where("quizId", "==", quizID).
				Where("status", "in", []string{"submitted", "evaluated"})

			if role == "student" {
				// Students see only their own submissions
				query = query.Where("studentId", "==", userID)
			} else if role == "teacher" {
				// Teachers can see all submissions for their quizzes
				if quiz.TeacherID != userID {
					utils.RespondError(w, http.StatusForbidden, "You can only view results for your own quizzes")
					return
				}
			}
			// Admins see all submissions

			// Execute query
			iter := query.OrderBy("submittedAt", "desc").Documents(ctx)
			defer iter.Stop()

			submissions := make([]models.QuizSubmission, 0)
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch results")
					return
				}

				var submission models.QuizSubmission
				if err := doc.DataTo(&submission); err != nil {
					continue
				}
				submission.ID = doc.Ref.ID

				// Get student details
				if role == "teacher" || role == "admin" {
					studentDoc, err := firestoreClient.Collection("users").Doc(submission.StudentID).Get(ctx)
					if err == nil {
						var student models.User
						if err := studentDoc.DataTo(&student); err == nil {
							submission.StudentName = student.Name
							submission.StudentEmail = student.Email
						}
					}
				}

				submissions = append(submissions, submission)
			}

			// Calculate statistics for teachers/admins
			stats := map[string]interface{}{
				"totalSubmissions": len(submissions),
			}

			if role == "teacher" || role == "admin" {
				totalScore := 0.0
				passed := 0
				for _, sub := range submissions {
					totalScore += sub.Score
					if sub.Passed {
						passed++
					}
				}

				avgScore := 0.0
				if len(submissions) > 0 {
					avgScore = totalScore / float64(len(submissions))
				}

				stats["averageScore"] = avgScore
				stats["passRate"] = float64(passed) / float64(len(submissions)) * 100
				stats["totalPassed"] = passed
				stats["totalFailed"] = len(submissions) - passed
			}

			utils.RespondSuccess(w, "Results fetched successfully", map[string]interface{}{
				"submissions": submissions,
				"statistics":  stats,
			})
		}
	})).ServeHTTP(w, r)
}

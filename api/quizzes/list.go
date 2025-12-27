package handler

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/api/iterator"
	"github.com/yourusername/lms/models"
	"github.com/yourusername/lms/utils"
)

// Handler lists quizzes (filtered by role and course)
func Handler(w http.ResponseWriter, r *http.Request) {
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
		courseID := r.URL.Query().Get("courseId")
		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
				limit = parsed
			}
		}

		// Get Firestore client
		firestoreClient := utils.GetFirestoreClient()
		if firestoreClient == nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize Firestore")
			return
		}

		// Build query based on role
		query := firestoreClient.Collection("quizzes").Query

		// Filter by course if provided
		if courseID != "" {
			query = query.Where("courseId", "==", courseID)
		}

		// Role-based filtering
		switch role {
		case "admin":
			// Admins see all quizzes
			break
		case "teacher":
			// Teachers see only their quizzes
			query = query.Where("teacherId", "==", userID)
		case "student":
			// Students see only published quizzes for courses they're enrolled in
			query = query.Where("isPublished", "==", true)
			
			// If no specific course, we need to get enrolled courses first
			if courseID == "" {
				// Get student's enrollments
				enrollmentsQuery := firestoreClient.Collection("enrollments").
					Where("studentId", "==", userID).
					Where("status", "==", "active")
				
				enrollDocs, err := enrollmentsQuery.Documents(ctx).GetAll()
				if err != nil {
					utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch enrollments")
					return
				}

				// Extract course IDs
				courseIDs := make([]string, 0)
				for _, doc := range enrollDocs {
					var enrollment models.Enrollment
					if err := doc.DataTo(&enrollment); err == nil {
						courseIDs = append(courseIDs, enrollment.CourseID)
					}
				}

				// If no enrollments, return empty array
				if len(courseIDs) == 0 {
					utils.RespondSuccess(w, "Quizzes fetched successfully", []models.Quiz{})
					return
				}

				// Filter quizzes by enrolled courses
				query = query.Where("courseId", "in", courseIDs)
			}
		default:
			utils.RespondError(w, http.StatusForbidden, "Invalid role")
			return
		}

		// Order by creation date (newest first)
		query = query.OrderBy("createdAt", "desc").Limit(limit)

		// Execute query
		iter := query.Documents(ctx)
		defer iter.Stop()

		quizzes := make([]models.Quiz, 0)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch quizzes")
				return
			}

			var quiz models.Quiz
			if err := doc.DataTo(&quiz); err != nil {
				continue
			}

			quiz.ID = doc.Ref.ID
			quizzes = append(quizzes, quiz)
		}

		utils.RespondSuccess(w, "Quizzes fetched successfully", quizzes)
	})).ServeHTTP(w, r)
}

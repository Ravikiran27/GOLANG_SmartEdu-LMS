package handler

import (
	"context"
	"lms-platform/models"
	"lms-platform/utils"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ListCourses retrieves all courses (filtered by role)
func ListCourses(w http.ResponseWriter, r *http.Request) {
	utils.EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	utils.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		ctx := r.Context()
		uid, _, role := utils.GetUserFromContext(ctx)

		firestoreClient, err := utils.GetFirestoreClient(ctx)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Failed to initialize firestore")
			return
		}

		var query firestore.Query
		coursesRef := firestoreClient.Collection("courses")

		// Filter based on role
		switch role {
		case "student":
			// Students see only published courses
			query = coursesRef.Where("isPublished", "==", true).Where("isDeleted", "==", false)
		case "teacher":
			// Teachers see their own courses + published courses
			teacherParam := r.URL.Query().Get("teacher")
			if teacherParam == "me" {
				query = coursesRef.Where("teacherId", "==", uid).Where("isDeleted", "==", false)
			} else {
				query = coursesRef.Where("isPublished", "==", true).Where("isDeleted", "==", false)
			}
		case "admin":
			// Admins see all courses
			query = coursesRef.Where("isDeleted", "==", false)
		default:
			query = coursesRef.Where("isPublished", "==", true).Where("isDeleted", "==", false)
		}

		// Pagination
		page, pageSize := utils.GetPaginationParams(r)
		query = query.OrderBy("createdAt", firestore.Desc).Limit(pageSize).Offset((page - 1) * pageSize)

		// Execute query
		iter := query.Documents(ctx)
		defer iter.Stop()

		var courses []models.Course
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch courses")
				return
			}

			var course models.Course
			doc.DataTo(&course)
			courses = append(courses, course)
		}

		utils.RespondSuccess(w, map[string]interface{}{
			"courses": courses,
			"page":    page,
			"pageSize": pageSize,
		})
	})(w, r)
}

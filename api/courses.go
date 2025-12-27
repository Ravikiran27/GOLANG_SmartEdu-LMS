package handler

import (
	"net/http"
	"strings"

	courseHandlers "github.com/Ravikiran27/GOLANG_SmartEdu-LMS/api/courses"
)

// Handler routes all course-related requests
func CoursesRouter(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /api/courses/
	path := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	path = strings.TrimPrefix(path, "courses/") // Handle both /api/courses and /api/courses/courses
	
	// Route to appropriate handler based on path
	switch path {
	case "create":
		courseHandlers.CreateCourse(w, r)
	case "list":
		courseHandlers.ListCourses(w, r)
	case "get":
		courseHandlers.GetCourse(w, r)
	case "update":
		courseHandlers.UpdateCourse(w, r)
	case "delete":
		courseHandlers.DeleteCourse(w, r)
	case "enroll":
		courseHandlers.EnrollStudent(w, r)
	case "progress":
		courseHandlers.UpdateProgress(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

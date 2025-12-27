package handler

import (
	"net/http"
	"strings"

	quizHandlers "github.com/Ravikiran27/GOLANG_SmartEdu-LMS/api/quizzes"
)

// Handler routes all quiz-related requests
func QuizzesRouter(w http.ResponseWriter, r *http.Request) {
	// Extract the path after /api/quizzes/
	path := strings.TrimPrefix(r.URL.Path, "/api/quizzes/")
	path = strings.TrimPrefix(path, "quizzes/") // Handle both /api/quizzes and /api/quizzes/quizzes
	
	// Route to appropriate handler based on path
	switch path {
	case "create":
		quizHandlers.CreateQuiz(w, r)
	case "list":
		quizHandlers.ListQuizzes(w, r)
	case "get":
		quizHandlers.GetQuiz(w, r)
	case "add-question":
		quizHandlers.AddQuestion(w, r)
	case "start":
		quizHandlers.StartQuiz(w, r)
	case "submit":
		quizHandlers.SubmitQuiz(w, r)
	case "results":
		quizHandlers.GetResults(w, r)
	case "resume":
		quizHandlers.ResumeQuiz(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	
	"github.com/yourusername/lms/models"
)

// ParseJSONBody decodes JSON request body into a struct
func ParseJSONBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// GetCurrentTimestamp returns the current time in Firestore format
func GetCurrentTimestamp() time.Time {
	return time.Now().UTC()
}

// CalculatePercentage calculates percentage (obtained/total * 100)
func CalculatePercentage(obtained, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (obtained / total) * 100
}

// Contains checks if a string exists in a slice
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateEmail checks if email format is valid (basic validation)
func ValidateEmail(email string) bool {
	// Basic check - Firebase will do detailed validation
	return len(email) > 3 && ContainsChar(email, '@') && ContainsChar(email, '.')
}

// ContainsChar checks if string contains a character
func ContainsChar(s string, c rune) bool {
	for _, ch := range s {
		if ch == c {
			return true
		}
	}
	return false
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// GetPaginationParams extracts pagination from query params
func GetPaginationParams(r *http.Request) (page, pageSize int) {
	page = 1
	pageSize = 20 // Default page size

	// Parse query parameters
	query := r.URL.Query()
	
	if p := query.Get("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	
	if ps := query.Get("pageSize"); ps != "" {
		if parsed, err := parseInt(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	
	return page, pageSize
}

// Helper to parse int from string
func parseInt(s string) (int, error) {
	var i int
	_, err := json.Unmarshal([]byte(s), &i)
	return i, err
}

// ShuffleQuestions randomizes question order for quiz
func ShuffleQuestions(questions *[]models.Question) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*questions), func(i, j int) {
		(*questions)[i], (*questions)[j] = (*questions)[j], (*questions)[i]
	})
}

// ShuffleOptions randomizes option order for a question
func ShuffleOptions(options *[]models.QuestionOption) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*options), func(i, j int) {
		(*options)[i], (*options)[j] = (*options)[j], (*options)[i]
	})
}

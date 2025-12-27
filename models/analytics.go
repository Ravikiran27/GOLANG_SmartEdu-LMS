package models

import "time"

// Analytics represents pre-aggregated analytics data
type Analytics struct {
	AnalyticsID string              `firestore:"analyticsId" json:"analyticsId"`
	Type        string              `firestore:"type" json:"type"` // student_performance | course_stats | quiz_stats | exam_stats
	EntityID    string              `firestore:"entityId" json:"entityId"`
	Period      string              `firestore:"period" json:"period"` // weekly | monthly | all_time
	Metrics     AnalyticsMetrics    `firestore:"metrics" json:"metrics"`
	ComputedAt  time.Time           `firestore:"computedAt" json:"computedAt"`
	StartDate   time.Time           `firestore:"startDate" json:"startDate"`
	EndDate     time.Time           `firestore:"endDate" json:"endDate"`
}

// AnalyticsMetrics contains various metrics
type AnalyticsMetrics struct {
	TotalAttempts    int     `firestore:"totalAttempts" json:"totalAttempts"`
	AverageScore     float64 `firestore:"averageScore" json:"averageScore"`
	HighestScore     float64 `firestore:"highestScore" json:"highestScore"`
	LowestScore      float64 `firestore:"lowestScore" json:"lowestScore"`
	CompletionRate   float64 `firestore:"completionRate" json:"completionRate"`
	EnrollmentCount  int     `firestore:"enrollmentCount" json:"enrollmentCount"`
	ActiveStudents   int     `firestore:"activeStudents" json:"activeStudents"`
}

// StudentPerformance represents student performance summary
type StudentPerformance struct {
	StudentID        string                  `json:"studentId"`
	StudentName      string                  `json:"studentName"`
	CoursesEnrolled  int                     `json:"coursesEnrolled"`
	CoursesCompleted int                     `json:"coursesCompleted"`
	QuizzesTaken     int                     `json:"quizzesTaken"`
	ExamsTaken       int                     `json:"examsTaken"`
	AverageScore     float64                 `json:"averageScore"`
	RecentActivities []ActivityLog           `json:"recentActivities"`
	CourseProgress   []CourseProgressSummary `json:"courseProgress"`
}

// CourseProgressSummary represents progress in a course
type CourseProgressSummary struct {
	CourseID    string  `json:"courseId"`
	CourseTitle string  `json:"courseTitle"`
	Progress    float64 `json:"progress"`
	Status      string  `json:"status"`
}

// ActivityLog represents a student activity
type ActivityLog struct {
	Type      string    `json:"type"` // enrollment | quiz | exam | assignment
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score,omitempty"`
}

// Notification represents a notification
type Notification struct {
	NotificationID string    `firestore:"notificationId" json:"notificationId"`
	UserID         string    `firestore:"userId" json:"userId"`
	Type           string    `firestore:"type" json:"type"` // course_update | quiz_published | assignment_due | grade_released
	Title          string    `firestore:"title" json:"title"`
	Message        string    `firestore:"message" json:"message"`
	ReferenceID    string    `firestore:"referenceId" json:"referenceId"`
	ReferenceType  string    `firestore:"referenceType" json:"referenceType"` // course | quiz | assignment | exam
	IsRead         bool      `firestore:"isRead" json:"isRead"`
	CreatedAt      time.Time `firestore:"createdAt" json:"createdAt"`
}

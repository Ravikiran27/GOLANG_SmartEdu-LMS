package models

import "time"

// Course represents a course in the system
type Course struct {
	CourseID        string           `firestore:"courseId" json:"courseId"`
	Title           string           `firestore:"title" json:"title"`
	Description     string           `firestore:"description" json:"description"`
	Syllabus        string           `firestore:"syllabus" json:"syllabus"`
	TeacherID       string           `firestore:"teacherId" json:"teacherId"`
	TeacherName     string           `firestore:"teacherName" json:"teacherName"`
	Category        string           `firestore:"category" json:"category"`
	Difficulty      string           `firestore:"difficulty" json:"difficulty"` // beginner | intermediate | advanced
	Thumbnail       string           `firestore:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	Materials       []CourseMaterial `firestore:"materials" json:"materials"`
	EnrollmentCount int              `firestore:"enrollmentCount" json:"enrollmentCount"`
	IsPublished     bool             `firestore:"isPublished" json:"isPublished"`
	CreatedAt       time.Time        `firestore:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `firestore:"updatedAt" json:"updatedAt"`
	IsDeleted       bool             `firestore:"isDeleted" json:"isDeleted"`
}

// CourseMaterial represents a course material
type CourseMaterial struct {
	ID         string    `firestore:"id" json:"id"`
	Name       string    `firestore:"name" json:"name"`
	Type       string    `firestore:"type" json:"type"` // pdf | ppt | video | doc
	URL        string    `firestore:"url" json:"url"`
	Size       int64     `firestore:"size" json:"size"`
	UploadedAt time.Time `firestore:"uploadedAt" json:"uploadedAt"`
}

// CreateCourseRequest represents course creation request
type CreateCourseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Syllabus    string `json:"syllabus"`
	Category    string `json:"category" validate:"required"`
	Difficulty  string `json:"difficulty" validate:"required,oneof=beginner intermediate advanced"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

// UpdateCourseRequest represents course update request
type UpdateCourseRequest struct {
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	Syllabus    string           `json:"syllabus,omitempty"`
	Category    string           `json:"category,omitempty"`
	Difficulty  string           `json:"difficulty,omitempty"`
	Thumbnail   string           `json:"thumbnail,omitempty"`
	Materials   []CourseMaterial `json:"materials,omitempty"`
	IsPublished bool             `json:"isPublished,omitempty"`
}

// Enrollment represents a course enrollment
type Enrollment struct {
	EnrollmentID       string    `firestore:"enrollmentId" json:"enrollmentId"`
	StudentID          string    `firestore:"studentId" json:"studentId"`
	StudentName        string    `firestore:"studentName" json:"studentName"`
	CourseID           string    `firestore:"courseId" json:"courseId"`
	CourseTitle        string    `firestore:"courseTitle" json:"courseTitle"`
	EnrolledAt         time.Time `firestore:"enrolledAt" json:"enrolledAt"`
	Progress           float64   `firestore:"progress" json:"progress"`
	CompletedMaterials []string  `firestore:"completedMaterials" json:"completedMaterials"`
	Status             string    `firestore:"status" json:"status"` // active | completed | dropped
	LastAccessedAt     time.Time `firestore:"lastAccessedAt" json:"lastAccessedAt"`
}

// EnrollmentRequest represents enrollment creation
type EnrollmentRequest struct {
	CourseID string `json:"courseId" validate:"required"`
}

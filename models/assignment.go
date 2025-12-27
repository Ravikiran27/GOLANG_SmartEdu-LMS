package models

import "time"

// Assignment represents an assignment
type Assignment struct {
	AssignmentID       string                 `firestore:"assignmentId" json:"assignmentId"`
	CourseID           string                 `firestore:"courseId" json:"courseId"`
	CourseTitle        string                 `firestore:"courseTitle" json:"courseTitle"`
	TeacherID          string                 `firestore:"teacherId" json:"teacherId"`
	Title              string                 `firestore:"title" json:"title"`
	Description        string                 `firestore:"description" json:"description"`
	Instructions       string                 `firestore:"instructions" json:"instructions"`
	Attachments        []AssignmentAttachment `firestore:"attachments" json:"attachments"`
	TotalMarks         float64                `firestore:"totalMarks" json:"totalMarks"`
	DueDate            time.Time              `firestore:"dueDate" json:"dueDate"`
	AllowLateSubmission bool                  `firestore:"allowLateSubmission" json:"allowLateSubmission"`
	LatePenalty        float64                `firestore:"latePenalty" json:"latePenalty"` // percentage per day
	IsPublished        bool                   `firestore:"isPublished" json:"isPublished"`
	CreatedAt          time.Time              `firestore:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time              `firestore:"updatedAt" json:"updatedAt"`
	IsDeleted          bool                   `firestore:"isDeleted" json:"isDeleted"`
}

// AssignmentAttachment represents an attachment in assignment
type AssignmentAttachment struct {
	Name string `firestore:"name" json:"name"`
	URL  string `firestore:"url" json:"url"`
	Type string `firestore:"type" json:"type"`
}

// CreateAssignmentRequest represents assignment creation
type CreateAssignmentRequest struct {
	CourseID            string                 `json:"courseId" validate:"required"`
	Title               string                 `json:"title" validate:"required"`
	Description         string                 `json:"description" validate:"required"`
	Instructions        string                 `json:"instructions"`
	Attachments         []AssignmentAttachment `json:"attachments"`
	TotalMarks          float64                `json:"totalMarks" validate:"required,min=1"`
	DueDate             time.Time              `json:"dueDate" validate:"required"`
	AllowLateSubmission bool                   `json:"allowLateSubmission"`
	LatePenalty         float64                `json:"latePenalty"`
}

// AssignmentSubmission represents student assignment submission
type AssignmentSubmission struct {
	SubmissionID     string                        `firestore:"submissionId" json:"submissionId"`
	AssignmentID     string                        `firestore:"assignmentId" json:"assignmentId"`
	StudentID        string                        `firestore:"studentId" json:"studentId"`
	StudentName      string                        `firestore:"studentName" json:"studentName"`
	SubmissionText   string                        `firestore:"submissionText,omitempty" json:"submissionText,omitempty"`
	Attachments      []AssignmentSubmissionFile    `firestore:"attachments" json:"attachments"`
	SubmittedAt      time.Time                     `firestore:"submittedAt" json:"submittedAt"`
	IsLateSubmission bool                          `firestore:"isLateSubmission" json:"isLateSubmission"`
	DaysLate         int                           `firestore:"daysLate" json:"daysLate"`
	MarksAwarded     float64                       `firestore:"marksAwarded" json:"marksAwarded"`
	Feedback         string                        `firestore:"feedback,omitempty" json:"feedback,omitempty"`
	Status           string                        `firestore:"status" json:"status"` // submitted | evaluated | returned
	EvaluatedAt      *time.Time                    `firestore:"evaluatedAt,omitempty" json:"evaluatedAt,omitempty"`
	EvaluatedBy      string                        `firestore:"evaluatedBy,omitempty" json:"evaluatedBy,omitempty"`
}

// AssignmentSubmissionFile represents a file in submission
type AssignmentSubmissionFile struct {
	Name       string    `firestore:"name" json:"name"`
	URL        string    `firestore:"url" json:"url"`
	Size       int64     `firestore:"size" json:"size"`
	UploadedAt time.Time `firestore:"uploadedAt" json:"uploadedAt"`
}

// SubmitAssignmentRequest represents assignment submission
type SubmitAssignmentRequest struct {
	AssignmentID   string                       `json:"assignmentId" validate:"required"`
	SubmissionText string                       `json:"submissionText"`
	Attachments    []AssignmentSubmissionFile   `json:"attachments"`
}

// EvaluateAssignmentRequest represents assignment evaluation by teacher
type EvaluateAssignmentRequest struct {
	SubmissionID string  `json:"submissionId" validate:"required"`
	MarksAwarded float64 `json:"marksAwarded" validate:"required,min=0"`
	Feedback     string  `json:"feedback"`
}

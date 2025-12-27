package models

import "time"

// Exam represents an exam
type Exam struct {
	ExamID                  string     `firestore:"examId" json:"examId"`
	CourseID                string     `firestore:"courseId" json:"courseId"`
	CourseTitle             string     `firestore:"courseTitle" json:"courseTitle"`
	TeacherID               string     `firestore:"teacherId" json:"teacherId"`
	Title                   string     `firestore:"title" json:"title"`
	Description             string     `firestore:"description" json:"description"`
	Duration                int        `firestore:"duration" json:"duration"` // minutes
	TotalMarks              float64    `firestore:"totalMarks" json:"totalMarks"`
	PassingMarks            float64    `firestore:"passingMarks" json:"passingMarks"`
	NegativeMarking         bool       `firestore:"negativeMarking" json:"negativeMarking"`
	NegativeMarkValue       float64    `firestore:"negativeMarkValue" json:"negativeMarkValue"`
	QuestionsCount          int        `firestore:"questionsCount" json:"questionsCount"`
	RandomizeQuestions      bool       `firestore:"randomizeQuestions" json:"randomizeQuestions"`
	ExamType                string     `firestore:"examType" json:"examType"` // midterm | final | practice
	StartTime               time.Time  `firestore:"startTime" json:"startTime"`
	EndTime                 time.Time  `firestore:"endTime" json:"endTime"`
	Instructions            string     `firestore:"instructions" json:"instructions"`
	IsPublished             bool       `firestore:"isPublished" json:"isPublished"`
	RequiresManualEvaluation bool      `firestore:"requiresManualEvaluation" json:"requiresManualEvaluation"`
	CreatedAt               time.Time  `firestore:"createdAt" json:"createdAt"`
	UpdatedAt               time.Time  `firestore:"updatedAt" json:"updatedAt"`
	IsDeleted               bool       `firestore:"isDeleted" json:"isDeleted"`
}

// CreateExamRequest represents exam creation
type CreateExamRequest struct {
	CourseID           string    `json:"courseId" validate:"required"`
	Title              string    `json:"title" validate:"required"`
	Description        string    `json:"description"`
	Duration           int       `json:"duration" validate:"required,min=1"`
	TotalMarks         float64   `json:"totalMarks" validate:"required,min=1"`
	PassingMarks       float64   `json:"passingMarks" validate:"required"`
	NegativeMarking    bool      `json:"negativeMarking"`
	NegativeMarkValue  float64   `json:"negativeMarkValue"`
	RandomizeQuestions bool      `json:"randomizeQuestions"`
	ExamType           string    `json:"examType" validate:"required,oneof=midterm final practice"`
	StartTime          time.Time `json:"startTime" validate:"required"`
	EndTime            time.Time `json:"endTime" validate:"required"`
	Instructions       string    `json:"instructions"`
}

// ExamSubmission represents an exam submission
type ExamSubmission struct {
	SubmissionID   string                `firestore:"submissionId" json:"submissionId"`
	ExamID         string                `firestore:"examId" json:"examId"`
	StudentID      string                `firestore:"studentId" json:"studentId"`
	StudentName    string                `firestore:"studentName" json:"studentName"`
	Answers        []ExamSubmittedAnswer `firestore:"answers" json:"answers"`
	StartedAt      time.Time             `firestore:"startedAt" json:"startedAt"`
	SubmittedAt    time.Time             `firestore:"submittedAt" json:"submittedAt"`
	TimeTaken      int                   `firestore:"timeTaken" json:"timeTaken"` // minutes
	TotalMarks     float64               `firestore:"totalMarks" json:"totalMarks"`
	MarksObtained  float64               `firestore:"marksObtained" json:"marksObtained"`
	Percentage     float64               `firestore:"percentage" json:"percentage"`
	Status         string                `firestore:"status" json:"status"` // in_progress | submitted | partially_evaluated | evaluated
	EvaluatedAt    *time.Time            `firestore:"evaluatedAt,omitempty" json:"evaluatedAt,omitempty"`
	EvaluatedBy    string                `firestore:"evaluatedBy,omitempty" json:"evaluatedBy,omitempty"`
	TeacherComments string               `firestore:"teacherComments,omitempty" json:"teacherComments,omitempty"`
}

// ExamSubmittedAnswer represents a single answer in exam submission
type ExamSubmittedAnswer struct {
	QuestionID      string  `firestore:"questionId" json:"questionId"`
	QuestionType    string  `firestore:"questionType" json:"questionType"`
	SelectedAnswer  string  `firestore:"selectedAnswer" json:"selectedAnswer"`
	IsCorrect       *bool   `firestore:"isCorrect,omitempty" json:"isCorrect,omitempty"`
	MarksAwarded    float64 `firestore:"marksAwarded" json:"marksAwarded"`
	TeacherFeedback string  `firestore:"teacherFeedback,omitempty" json:"teacherFeedback,omitempty"`
}

// SubmitExamRequest represents exam submission
type SubmitExamRequest struct {
	ExamID  string                    `json:"examId" validate:"required"`
	Answers []SubmitExamAnswerRequest `json:"answers" validate:"required"`
}

// SubmitExamAnswerRequest represents individual exam answer
type SubmitExamAnswerRequest struct {
	QuestionID     string `json:"questionId" validate:"required"`
	SelectedAnswer string `json:"selectedAnswer" validate:"required"`
}

// EvaluateExamRequest represents manual evaluation by teacher
type EvaluateExamRequest struct {
	SubmissionID    string                  `json:"submissionId" validate:"required"`
	Evaluations     []AnswerEvaluation      `json:"evaluations" validate:"required"`
	TeacherComments string                  `json:"teacherComments"`
}

// AnswerEvaluation represents evaluation of a single answer
type AnswerEvaluation struct {
	QuestionID      string  `json:"questionId" validate:"required"`
	MarksAwarded    float64 `json:"marksAwarded" validate:"required,min=0"`
	TeacherFeedback string  `json:"teacherFeedback"`
}

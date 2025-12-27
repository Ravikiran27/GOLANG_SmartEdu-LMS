package models

import "time"

// Quiz represents a quiz
type Quiz struct {
	ID                 string    `firestore:"id" json:"id"`
	QuizID             string    `firestore:"quizId" json:"quizId"`
	CourseID           string    `firestore:"courseId" json:"courseId"`
	CourseTitle        string    `firestore:"courseTitle" json:"courseTitle"`
	TeacherID          string    `firestore:"teacherId" json:"teacherId"`
	Title              string    `firestore:"title" json:"title"`
	Description        string    `firestore:"description" json:"description"`
	Duration           int       `firestore:"duration" json:"duration"` // minutes
	TotalMarks         float64   `firestore:"totalMarks" json:"totalMarks"`
	PassingMarks       float64   `firestore:"passingMarks" json:"passingMarks"`
	NegativeMarking    bool      `firestore:"negativeMarking" json:"negativeMarking"`
	NegativeMarkValue  float64   `firestore:"negativeMarkValue" json:"negativeMarkValue"`
	QuestionsCount     int       `firestore:"questionsCount" json:"questionsCount"`
	QuestionCount      int       `firestore:"questionCount" json:"questionCount"`
	RandomizeQuestions bool      `firestore:"randomizeQuestions" json:"randomizeQuestions"`
	ShowResults        bool      `firestore:"showResults" json:"showResults"`
	ShowResultsAfterSubmit bool  `firestore:"showResultsAfterSubmit" json:"showResultsAfterSubmit"`
	ShuffleQuestions   bool      `firestore:"shuffleQuestions" json:"shuffleQuestions"`
	ShuffleOptions     bool      `firestore:"shuffleOptions" json:"shuffleOptions"`
	AllowReview        bool      `firestore:"allowReview" json:"allowReview"`
	AllowedAttempts    int       `firestore:"allowedAttempts" json:"allowedAttempts"` // 0 = unlimited
	MaxAttempts        int       `firestore:"maxAttempts" json:"maxAttempts"`
	Instructions       string    `firestore:"instructions" json:"instructions"`
	Deadline           time.Time `firestore:"deadline" json:"deadline"`
	StartDate          *time.Time `firestore:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate            *time.Time `firestore:"endDate,omitempty" json:"endDate,omitempty"`
	
	// Cheating prevention features
	PreventTabSwitch       bool `firestore:"preventTabSwitch" json:"preventTabSwitch"`
	MaxTabSwitches         int  `firestore:"maxTabSwitches" json:"maxTabSwitches"`
	RequireFullscreen      bool `firestore:"requireFullscreen" json:"requireFullscreen"`
	DisableCopyPaste       bool `firestore:"disableCopyPaste" json:"disableCopyPaste"`
	EnableProctoring       bool `firestore:"enableProctoring" json:"enableProctoring"`
	RandomizeQuestionOrder bool `firestore:"randomizeQuestionOrder" json:"randomizeQuestionOrder"`
	TimePerQuestion        int  `firestore:"timePerQuestion" json:"timePerQuestion"` // seconds
	LockAfterSubmit        bool `firestore:"lockAfterSubmit" json:"lockAfterSubmit"`
	
	// Teacher permissions
	AllowTeacherResume     bool `firestore:"allowTeacherResume" json:"allowTeacherResume"`
	AllowTeacherExtendTime bool `firestore:"allowTeacherExtendTime" json:"allowTeacherExtendTime"`
	
	IsPublished        bool      `firestore:"isPublished" json:"isPublished"`
	CreatedAt          time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time `firestore:"updatedAt" json:"updatedAt"`
	IsDeleted          bool      `firestore:"isDeleted" json:"isDeleted"`
}

// Question represents a quiz/exam question
type Question struct {
	ID           string          `firestore:"id" json:"id"`
	QuestionID   string          `firestore:"questionId" json:"questionId"`
	QuizID       string          `firestore:"quizId,omitempty" json:"quizId,omitempty"`
	ExamID       string          `firestore:"examId,omitempty" json:"examId,omitempty"`
	Type         string          `firestore:"type" json:"type"` // mcq | true_false | short_answer | descriptive | long_answer
	Text         string          `firestore:"text" json:"text"`
	QuestionText string          `firestore:"questionText" json:"questionText"`
	ImageURL     string          `firestore:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Marks        float64         `firestore:"marks" json:"marks"`
	Points       float64         `firestore:"points" json:"points"`
	Options      []QuestionOption `firestore:"options,omitempty" json:"options,omitempty"`
	CorrectAnswer string         `firestore:"correctAnswer,omitempty" json:"correctAnswer,omitempty"`
	Explanation  string          `firestore:"explanation,omitempty" json:"explanation,omitempty"`
	Order        int             `firestore:"order" json:"order"`
	CreatedAt    time.Time       `firestore:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time       `firestore:"updatedAt" json:"updatedAt"`
}

// QuestionOption represents an option for MCQ/True-False
type QuestionOption struct {
	ID        string `firestore:"id" json:"id"` // A, B, C, D or True, False
	Text      string `firestore:"text" json:"text"`
	IsCorrect bool   `firestore:"isCorrect" json:"isCorrect"`
}

// CreateQuizRequest represents quiz creation
type CreateQuizRequest struct {
	CourseID           string     `json:"courseId" validate:"required"`
	Title              string     `json:"title" validate:"required"`
	Description        string     `json:"description"`
	Duration           int        `json:"duration" validate:"required,min=1"`
	TotalMarks         float64    `json:"totalMarks" validate:"required,min=1"`
	PassingMarks       float64    `json:"passingMarks" validate:"required"`
	NegativeMarking    bool       `json:"negativeMarking"`
	NegativeMarkValue  float64    `json:"negativeMarkValue"`
	Instructions       string     `json:"instructions"`
	Deadline           time.Time  `json:"deadline"`
	RandomizeQuestions bool       `json:"randomizeQuestions"`
	ShowResults        bool       `json:"showResults"`
	ShowResultsAfterSubmit bool   `json:"showResultsAfterSubmit"`
	ShuffleQuestions   bool       `json:"shuffleQuestions"`
	ShuffleOptions     bool       `json:"shuffleOptions"`
	AllowReview        bool       `json:"allowReview"`
	AllowedAttempts    int        `json:"allowedAttempts"`
	MaxAttempts        int        `json:"maxAttempts"`
	StartDate          *time.Time `json:"startDate,omitempty"`
	EndDate            *time.Time `json:"endDate,omitempty"`
	
	// Cheating prevention
	PreventTabSwitch       bool `json:"preventTabSwitch"`
	MaxTabSwitches         int  `json:"maxTabSwitches"`
	RequireFullscreen      bool `json:"requireFullscreen"`
	DisableCopyPaste       bool `json:"disableCopyPaste"`
	EnableProctoring       bool `json:"enableProctoring"`
	RandomizeQuestionOrder bool `json:"randomizeQuestionOrder"`
	TimePerQuestion        int  `json:"timePerQuestion"`
	
	// Teacher permissions
	AllowTeacherResume     bool `json:"allowTeacherResume"`
	AllowTeacherExtendTime bool `json:"allowTeacherExtendTime"`
	
	IsPublished        bool `json:"isPublished"`
}

// CreateQuestionRequest represents question creation
type CreateQuestionRequest struct {
	QuizID       string           `json:"quizId,omitempty"`
	ExamID       string           `json:"examId,omitempty"`
	Type         string           `json:"type" validate:"required,oneof=mcq true_false short_answer descriptive"`
	QuestionText string           `json:"questionText" validate:"required"`
	ImageURL     string           `json:"imageUrl,omitempty"`
	Marks        float64          `json:"marks" validate:"required,min=0.5"`
	Options      []QuestionOption `json:"options,omitempty"`
	CorrectAnswer string          `json:"correctAnswer,omitempty"`
	Explanation  string           `json:"explanation,omitempty"`
}

// QuizSubmission represents a quiz submission
type QuizSubmission struct {
	ID            string           `firestore:"id" json:"id"`
	SubmissionID  string           `firestore:"submissionId" json:"submissionId"`
	QuizID        string           `firestore:"quizId" json:"quizId"`
	StudentID     string           `firestore:"studentId" json:"studentId"`
	StudentName   string           `firestore:"studentName" json:"studentName"`
	StudentEmail  string           `firestore:"studentEmail,omitempty" json:"studentEmail,omitempty"`
	CourseID      string           `firestore:"courseId" json:"courseId"`
	AttemptNumber int              `firestore:"attemptNumber" json:"attemptNumber"`
	Answers       []Answer         `firestore:"answers" json:"answers"`
	Questions     []Question       `firestore:"questions" json:"questions"`
	StartedAt     time.Time        `firestore:"startedAt" json:"startedAt"`
	SubmittedAt   time.Time        `firestore:"submittedAt" json:"submittedAt"`
	TimeTaken     int              `firestore:"timeTaken" json:"timeTaken"` // minutes
	TimeLimit     int              `firestore:"timeLimit" json:"timeLimit"` // minutes
	TotalMarks    float64          `firestore:"totalMarks" json:"totalMarks"`
	MarksObtained float64          `firestore:"marksObtained" json:"marksObtained"`
	Score         float64          `firestore:"score" json:"score"`
	Percentage    float64          `firestore:"percentage" json:"percentage"`
	Passed        bool             `firestore:"passed" json:"passed"`
	Status        string           `firestore:"status" json:"status"` // in_progress | submitted | evaluated
	EvaluatedAt   *time.Time       `firestore:"evaluatedAt,omitempty" json:"evaluatedAt,omitempty"`
	EvaluatedBy   string           `firestore:"evaluatedBy,omitempty" json:"evaluatedBy,omitempty"`
	
	// Cheating detection
	TabSwitchCount     int      `firestore:"tabSwitchCount" json:"tabSwitchCount"`
	FullscreenExits    int      `firestore:"fullscreenExits" json:"fullscreenExits"`
	SuspiciousActivity []string `firestore:"suspiciousActivity" json:"suspiciousActivity"`
	
	// Teacher resume tracking
	ResumedBy    string     `firestore:"resumedBy,omitempty" json:"resumedBy,omitempty"`
	ResumedAt    time.Time  `firestore:"resumedAt,omitempty" json:"resumedAt,omitempty"`
	ResumeReason string     `firestore:"resumeReason,omitempty" json:"resumeReason,omitempty"`
	
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `firestore:"updatedAt" json:"updatedAt"`
}

// Answer represents a single answer in a submission
type Answer struct {
	QuestionID      string   `firestore:"questionId" json:"questionId"`
	SelectedOptions []string `firestore:"selectedOptions" json:"selectedOptions"`
	TextAnswer      string   `firestore:"textAnswer" json:"textAnswer"`
	IsCorrect       bool     `firestore:"isCorrect" json:"isCorrect"`
	PointsAwarded   float64  `firestore:"pointsAwarded" json:"pointsAwarded"`
}

// SubmittedAnswer represents a single answer in a submission
type SubmittedAnswer struct {
	QuestionID     string  `firestore:"questionId" json:"questionId"`
	SelectedAnswer string  `firestore:"selectedAnswer" json:"selectedAnswer"`
	IsCorrect      bool    `firestore:"isCorrect" json:"isCorrect"`
	MarksAwarded   float64 `firestore:"marksAwarded" json:"marksAwarded"`
}

// SubmitQuizRequest represents quiz submission from student
type SubmitQuizRequest struct {
	SubmissionID    string   `json:"submissionId"`
	QuizID          string   `json:"quizId" validate:"required"`
	Answers         []Answer `json:"answers" validate:"required"`
	TabSwitches     int      `json:"tabSwitches"`
	FullscreenExits int      `json:"fullscreenExits"`
	TimedOut        bool     `json:"timedOut"`
}

// SubmitQuizAnswerRequest represents individual answer
type SubmitQuizAnswerRequest struct {
	QuestionID     string `json:"questionId" validate:"required"`
	SelectedAnswer string `json:"selectedAnswer" validate:"required"`
}

# Firestore Schema Design - LMS Platform

## Design Principles
- Flat collections (minimal nesting for scalability)
- Indexed fields for common queries
- Denormalization where needed for performance
- Timestamp tracking for all documents
- Soft delete support

---

## Collections

### 1. users
**Path:** `/users/{userId}`

```json
{
  "uid": "string (Firebase Auth UID)",
  "email": "string",
  "displayName": "string",
  "role": "string (admin | teacher | student)",
  "photoURL": "string (optional)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "isActive": "boolean",
  "metadata": {
    "lastLogin": "timestamp",
    "department": "string (optional)",
    "rollNumber": "string (for students)",
    "employeeId": "string (for teachers)"
  }
}
```

**Indexes:**
- role (ascending)
- email (ascending)
- isActive (ascending)

---

### 2. courses
**Path:** `/courses/{courseId}`

```json
{
  "courseId": "string (auto-generated)",
  "title": "string",
  "description": "string",
  "syllabus": "string (long text)",
  "teacherId": "string (ref to users)",
  "teacherName": "string (denormalized)",
  "category": "string",
  "difficulty": "string (beginner | intermediate | advanced)",
  "thumbnail": "string (Storage URL)",
  "materials": [
    {
      "id": "string",
      "name": "string",
      "type": "string (pdf | ppt | video | doc)",
      "url": "string (Storage URL)",
      "size": "number (bytes)",
      "uploadedAt": "timestamp"
    }
  ],
  "enrollmentCount": "number (denormalized for performance)",
  "isPublished": "boolean",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "isDeleted": "boolean"
}
```

**Indexes:**
- teacherId (ascending)
- isPublished (ascending)
- category (ascending)
- createdAt (descending)

---

### 3. enrollments
**Path:** `/enrollments/{enrollmentId}`

```json
{
  "enrollmentId": "string (auto-generated)",
  "studentId": "string (ref to users)",
  "studentName": "string (denormalized)",
  "courseId": "string (ref to courses)",
  "courseTitle": "string (denormalized)",
  "enrolledAt": "timestamp",
  "progress": "number (0-100)",
  "completedMaterials": ["string (material IDs)"],
  "status": "string (active | completed | dropped)",
  "lastAccessedAt": "timestamp"
}
```

**Indexes:**
- studentId + courseId (composite, unique)
- courseId (ascending)
- studentId (ascending)
- status (ascending)

---

### 4. quizzes
**Path:** `/quizzes/{quizId}`

```json
{
  "quizId": "string (auto-generated)",
  "courseId": "string (ref to courses)",
  "courseTitle": "string (denormalized)",
  "teacherId": "string (ref to users)",
  "title": "string",
  "description": "string",
  "duration": "number (minutes)",
  "totalMarks": "number",
  "passingMarks": "number",
  "negativeMarking": "boolean",
  "negativeMarkValue": "number (e.g., 0.25)",
  "questionsCount": "number (denormalized)",
  "randomizeQuestions": "boolean",
  "showResults": "boolean (immediately after submission)",
  "allowedAttempts": "number (0 = unlimited)",
  "startDate": "timestamp (nullable)",
  "endDate": "timestamp (nullable)",
  "isPublished": "boolean",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "isDeleted": "boolean"
}
```

**Indexes:**
- courseId (ascending)
- teacherId (ascending)
- isPublished (ascending)

---

### 5. questions
**Path:** `/questions/{questionId}`

```json
{
  "questionId": "string (auto-generated)",
  "quizId": "string (ref to quizzes)",
  "examId": "string (ref to exams, nullable)",
  "type": "string (mcq | true_false | short_answer | descriptive)",
  "questionText": "string (supports markdown)",
  "imageUrl": "string (optional)",
  "marks": "number",
  "options": [
    {
      "id": "string (A, B, C, D)",
      "text": "string",
      "isCorrect": "boolean (for MCQ/True-False)"
    }
  ],
  "correctAnswer": "string (for short_answer/descriptive - teacher reference)",
  "explanation": "string (optional, shown after submission)",
  "order": "number",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Indexes:**
- quizId (ascending)
- examId (ascending)
- type (ascending)

---

### 6. exams
**Path:** `/exams/{examId}`

```json
{
  "examId": "string (auto-generated)",
  "courseId": "string (ref to courses)",
  "courseTitle": "string (denormalized)",
  "teacherId": "string (ref to users)",
  "title": "string",
  "description": "string",
  "duration": "number (minutes)",
  "totalMarks": "number",
  "passingMarks": "number",
  "negativeMarking": "boolean",
  "negativeMarkValue": "number",
  "questionsCount": "number (denormalized)",
  "randomizeQuestions": "boolean",
  "examType": "string (midterm | final | practice | assignment_exam)",
  "startTime": "timestamp (scheduled start)",
  "endTime": "timestamp (scheduled end)",
  "instructions": "string",
  "isPublished": "boolean",
  "requiresManualEvaluation": "boolean",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "isDeleted": "boolean"
}
```

**Indexes:**
- courseId (ascending)
- teacherId (ascending)
- examType (ascending)
- startTime (ascending)

---

### 7. quiz_submissions
**Path:** `/quiz_submissions/{submissionId}`

```json
{
  "submissionId": "string (auto-generated)",
  "quizId": "string (ref to quizzes)",
  "studentId": "string (ref to users)",
  "studentName": "string (denormalized)",
  "attemptNumber": "number",
  "answers": [
    {
      "questionId": "string",
      "selectedAnswer": "string (option ID or text)",
      "isCorrect": "boolean (for auto-eval)",
      "marksAwarded": "number"
    }
  ],
  "startedAt": "timestamp",
  "submittedAt": "timestamp",
  "timeTaken": "number (minutes)",
  "totalMarks": "number",
  "marksObtained": "number",
  "percentage": "number",
  "status": "string (in_progress | submitted | evaluated)",
  "evaluatedAt": "timestamp (nullable)",
  "evaluatedBy": "string (teacherId, nullable)"
}
```

**Indexes:**
- quizId + studentId (composite)
- studentId (ascending)
- quizId (ascending)
- status (ascending)

---

### 8. exam_submissions
**Path:** `/exam_submissions/{submissionId}`

```json
{
  "submissionId": "string (auto-generated)",
  "examId": "string (ref to exams)",
  "studentId": "string (ref to users)",
  "studentName": "string (denormalized)",
  "answers": [
    {
      "questionId": "string",
      "questionType": "string",
      "selectedAnswer": "string (option ID or text)",
      "isCorrect": "boolean (for auto-eval, null for manual)",
      "marksAwarded": "number (nullable until evaluated)",
      "teacherFeedback": "string (for manual evaluation)"
    }
  ],
  "startedAt": "timestamp",
  "submittedAt": "timestamp",
  "timeTaken": "number (minutes)",
  "totalMarks": "number",
  "marksObtained": "number (nullable until fully evaluated)",
  "percentage": "number (nullable)",
  "status": "string (in_progress | submitted | partially_evaluated | evaluated)",
  "evaluatedAt": "timestamp (nullable)",
  "evaluatedBy": "string (teacherId, nullable)",
  "teacherComments": "string (overall feedback)"
}
```

**Indexes:**
- examId + studentId (composite, unique)
- studentId (ascending)
- examId (ascending)
- status (ascending)

---

### 9. assignments
**Path:** `/assignments/{assignmentId}`

```json
{
  "assignmentId": "string (auto-generated)",
  "courseId": "string (ref to courses)",
  "courseTitle": "string (denormalized)",
  "teacherId": "string (ref to users)",
  "title": "string",
  "description": "string",
  "instructions": "string",
  "attachments": [
    {
      "name": "string",
      "url": "string (Storage URL)",
      "type": "string"
    }
  ],
  "totalMarks": "number",
  "dueDate": "timestamp",
  "allowLateSubmission": "boolean",
  "latePenalty": "number (percentage deduction per day)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp",
  "isPublished": "boolean",
  "isDeleted": "boolean"
}
```

**Indexes:**
- courseId (ascending)
- teacherId (ascending)
- dueDate (ascending)

---

### 10. assignment_submissions
**Path:** `/assignment_submissions/{submissionId}`

```json
{
  "submissionId": "string (auto-generated)",
  "assignmentId": "string (ref to assignments)",
  "studentId": "string (ref to users)",
  "studentName": "string (denormalized)",
  "submissionText": "string (optional)",
  "attachments": [
    {
      "name": "string",
      "url": "string (Storage URL)",
      "size": "number",
      "uploadedAt": "timestamp"
    }
  ],
  "submittedAt": "timestamp",
  "isLateSubmission": "boolean",
  "daysLate": "number",
  "marksAwarded": "number (nullable)",
  "feedback": "string",
  "status": "string (submitted | evaluated | returned)",
  "evaluatedAt": "timestamp (nullable)",
  "evaluatedBy": "string (teacherId, nullable)"
}
```

**Indexes:**
- assignmentId + studentId (composite, unique)
- assignmentId (ascending)
- studentId (ascending)
- status (ascending)

---

### 11. analytics
**Path:** `/analytics/{analyticsId}`

Pre-aggregated analytics for performance dashboards.

```json
{
  "analyticsId": "string (pattern: {type}_{entityId}_{period})",
  "type": "string (student_performance | course_stats | quiz_stats | exam_stats)",
  "entityId": "string (studentId | courseId | quizId | examId)",
  "period": "string (weekly | monthly | all_time)",
  "metrics": {
    "totalAttempts": "number",
    "averageScore": "number",
    "highestScore": "number",
    "lowestScore": "number",
    "completionRate": "number",
    "enrollmentCount": "number",
    "activeStudents": "number"
  },
  "computedAt": "timestamp",
  "startDate": "timestamp",
  "endDate": "timestamp"
}
```

**Indexes:**
- type + entityId (composite)
- period (ascending)

---

### 12. notifications
**Path:** `/notifications/{notificationId}`

```json
{
  "notificationId": "string (auto-generated)",
  "userId": "string (recipient)",
  "type": "string (course_update | quiz_published | assignment_due | grade_released)",
  "title": "string",
  "message": "string",
  "referenceId": "string (courseId | quizId | assignmentId)",
  "referenceType": "string (course | quiz | assignment | exam)",
  "isRead": "boolean",
  "createdAt": "timestamp"
}
```

**Indexes:**
- userId + isRead (composite)
- userId + createdAt (composite)

---

## Security Rules Strategy

```javascript
// User can only read/write their own user document
match /users/{userId} {
  allow read: if request.auth.uid == userId;
  allow write: if request.auth.uid == userId && 
               request.resource.data.role == resource.data.role; // Can't change own role
}

// Admin can manage all users
match /users/{userId} {
  allow read, write: if request.auth.token.role == 'admin';
}

// Course read access: public if published, teacher can read own
match /courses/{courseId} {
  allow read: if resource.data.isPublished == true || 
              request.auth.token.role in ['admin', 'teacher'];
  allow write: if request.auth.token.role in ['admin', 'teacher'];
}

// Enrollment: student can create own, read own
match /enrollments/{enrollmentId} {
  allow read: if request.auth.uid == resource.data.studentId ||
              request.auth.token.role in ['admin', 'teacher'];
  allow create: if request.auth.uid == request.resource.data.studentId &&
                request.auth.token.role == 'student';
}

// Submissions: student can create/read own, teacher can read all for their courses
match /quiz_submissions/{submissionId} {
  allow read: if request.auth.uid == resource.data.studentId ||
              request.auth.token.role in ['admin', 'teacher'];
  allow create: if request.auth.uid == request.resource.data.studentId;
}
```

---

## Data Access Patterns

### Common Queries
1. **Get student enrollments:** `enrollments where studentId == {uid}`
2. **Get course students:** `enrollments where courseId == {courseId}`
3. **Get teacher courses:** `courses where teacherId == {uid}`
4. **Get published quizzes for course:** `quizzes where courseId == {id} AND isPublished == true`
5. **Get student quiz attempts:** `quiz_submissions where quizId == {id} AND studentId == {uid}`
6. **Get pending evaluations:** `exam_submissions where status == 'submitted' AND examId in teacherExams`

---

## Scalability Considerations
1. **Denormalization:** Teacher names, course titles stored in enrollments for fast reads
2. **Counters:** enrollmentCount, questionsCount pre-calculated
3. **Analytics Collection:** Pre-aggregated data to avoid real-time calculations
4. **Batch Writes:** Update enrollment counts via Cloud Functions (if needed later)
5. **Pagination:** All list queries use cursor-based pagination with createdAt
6. **Composite Indexes:** Created for common multi-field queries

---

## Migration & Versioning
- Include `schemaVersion: 1` in all new documents
- Support backward compatibility when schema changes
- Use Cloud Functions for data migrations (future)

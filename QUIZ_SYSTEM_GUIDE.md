# Quiz System with Cheating Prevention - Complete Guide

## Overview

The LMS Quiz System includes **comprehensive cheating prevention** features and **teacher override permissions** to ensure fair assessment while maintaining flexibility for exceptional circumstances.

---

## 🛡️ Cheating Prevention Features

### 1. Tab Switching Detection

**Purpose**: Prevent students from opening other tabs during the quiz.

**Implementation**:
```json
{
  "preventTabSwitch": true,
  "maxTabSwitches": 3
}
```

**How it works**:
- Frontend tracks `visibilitychange` and `blur` events
- Each tab switch increments the counter
- If `tabSwitchCount > maxTabSwitches`, quiz is flagged as suspicious
- Counter is sent with submission

**Frontend Implementation** (React):
```typescript
useEffect(() => {
  let tabSwitches = 0;
  
  const handleVisibilityChange = () => {
    if (document.hidden && quizSettings.preventTabSwitch) {
      tabSwitches++;
      
      if (tabSwitches >= quizSettings.maxTabSwitches) {
        alert('Maximum tab switches exceeded! Your submission will be flagged.');
      }
    }
  };
  
  document.addEventListener('visibilitychange', handleVisibilityChange);
  
  return () => {
    document.removeEventListener('visibilitychange', handleVisibilityChange);
  };
}, []);
```

---

### 2. Fullscreen Mode Requirement

**Purpose**: Force students to stay in fullscreen mode.

**Implementation**:
```json
{
  "requireFullscreen": true
}
```

**How it works**:
- Quiz automatically enters fullscreen on start
- Tracks `fullscreenchange` events
- Each exit from fullscreen increments counter
- Flagged as suspicious activity

**Frontend Implementation**:
```typescript
const enterFullscreen = () => {
  const elem = document.documentElement;
  if (elem.requestFullscreen) {
    elem.requestFullscreen();
  }
};

useEffect(() => {
  let fullscreenExits = 0;
  
  const handleFullscreenChange = () => {
    if (!document.fullscreenElement && quizSettings.requireFullscreen) {
      fullscreenExits++;
      alert('Please stay in fullscreen mode!');
      enterFullscreen();
    }
  };
  
  document.addEventListener('fullscreenchange', handleFullscreenChange);
  
  return () => {
    document.removeEventListener('fullscreenchange', handleFullscreenChange);
  };
}, []);
```

---

### 3. Disable Copy/Paste

**Purpose**: Prevent copying questions or pasting answers.

**Implementation**:
```json
{
  "disableCopyPaste": true
}
```

**Frontend Implementation**:
```typescript
useEffect(() => {
  const preventCopyPaste = (e: ClipboardEvent) => {
    if (quizSettings.disableCopyPaste) {
      e.preventDefault();
      alert('Copy/paste is disabled during the quiz');
    }
  };
  
  document.addEventListener('copy', preventCopyPaste);
  document.addEventListener('paste', preventCopyPaste);
  document.addEventListener('cut', preventCopyPaste);
  
  return () => {
    document.removeEventListener('copy', preventCopyPaste);
    document.removeEventListener('paste', preventCopyPaste);
    document.removeEventListener('cut', preventCopyPaste);
  };
}, []);
```

---

### 4. Question Randomization

**Purpose**: Ensure each student gets questions in different order.

**Implementation**:
```json
{
  "shuffleQuestions": true,
  "shuffleOptions": true,
  "randomizeQuestionOrder": true
}
```

**How it works**:
- `shuffleQuestions`: Randomizes question order per student
- `shuffleOptions`: Randomizes MCQ option order
- Backend shuffles on quiz start (in `start.go`)
- Original order preserved in `questions` collection

---

### 5. Time Limits

**Purpose**: Prevent students from taking excessive time.

**Implementation**:
```json
{
  "duration": 60,
  "timePerQuestion": 120
}
```

**How it works**:
- `duration`: Total quiz time in minutes
- `timePerQuestion`: Max seconds per question (optional)
- Frontend countdown timer
- Auto-submit when time expires
- Backend validates submission time

**Frontend Timer**:
```typescript
const [timeLeft, setTimeLeft] = useState(quizDuration * 60);

useEffect(() => {
  const timer = setInterval(() => {
    setTimeLeft(prev => {
      if (prev <= 1) {
        submitQuiz(true); // Auto-submit with timedOut flag
        return 0;
      }
      return prev - 1;
    });
  }, 1000);
  
  return () => clearInterval(timer);
}, []);
```

---

### 6. Lock After Submit

**Purpose**: Prevent students from retaking quiz after submission.

**Implementation**:
```json
{
  "lockAfterSubmit": true,
  "maxAttempts": 1
}
```

**How it works**:
- After submission, status changes to `"submitted"` or `"evaluated"`
- Student cannot start new attempt if `maxAttempts` reached
- Teacher can override with `resume` permission

---

### 7. Proctoring Support (Future Enhancement)

**Purpose**: Enable webcam/screen recording.

**Implementation**:
```json
{
  "enableProctoring": true
}
```

**Planned Features**:
- Webcam snapshot every 30 seconds
- Screen recording (with permission)
- AI-based suspicious behavior detection
- Facial recognition verification

---

## 👨‍🏫 Teacher Override Permissions

### 1. Resume Locked Quiz

**Purpose**: Allow teacher to unlock a student's quiz after submission (e.g., technical issues, unfair flagging).

**API Endpoint**: `POST /api/quizzes/resume`

**Request**:
```json
{
  "submissionId": "sub123",
  "extendTime": 15,
  "reason": "Student experienced internet disconnection"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Quiz resumed successfully",
  "data": {
    "submissionId": "sub123",
    "extendedTime": 15,
    "newTimeLimit": 75
  }
}
```

**Usage Flow**:
1. Student submits quiz (status: `"submitted"`)
2. Student reports technical issue to teacher
3. Teacher calls resume API
4. Submission status changes back to `"in_progress"`
5. Student can continue from where they left off
6. Resume action logged in `audit_logs` collection

**Security**:
- Only teacher who owns the quiz (or admin) can resume
- Quiz must have `allowTeacherResume: true`
- Creates notification for student
- Logs all resume actions with reason

---

### 2. Extend Time During Quiz

**Purpose**: Give extra time to student during active quiz.

**API Endpoint**: `POST /api/quizzes/resume`

**Request**:
```json
{
  "submissionId": "sub123",
  "extendTime": 20,
  "reason": "Accommodations for student with disability"
}
```

**How it works**:
- Updates `timeLimit` field in submission
- Student's timer increases by extended minutes
- Requires `allowTeacherExtendTime: true` in quiz

---

## 📊 Suspicious Activity Detection

### Backend Flagging

When a quiz is submitted, the backend checks:

```go
// Check for suspicious activity
suspiciousFlags := make([]string, 0)

if quiz.PreventTabSwitch && req.TabSwitches > quiz.MaxTabSwitches {
    suspiciousFlags = append(suspiciousFlags, "Excessive tab switching detected")
}

if quiz.RequireFullscreen && req.FullscreenExits > 0 {
    suspiciousFlags = append(suspiciousFlags, "Exited fullscreen mode")
}

if req.TimedOut {
    suspiciousFlags = append(suspiciousFlags, "Time limit exceeded")
}

// Save to submission
updates = append(updates, firestore.Update{
    Path: "suspiciousActivity",
    Value: suspiciousFlags,
})
```

### Teacher Dashboard

Teachers can view suspicious submissions:

**GET /api/quizzes/results?quizId=quiz123**

Response includes:
```json
{
  "submissions": [
    {
      "id": "sub123",
      "studentName": "John Doe",
      "score": 85,
      "suspiciousActivity": [
        "Excessive tab switching detected",
        "Exited fullscreen mode"
      ],
      "tabSwitchCount": 5,
      "fullscreenExits": 2
    }
  ]
}
```

---

## 🎯 Complete Quiz Creation Example

### Backend Request

**POST /api/quizzes/create**

```json
{
  "title": "Midterm Exam - Database Systems",
  "description": "Covers SQL, NoSQL, and transaction management",
  "courseId": "course123",
  "totalMarks": 100,
  "passingMarks": 60,
  "duration": 90,
  "instructions": "Read all questions carefully. No external resources allowed.",
  "deadline": "2025-12-31T23:59:59Z",
  "showResultsAfterSubmit": false,
  "shuffleQuestions": true,
  "shuffleOptions": true,
  "maxAttempts": 1,
  "allowReview": true,
  
  "preventTabSwitch": true,
  "maxTabSwitches": 3,
  "requireFullscreen": true,
  "disableCopyPaste": true,
  "enableProctoring": false,
  "randomizeQuestionOrder": true,
  "timePerQuestion": 180,
  
  "allowTeacherResume": true,
  "allowTeacherExtendTime": true,
  
  "isPublished": true
}
```

---

## 🔄 Complete Quiz Flow

### Student Perspective

1. **View Available Quizzes**
   - GET `/api/quizzes/list?courseId=course123`
   - Only see published quizzes for enrolled courses

2. **Start Quiz**
   - POST `/api/quizzes/start` with `{ "quizId": "quiz123" }`
   - Receives shuffled questions (without correct answers)
   - Gets cheating prevention settings
   - Frontend enables fullscreen, tab detection, etc.

3. **Take Quiz**
   - Answer questions within time limit
   - Frontend tracks tab switches and fullscreen exits
   - Timer counts down

4. **Submit Quiz**
   - POST `/api/quizzes/submit`
   - Includes answers, tab switches, fullscreen exits
   - Backend auto-evaluates MCQ/True-False
   - Returns score (if enabled)

5. **View Results**
   - GET `/api/quizzes/results?submissionId=sub123`
   - See score, answers, explanations

### Teacher Perspective

1. **Create Quiz**
   - POST `/api/quizzes/create`
   - Enable cheating prevention features

2. **Add Questions**
   - POST `/api/quizzes/add-question` (multiple times)

3. **Publish Quiz**
   - Update `isPublished: true`

4. **Monitor Submissions**
   - GET `/api/quizzes/results?quizId=quiz123`
   - View all student submissions
   - See suspicious activity flags

5. **Resume Quiz (if needed)**
   - POST `/api/quizzes/resume`
   - Provide reason for override

---

## 🔐 Security Considerations

### Backend Validation

1. **Question Answers Hidden**: Correct answers removed from API response in `start.go`
2. **Time Validation**: Backend checks if submission is within time limit
3. **Attempt Validation**: Verifies student hasn't exceeded max attempts
4. **Enrollment Check**: Ensures student is enrolled in course
5. **Deadline Check**: Validates quiz deadline not passed

### Firestore Security Rules

```javascript
match /quiz_submissions/{submissionId} {
  allow read: if request.auth != null && (
    resource.data.studentId == request.auth.uid ||
    isTeacher() ||
    isAdmin()
  );
  
  allow create: if request.auth != null && 
    isStudent() &&
    request.resource.data.studentId == request.auth.uid;
  
  allow update: if request.auth != null && (
    (isStudent() && resource.data.studentId == request.auth.uid && resource.data.status == 'in_progress') ||
    (isTeacher() && resource.data.teacherId == request.auth.uid) ||
    isAdmin()
  );
}
```

---

## 📈 Analytics Integration

Quiz submissions update student analytics:

```go
// Update analytics after submission
analyticsRef.Update(ctx, []firestore.Update{
    {Path: "quizzesCompleted", Value: firestore.Increment(1)},
    {Path: "totalQuizScore", Value: firestore.Increment(int(totalScore))},
    {Path: "updatedAt", Value: now},
})
```

---

## 🧪 Testing Checklist

### Cheating Prevention

- [ ] Tab switching detection works
- [ ] Fullscreen enforcement works
- [ ] Copy/paste is disabled
- [ ] Questions are shuffled per student
- [ ] Timer auto-submits quiz
- [ ] Suspicious activity is flagged

### Teacher Permissions

- [ ] Teacher can resume locked quiz
- [ ] Teacher can extend time
- [ ] Resume creates audit log
- [ ] Student receives notification

### Security

- [ ] Correct answers not exposed in API
- [ ] Time limits enforced on backend
- [ ] Max attempts validated
- [ ] Only enrolled students can take quiz

---

## 🚀 API Endpoints Summary

| Endpoint | Method | Description | Auth |
|----------|--------|-------------|------|
| `/api/quizzes/create` | POST | Create quiz | Teacher/Admin |
| `/api/quizzes/list` | GET | List quizzes | All roles |
| `/api/quizzes/get` | GET | Get quiz details | All roles |
| `/api/quizzes/add-question` | POST | Add question | Teacher/Admin |
| `/api/quizzes/start` | POST | Start quiz attempt | Student |
| `/api/quizzes/submit` | POST | Submit quiz | Student |
| `/api/quizzes/results` | GET | Get results | All roles |
| `/api/quizzes/resume` | POST | Resume locked quiz | Teacher/Admin |

---

## 📝 Next Steps

1. **Implement Frontend**:
   - Quiz taking component with timer
   - Tab switching detection
   - Fullscreen enforcement
   - Copy/paste prevention

2. **Add Proctoring**:
   - Webcam integration
   - Screen recording
   - AI behavior analysis

3. **Enhanced Analytics**:
   - Quiz difficulty analysis
   - Question-level statistics
   - Cheating pattern detection

4. **Mobile Support**:
   - Responsive quiz interface
   - Mobile-specific cheating prevention
   - Offline submission support

# LMS Platform - Complete Project Structure

This document provides the complete folder structure and remaining implementation guide for the full-fledged LMS platform.

## ✅ Completed Components

### Backend (Golang)
1. **Firebase Utilities** (`/utils`)
   - ✅ Firebase Admin SDK initialization
   - ✅ Auth middleware with role-based access control
   - ✅ Response helpers (JSON, CORS, error handling)
   - ✅ Common utilities (pagination, helpers)

2. **Models** (`/models`)
   - ✅ User, Course, Quiz, Exam, Assignment, Analytics models
   - ✅ Complete request/response structures
   - ✅ Firestore-compatible annotations

3. **Auth APIs** (`/api/auth`)
   - ✅ POST `/api/auth/register` - User registration
   - ✅ GET `/api/auth/profile` - Get user profile
   - ✅ PUT `/api/auth/update` - Update profile
   - ✅ POST `/api/auth/set-role` - Admin role management

4. **Course APIs** (`/api/courses`)
   - ✅ POST `/api/courses/create` - Create course
   - ✅ GET `/api/courses/list` - List courses (role-filtered)
   - ✅ GET `/api/courses/get?id=X` - Get single course
   - ✅ PUT `/api/courses/update?id=X` - Update course
   - ✅ DELETE `/api/courses/delete?id=X` - Soft delete course
   - ✅ POST `/api/courses/enroll` - Student enrollment
   - ✅ GET `/api/courses/my-enrollments` - Student's enrollments

---

## 📋 Remaining Implementation Guide

### Quiz APIs (`/api/quizzes`)

**Required Files:**

#### `/api/quizzes/create.go`
- Handler: `CreateQuiz`
- Role: Teacher, Admin
- Creates quiz with validation
- Auto-calculate total marks based on questions

#### `/api/quizzes/list.go`
- Handler: `ListQuizzes`
- Filter by courseId, teacherId
- Students see only published quizzes

#### `/api/quizzes/get.go`
- Handler: `GetQuiz`
- Returns quiz details
- Students see questions only after starting

#### `/api/quizzes/add-question.go`
- Handler: `AddQuestion`
- Adds question to quiz
- Updates quiz totalMarks and questionsCount

#### `/api/quizzes/start.go`
- Handler: `StartQuiz`
- Creates submission with status="in_progress"
- Validates time window and attempt count
- Returns randomized questions if enabled

#### `/api/quizzes/submit.go`
- Handler: `SubmitQuiz`
- Auto-evaluates MCQ/True-False
- Calculates marks with negative marking
- Marks submission as "evaluated"

#### `/api/quizzes/results.go`
- Handler: `GetQuizResults`
- Student can view own results
- Teacher can view all results for their quiz

---

### Exam APIs (`/api/exams`)

**Required Files:**

#### `/api/exams/create.go`
- Handler: `CreateExam`
- Role: Teacher, Admin
- Scheduled exams with start/end time

#### `/api/exams/list.go`
- Handler: `ListExams`
- Filter by courseId, examType, date range

#### `/api/exams/add-question.go`
- Handler: `AddExamQuestion`
- Similar to quiz questions
- Set requiresManualEvaluation if descriptive questions

#### `/api/exams/start.go`
- Handler: `StartExam`
- Validates exam time window
- Creates submission record

#### `/api/exams/submit.go`
- Handler: `SubmitExam`
- Auto-eval MCQ, mark others for manual review
- Status: "partially_evaluated" if manual eval needed

#### `/api/exams/evaluate.go`
- Handler: `EvaluateExam`
- Role: Teacher
- Manual evaluation for descriptive answers
- Updates marksAwarded, feedback, status="evaluated"

#### `/api/exams/results.go`
- Handler: `GetExamResults`
- Returns evaluated submissions

---

### Assignment APIs (`/api/assignments`)

**Required Files:**

#### `/api/assignments/create.go`
- Handler: `CreateAssignment`
- Role: Teacher, Admin
- Includes attachments (URLs from Firebase Storage)

#### `/api/assignments/list.go`
- Handler: `ListAssignments`
- Filter by courseId
- Students see published assignments

#### `/api/assignments/submit.go`
- Handler: `SubmitAssignment`
- Role: Student
- Check late submission, calculate penalty
- File upload URLs passed from frontend

#### `/api/assignments/evaluate.go`
- Handler: `EvaluateAssignment`
- Role: Teacher
- Award marks, provide feedback

#### `/api/assignments/my-submissions.go`
- Handler: `GetMySubmissions`
- Student views own submissions

---

### Analytics APIs (`/api/analytics`)

**Required Files:**

#### `/api/analytics/student-performance.go`
- Handler: `GetStudentPerformance`
- Aggregates quiz/exam scores
- Course progress, recent activities

#### `/api/analytics/course-stats.go`
- Handler: `GetCourseStats`
- Role: Teacher, Admin
- Enrollment count, completion rate, average scores

#### `/api/analytics/quiz-stats.go`
- Handler: `GetQuizStats`
- Average score, attempt count, pass rate

---

### Admin APIs (`/api/admin`)

**Required Files:**

#### `/api/admin/users.go`
- Handler: `ListAllUsers`
- Role: Admin only
- Pagination, filter by role

#### `/api/admin/activate-user.go`
- Handler: `ActivateUser`
- Enable/disable user accounts

#### `/api/admin/reports.go`
- Handler: `GetPlatformReports`
- Total users, courses, quizzes, exams
- Activity metrics

---

## 🎨 Frontend Structure (Next.js App Router)

### Folder Structure

```
/frontend
├── /app
│   ├── layout.tsx                 # Root layout
│   ├── page.tsx                   # Landing page
│   ├── /auth
│   │   ├── /login
│   │   │   └── page.tsx
│   │   └── /register
│   │       └── page.tsx
│   ├── /dashboard
│   │   ├── layout.tsx             # Dashboard layout with sidebar
│   │   ├── /student
│   │   │   ├── page.tsx           # Student home
│   │   │   ├── /courses
│   │   │   │   ├── page.tsx       # Browse courses
│   │   │   │   └── /[id]
│   │   │   │       └── page.tsx   # Course details
│   │   │   ├── /quizzes
│   │   │   │   ├── page.tsx
│   │   │   │   └── /[id]
│   │   │   │       └── page.tsx   # Take quiz
│   │   │   ├── /exams
│   │   │   │   └── page.tsx
│   │   │   ├── /assignments
│   │   │   │   └── page.tsx
│   │   │   └── /performance
│   │   │       └── page.tsx       # Analytics
│   │   ├── /teacher
│   │   │   ├── page.tsx           # Teacher home
│   │   │   ├── /courses
│   │   │   │   ├── page.tsx       # My courses
│   │   │   │   ├── /create
│   │   │   │   │   └── page.tsx
│   │   │   │   └── /[id]
│   │   │   │       ├── page.tsx   # Edit course
│   │   │   │       ├── /quizzes
│   │   │   │       │   └── page.tsx
│   │   │   │       ├── /exams
│   │   │   │       │   └── page.tsx
│   │   │   │       └── /assignments
│   │   │   │           └── page.tsx
│   │   │   └── /evaluations
│   │   │       └── page.tsx       # Pending evaluations
│   │   └── /admin
│   │       ├── page.tsx           # Admin dashboard
│   │       ├── /users
│   │       │   └── page.tsx
│   │       ├── /courses
│   │       │   └── page.tsx
│   │       └── /reports
│   │           └── page.tsx
│   └── /api                        # (Optional) Next.js API routes if needed
├── /components
│   ├── /ui
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   └── Table.tsx
│   ├── /dashboard
│   │   ├── Sidebar.tsx
│   │   ├── Navbar.tsx
│   │   └── StatCard.tsx
│   ├── /courses
│   │   ├── CourseCard.tsx
│   │   ├── CourseList.tsx
│   │   └── MaterialUpload.tsx
│   ├── /quizzes
│   │   ├── QuizCard.tsx
│   │   ├── QuestionBuilder.tsx
│   │   └── QuizPlayer.tsx
│   └── /auth
│       └── ProtectedRoute.tsx
├── /lib
│   ├── firebase.ts                # Firebase client config
│   ├── api.ts                     # API client
│   └── auth-context.tsx           # Auth context provider
├── /public
│   └── /assets
└── package.json
```

---

### Key Frontend Files

#### `/lib/firebase.ts` - Firebase Client Configuration
```typescript
import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';
import { getStorage } from 'firebase/storage';

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
export const storage = getStorage(app);
```

#### `/lib/api.ts` - API Client
```typescript
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';

export async function apiRequest(
  endpoint: string,
  options: RequestInit = {}
) {
  const token = await getIdToken(); // From Firebase Auth
  
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers,
    },
  });
  
  return response.json();
}
```

#### `/lib/auth-context.tsx` - Auth Context
```typescript
'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { User, onAuthStateChanged } from 'firebase/auth';
import { auth } from './firebase';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  role: string | null;
}

const AuthContext = createContext<AuthContextType>({ user: null, loading: true, role: null });

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [role, setRole] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      setUser(user);
      if (user) {
        const tokenResult = await user.getIdTokenResult();
        setRole(tokenResult.claims.role as string);
      }
      setLoading(false);
    });
    return unsubscribe;
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, role }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
```

---

## 🚀 Deployment Configuration

### Vercel Deployment Steps

1. **Install Vercel CLI**
   ```bash
   npm i -g vercel
   ```

2. **Link Project**
   ```bash
   vercel link
   ```

3. **Set Environment Variables**
   In Vercel Dashboard, add:
   - `FIREBASE_PROJECT_ID`
   - `FIREBASE_PRIVATE_KEY`
   - `FIREBASE_CLIENT_EMAIL`
   - `FIREBASE_STORAGE_BUCKET`

4. **Deploy**
   ```bash
   vercel --prod
   ```

### Frontend `.env.local` Example
```env
NEXT_PUBLIC_FIREBASE_API_KEY=your_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your_project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123:web:abc
NEXT_PUBLIC_API_URL=https://your-backend.vercel.app
```

---

## 🔐 Security Best Practices

1. **Firebase Security Rules** (Firestore)
   - Users can only read/write their own documents
   - Students cannot modify grades
   - Teachers can only modify their own courses
   - Admin has full access

2. **API Security**
   - All APIs validate Firebase tokens
   - Role-based middleware enforced
   - Input validation on all endpoints
   - Rate limiting (implement with Vercel Edge Middleware)

3. **CORS**
   - Restrict origins in production
   - Update `utils/response.go` with allowed origins

4. **File Upload**
   - Frontend uploads directly to Firebase Storage
   - Get signed URLs from Firebase
   - Pass URLs to backend APIs
   - Validate file types and sizes

---

## 📊 Testing Strategy

1. **Backend Testing**
   - Test each API with Postman/Thunder Client
   - Validate role-based access
   - Test edge cases (invalid IDs, expired tokens)

2. **Frontend Testing**
   - Manual testing for each role
   - Test authentication flow
   - Verify file uploads
   - Test quiz timer functionality

3. **Integration Testing**
   - Complete user journey: Register → Login → Enroll → Take Quiz → View Results

---

## 🎯 Next Steps to Complete

1. Implement remaining Quiz APIs (6 files)
2. Implement Exam APIs (6 files)
3. Implement Assignment APIs (5 files)
4. Implement Analytics APIs (3 files)
5. Implement Admin APIs (3 files)
6. Build Next.js frontend (est. 40+ components)
7. Configure Firebase Security Rules
8. Deploy to Vercel
9. Test end-to-end

**Estimated Total**: ~70 more files to create for a fully production-ready system.

---

This platform is designed to be **scalable**, **modular**, and **production-ready** with proper architecture patterns for a real-world LMS system.

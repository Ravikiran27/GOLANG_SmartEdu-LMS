# 🎓 Full-Fledged LMS + Quiz + Exam Platform

A **production-ready**, **scalable** Learning Management System built with **Golang serverless functions**, **Next.js**, **Firebase**, and **Vercel**.

## 🚀 Features

### 👥 User Roles
- **Admin**: Platform management, user control, reports
- **Teacher**: Course creation, quiz/exam management, evaluation
- **Student**: Course enrollment, quiz/exam taking, performance tracking

### 📚 LMS Capabilities
- ✅ Course Management (CRUD)
- ✅ Material Upload (PDF, PPT, Video via Firebase Storage)
- ✅ Student Enrollment & Progress Tracking
- ✅ Role-based Access Control

### 📝 Quiz & Exam System
- ✅ Multiple Question Types (MCQ, True/False, Short Answer, Descriptive)
- ✅ Timed Assessments
- ✅ Randomized Questions
- ✅ Auto-Evaluation (MCQ/True-False)
- ✅ Manual Evaluation (Descriptive Answers)
- ✅ Negative Marking Support
- ✅ Instant/Scheduled Result Release

### 📄 Assignment Management
- ✅ File Upload Submissions
- ✅ Deadline Enforcement
- ✅ Late Submission Penalty
- ✅ Teacher Feedback System

### 📊 Analytics & Performance
- ✅ Student Performance Dashboard
- ✅ Course-wise Analytics
- ✅ Quiz/Exam Score History
- ✅ Teacher Insights

### 🛠️ Admin Panel
- ✅ User Management (Activate/Deactivate)
- ✅ Role Assignment
- ✅ Platform Statistics
- ✅ Content Monitoring

---

## 🏗️ Architecture

### Tech Stack

#### Backend
- **Language**: Golang 1.21
- **Runtime**: Vercel Serverless Functions
- **Auth**: Firebase Authentication + Custom Claims
- **Database**: Firebase Firestore
- **Storage**: Firebase Storage
- **Architecture**: Stateless, RESTful APIs

#### Frontend
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: TailwindCSS (recommended)
- **State Management**: React Context + SWR
- **Auth**: Firebase Client SDK

#### Infrastructure
- **Hosting**: Vercel (Backend + Frontend)
- **CDN**: Vercel Edge Network
- **Database**: Firebase (Auto-scaling)
- **File Storage**: Firebase Storage

---

## 📁 Project Structure

```
lms-platform/
├── api/                          # Golang serverless functions
│   ├── auth/
│   │   ├── register.go
│   │   ├── profile.go
│   │   ├── update.go
│   │   └── set-role.go
│   ├── courses/
│   │   ├── create.go
│   │   ├── list.go
│   │   ├── get.go
│   │   ├── update.go
│   │   ├── delete.go
│   │   ├── enroll.go
│   │   └── my-enrollments.go
│   ├── quizzes/
│   ├── exams/
│   ├── assignments/
│   ├── analytics/
│   └── admin/
├── utils/
│   ├── firebase.go              # Firebase Admin SDK init
│   ├── auth.go                  # Auth middleware
│   ├── response.go              # Response helpers
│   └── helpers.go               # Utility functions
├── models/
│   ├── user.go
│   ├── course.go
│   ├── quiz.go
│   ├── exam.go
│   ├── assignment.go
│   └── analytics.go
├── frontend/                     # Next.js app
│   ├── app/
│   │   ├── auth/
│   │   └── dashboard/
│   ├── components/
│   └── lib/
├── vercel.json                   # Vercel config
├── go.mod
├── FIRESTORE_SCHEMA.md
├── IMPLEMENTATION_GUIDE.md
└── README.md
```

---

## 🔑 Environment Variables

### Backend (.env)
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
FIREBASE_CLIENT_EMAIL=firebase-adminsdk@your-project.iam.gserviceaccount.com
FIREBASE_STORAGE_BUCKET=your-project.appspot.com
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_FIREBASE_API_KEY=your_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123:web:abc
NEXT_PUBLIC_API_URL=https://your-backend.vercel.app
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Firebase Project
- Vercel Account

### Backend Setup

1. **Clone repository**
   ```bash
   git clone <your-repo>
   cd lms-platform
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Configure Firebase**
   - Create Firebase project
   - Enable Authentication (Email + Google)
   - Create Firestore database
   - Enable Firebase Storage
   - Download service account key

4. **Set environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your Firebase credentials
   ```

5. **Test locally with Vercel CLI**
   ```bash
   npm i -g vercel
   vercel dev
   ```

### Frontend Setup

1. **Navigate to frontend**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env.local
   # Add Firebase client config
   ```

4. **Run dev server**
   ```bash
   npm run dev
   ```

---

## 📡 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `GET /api/auth/profile` - Get user profile
- `PUT /api/auth/update` - Update profile
- `POST /api/auth/set-role` - Set user role (Admin)

### Courses
- `POST /api/courses/create` - Create course (Teacher)
- `GET /api/courses/list` - List courses
- `GET /api/courses/get?id=X` - Get course details
- `PUT /api/courses/update?id=X` - Update course
- `DELETE /api/courses/delete?id=X` - Delete course
- `POST /api/courses/enroll` - Enroll in course (Student)
- `GET /api/courses/my-enrollments` - Get enrollments

### Quizzes
- `POST /api/quizzes/create` - Create quiz
- `POST /api/quizzes/add-question` - Add question
- `POST /api/quizzes/start` - Start quiz attempt
- `POST /api/quizzes/submit` - Submit quiz
- `GET /api/quizzes/results` - Get results

### Exams
- `POST /api/exams/create` - Create exam
- `POST /api/exams/add-question` - Add question
- `POST /api/exams/start` - Start exam
- `POST /api/exams/submit` - Submit exam
- `POST /api/exams/evaluate` - Manual evaluation
- `GET /api/exams/results` - Get results

### Assignments
- `POST /api/assignments/create` - Create assignment
- `POST /api/assignments/submit` - Submit assignment
- `POST /api/assignments/evaluate` - Evaluate submission

### Analytics
- `GET /api/analytics/student-performance` - Student metrics
- `GET /api/analytics/course-stats` - Course statistics
- `GET /api/analytics/quiz-stats` - Quiz analytics

### Admin
- `GET /api/admin/users` - List all users
- `POST /api/admin/activate-user` - Activate/deactivate user
- `GET /api/admin/reports` - Platform reports

---

## 🔒 Security

### Authentication
- Firebase Auth token validation on every request
- Custom claims for role-based access
- Secure token refresh mechanism

### Authorization
- Middleware-based role checking
- Resource ownership validation
- Admin-only endpoints protected

### Data Protection
- Firestore Security Rules enforced
- Student data isolated
- Teacher can only access own courses

### File Upload
- Direct client-to-Firebase uploads
- Signed URLs for security
- File type and size validation

---

## 📊 Database Schema

See [FIRESTORE_SCHEMA.md](./FIRESTORE_SCHEMA.md) for complete schema documentation.

### Key Collections
- `users` - User profiles
- `courses` - Course catalog
- `enrollments` - Student enrollments
- `quizzes` - Quiz definitions
- `questions` - Question bank
- `quiz_submissions` - Quiz attempts
- `exams` - Exam definitions
- `exam_submissions` - Exam attempts
- `assignments` - Assignment definitions
- `assignment_submissions` - Submissions
- `analytics` - Pre-aggregated metrics
- `notifications` - User notifications

---

## 🎯 Roadmap

### Phase 1: Core Features (Current)
- ✅ Authentication & Authorization
- ✅ Course Management
- ✅ User Roles
- 🚧 Quiz System (partial)
- 🚧 Exam System (partial)
- 🚧 Assignments (partial)

### Phase 2: Advanced Features
- [ ] Real-time notifications
- [ ] Discussion forums
- [ ] Live classes integration
- [ ] Certificate generation
- [ ] Payment integration

### Phase 3: AI & Analytics
- [ ] AI-powered grading assistance
- [ ] Plagiarism detection
- [ ] Personalized learning paths
- [ ] Predictive analytics

---

## 🤝 Contributing

This is a production-ready template. To extend:

1. Follow existing patterns in `/api` folder
2. Add models in `/models`
3. Update Firestore schema documentation
4. Implement frontend components
5. Add comprehensive tests

---

## 📝 License

MIT License - Feel free to use for educational or commercial purposes.

---

## 🙏 Acknowledgments

- **Firebase** for backend infrastructure
- **Vercel** for serverless deployment
- **Next.js** for frontend framework

---

## 📧 Support

For implementation questions, refer to [IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md).

For schema details, see [FIRESTORE_SCHEMA.md](./FIRESTORE_SCHEMA.md).

---

**Built with ❤️ for scalable education platforms**

# System Architecture Diagram

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │          Next.js Frontend (Vercel)                       │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │   Student   │  │   Teacher   │  │    Admin    │      │   │
│  │  │  Dashboard  │  │  Dashboard  │  │  Dashboard  │      │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │   │
│  │         │                │                 │             │   │
│  │         └────────────────┴─────────────────┘             │   │
│  │                          │                               │   │
│  │                 Firebase Client SDK                      │   │
│  │              (Auth + Storage + Firestore)                │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AUTHENTICATION                             │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Firebase Authentication                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │   │
│  │   │  Email   │  │  Google  │  │  Custom Claims       │  │   │
│  │   │  /Pass   │  │  OAuth   │  │  (role: admin/...)   │  │   │
│  │   └──────────┘  └──────────┘  └──────────────────────┘  │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                    JWT Token (Bearer)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      BACKEND LAYER                              │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │     Golang Serverless Functions (Vercel Edge)            │   │
│  │                                                           │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐         │   │
│  │  │    Auth    │  │   Course   │  │    Quiz    │         │   │
│  │  │    APIs    │  │    APIs    │  │    APIs    │         │   │
│  │  └────────────┘  └────────────┘  └────────────┘         │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐         │   │
│  │  │    Exam    │  │ Assignment │  │  Analytics │         │   │
│  │  │    APIs    │  │    APIs    │  │    APIs    │         │   │
│  │  └────────────┘  └────────────┘  └────────────┘         │   │
│  │                                                           │   │
│  │  ┌───────────────────────────────────────────────────┐   │   │
│  │  │          Middleware Layer                        │   │   │
│  │  │  • Token Verification                            │   │   │
│  │  │  • Role-Based Access Control (RBAC)              │   │   │
│  │  │  • CORS Management                               │   │   │
│  │  │  • Response Standardization                      │   │   │
│  │  └───────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        DATA LAYER                               │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Firebase Firestore (NoSQL)                   │   │
│  │                                                           │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │   │
│  │  │  users   │  │ courses  │  │ quizzes  │  │  exams   │ │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │   │
│  │  │questions │  │quiz_sub  │  │exam_sub  │  │analytics │ │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │   │
│  │  │assignments│ │assign_sub│  │enrollment│  │  notifs  │ │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │   │
│  │                                                           │   │
│  │  ┌───────────────────────────────────────────────────┐   │   │
│  │  │          Security Rules Engine                   │   │   │
│  │  │  • Role-based document access                    │   │   │
│  │  │  • Ownership validation                          │   │   │
│  │  │  • Field-level permissions                       │   │   │
│  │  └───────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      STORAGE LAYER                              │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Firebase Storage (CDN)                       │   │
│  │                                                           │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │   │
│  │  │   Course     │  │  Assignment  │  │    User      │   │   │
│  │  │  Materials   │  │ Submissions  │  │   Photos     │   │   │
│  │  │ (PDF,Video)  │  │  (Files)     │  │   (Images)   │   │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘   │   │
│  │                                                           │   │
│  │  ┌───────────────────────────────────────────────────┐   │   │
│  │  │          Storage Security Rules                  │   │   │
│  │  │  • File type validation                          │   │   │
│  │  │  • Size limits enforcement                       │   │   │
│  │  │  • Upload permissions by role                    │   │   │
│  │  └───────────────────────────────────────────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Flow Diagrams

### 1. User Authentication Flow

```
┌────────┐
│ User   │
└───┬────┘
    │ 1. Enter email/password
    ▼
┌────────────────┐
│ Next.js Login  │
│     Page       │
└───┬────────────┘
    │ 2. signInWithEmailAndPassword()
    ▼
┌──────────────────┐
│ Firebase Auth    │
│                  │
│ • Validates      │
│ • Generates JWT  │
│ • Adds claims    │
└───┬──────────────┘
    │ 3. Returns token + user
    ▼
┌────────────────┐
│ Auth Context   │
│ (Frontend)     │
│                │
│ Stores:        │
│ • User object  │
│ • JWT token    │
│ • Role         │
└───┬────────────┘
    │ 4. Redirect to dashboard
    ▼
┌────────────────┐
│ Role-based     │
│ Dashboard      │
│                │
│ /student       │
│ /teacher       │
│ /admin         │
└────────────────┘
```

### 2. API Request Flow

```
┌────────┐
│ User   │
└───┬────┘
    │ 1. Click "Create Course"
    ▼
┌────────────────┐
│ React Component│
└───┬────────────┘
    │ 2. api.courses.create(data)
    ▼
┌────────────────┐
│ API Client     │
│ (lib/api.ts)   │
│                │
│ • Get JWT token│
│ • Add to header│
└───┬────────────┘
    │ 3. POST /api/courses/create
    │    Authorization: Bearer <token>
    ▼
┌────────────────┐
│ Vercel Edge    │
│ (Load Balancer)│
└───┬────────────┘
    │ 4. Route to function
    ▼
┌────────────────┐
│ Go Handler     │
│ create.go      │
└───┬────────────┘
    │ 5. AuthMiddleware()
    ▼
┌────────────────┐
│ Auth Middleware│
│                │
│ • Verify token │ ─────────────┐
│ • Extract UID  │              │
│ • Check role   │              │ 6. Verify with
└───┬────────────┘              │    Firebase Admin SDK
    │ 7. Valid? Continue        ▼
    ▼                     ┌──────────────┐
┌────────────────┐        │ Firebase     │
│ Business Logic │        │ Auth Service │
│                │        └──────────────┘
│ • Validate data│
│ • Create course│
└───┬────────────┘
    │ 8. Save to Firestore
    ▼
┌────────────────┐
│ Firestore      │
│                │
│ • Check rules  │
│ • Save doc     │
│ • Return ID    │
└───┬────────────┘
    │ 9. Course created
    ▼
┌────────────────┐
│ Response       │
│ {              │
│   success: true│
│   data: {...}  │
│ }              │
└───┬────────────┘
    │ 10. Send to client
    ▼
┌────────────────┐
│ React Component│
│                │
│ • Update UI    │
│ • Show success │
└────────────────┘
```

### 3. Quiz Submission Flow (Future)

```
Student                Teacher                Backend                Firestore
   │                      │                      │                      │
   │ 1. Start Quiz        │                      │                      │
   ├──────────────────────────────────────────────────────────────────>│
   │                      │                      │  2. Create          │
   │                      │                      │  quiz_submission     │
   │                      │                      │  (status: in_progress│
   │<─────────────────────────────────────────────────────────────────┤
   │ 3. Display questions │                      │                      │
   │                      │                      │                      │
   │ 4. Answer questions  │                      │                      │
   │ (client-side timer)  │                      │                      │
   │                      │                      │                      │
   │ 5. Submit answers    │                      │                      │
   ├──────────────────────────────────────────────────────────────────>│
   │                      │                      │ 6. Auto-evaluate     │
   │                      │                      │    MCQ/True-False    │
   │                      │                      │                      │
   │                      │                      │ 7. Update submission │
   │                      │                      │    (status: evaluated│
   │<─────────────────────────────────────────────────────────────────┤
   │ 8. Show results      │                      │                      │
   │ (if enabled)         │                      │                      │
   │                      │                      │                      │
   │                      │ 9. View all results  │                      │
   │                      ├──────────────────────────────────────────>│
   │                      │                      │ 10. Aggregate data   │
   │                      │<─────────────────────────────────────────┤
   │                      │ 11. Display analytics│                      │
```

### 4. File Upload Flow

```
Student                Frontend               Firebase Storage        Backend API
   │                      │                      │                      │
   │ 1. Select file       │                      │                      │
   │──────────────────────>│                      │                      │
   │                      │ 2. Upload to Storage │                      │
   │                      │──────────────────────>│                      │
   │                      │                      │ 3. Security rules    │
   │                      │                      │    validate          │
   │                      │<─────────────────────┤ 4. Return URL        │
   │                      │                      │                      │
   │                      │ 5. POST /api/assignments/submit             │
   │                      │     { fileUrl: "..." }                      │
   │                      ├──────────────────────────────────────────────>│
   │                      │                      │                      │
   │                      │                      │ 6. Save submission   │
   │                      │                      │    with file URL     │
   │                      │<─────────────────────────────────────────────┤
   │<─────────────────────┤ 7. Success           │                      │
   │ 8. Show confirmation │                      │                      │
```

---

## Security Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                          │
│                                                              │
│  Layer 1: Client-side                                       │
│  ┌────────────────────────────────────────────────────┐     │
│  │  • Firebase Auth required for all actions          │     │
│  │  • Protected routes redirect to login              │     │
│  │  • Role-based UI rendering                         │     │
│  └────────────────────────────────────────────────────┘     │
│                          │                                   │
│                          ▼                                   │
│  Layer 2: API Gateway (Vercel Edge)                         │
│  ┌────────────────────────────────────────────────────┐     │
│  │  • HTTPS enforced                                  │     │
│  │  • CORS validation                                 │     │
│  │  • Rate limiting (future)                          │     │
│  └────────────────────────────────────────────────────┘     │
│                          │                                   │
│                          ▼                                   │
│  Layer 3: Backend (Go Functions)                            │
│  ┌────────────────────────────────────────────────────┐     │
│  │  • JWT token verification                          │     │
│  │  • Custom claims validation (role check)           │     │
│  │  • Ownership verification (teacherId, studentId)   │     │
│  │  • Input validation                                │     │
│  └────────────────────────────────────────────────────┘     │
│                          │                                   │
│                          ▼                                   │
│  Layer 4: Database (Firestore Security Rules)               │
│  ┌────────────────────────────────────────────────────┐     │
│  │  • Request.auth validation                         │     │
│  │  • Role-based read/write permissions               │     │
│  │  • Field-level access control                      │     │
│  │  • Resource ownership checks                       │     │
│  └────────────────────────────────────────────────────┘     │
│                          │                                   │
│                          ▼                                   │
│  Layer 5: Storage (Firebase Storage Rules)                  │
│  ┌────────────────────────────────────────────────────┐     │
│  │  • File type validation                            │     │
│  │  • Size limit enforcement                          │     │
│  │  • Upload permission by role                       │     │
│  │  • Download permission by ownership                │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Setup                          │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │                 DNS / Domain                       │     │
│  │  • lms.yourdomain.com → Frontend (Vercel)          │     │
│  │  • api.yourdomain.com → Backend (Vercel)           │     │
│  └────────┬───────────────────────────────────────────┘     │
│           │                                                  │
│           ▼                                                  │
│  ┌────────────────────────────────────────────────────┐     │
│  │           Vercel Edge Network (CDN)                │     │
│  │  • 200+ locations globally                         │     │
│  │  • Automatic SSL certificates                      │     │
│  │  • DDoS protection                                 │     │
│  └────────┬───────────────────┬───────────────────────┘     │
│           │                   │                             │
│           ▼                   ▼                             │
│  ┌─────────────────┐ ┌─────────────────┐                   │
│  │  Next.js App    │ │  Go Functions   │                   │
│  │  (Frontend)     │ │  (Backend)      │                   │
│  │                 │ │                 │                   │
│  │  • Auto-scaled  │ │  • Auto-scaled  │                   │
│  │  • Serverless   │ │  • Serverless   │                   │
│  └─────────────────┘ └────────┬────────┘                   │
│                               │                             │
│                               ▼                             │
│  ┌────────────────────────────────────────────────────┐     │
│  │         Firebase (Google Cloud Platform)           │     │
│  │                                                     │     │
│  │  ┌─────────────┐  ┌─────────────┐  ┌───────────┐ │     │
│  │  │    Auth     │  │  Firestore  │  │  Storage  │ │     │
│  │  │  (Global)   │  │ (Regional)  │  │ (Regional)│ │     │
│  │  └─────────────┘  └─────────────┘  └───────────┘ │     │
│  │                                                     │     │
│  │  • Multi-region replication (optional)             │     │
│  │  • Automatic backups                               │     │
│  │  • 99.95% SLA                                      │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

This architecture provides:

✅ **Scalability**: Auto-scales from 0 to thousands of users  
✅ **Reliability**: Multi-layer redundancy, 99.95% uptime  
✅ **Security**: Defense-in-depth with 5 security layers  
✅ **Performance**: Edge caching, CDN, indexed queries  
✅ **Cost-Effective**: Pay only for what you use  

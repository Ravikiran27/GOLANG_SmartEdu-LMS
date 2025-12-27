# Architecture Decision Record (ADR)

## System Architecture - LMS Platform

### 1. Context & Goals

**Objective**: Build a production-ready LMS + Quiz + Exam platform that is:
- Scalable to thousands of concurrent users
- Cost-effective (serverless, pay-per-use)
- Easy to maintain and extend
- Deployable on Vercel only
- Secure and compliant

---

### 2. Technology Decisions

#### 2.1 Backend: Golang Serverless Functions

**Decision**: Use Golang for serverless API functions on Vercel.

**Rationale**:
- ✅ **Performance**: Go's compiled nature provides ~10x better cold start than Node.js
- ✅ **Concurrency**: Built-in goroutines handle concurrent requests efficiently
- ✅ **Type Safety**: Strong typing reduces runtime errors
- ✅ **Memory Efficient**: Lower memory footprint = lower costs
- ✅ **Vercel Support**: Native Go runtime available

**Alternatives Considered**:
- ❌ Node.js: Slower cold starts, higher memory usage
- ❌ Python: Not optimized for serverless performance
- ❌ Traditional server: Not compatible with Vercel constraint

**Trade-offs**:
- Slightly longer development time vs Node.js
- Smaller ecosystem for specific libraries

---

#### 2.2 Database: Firebase Firestore

**Decision**: Use Firestore as primary database.

**Rationale**:
- ✅ **NoSQL Flexibility**: Schema evolution without migrations
- ✅ **Real-time Sync**: Built-in real-time updates (future feature)
- ✅ **Scalability**: Auto-scales with zero configuration
- ✅ **Security Rules**: Declarative client-side security
- ✅ **Firebase Integration**: Seamless with Auth and Storage
- ✅ **Free Tier**: Generous for development/testing

**Alternatives Considered**:
- ❌ PostgreSQL/MySQL: Requires separate hosting, connection pooling issues in serverless
- ❌ MongoDB Atlas: Similar benefits but less Firebase ecosystem integration
- ❌ DynamoDB: AWS lock-in, more complex to configure

**Data Model Strategy**:
- **Denormalization**: Store teacherName, courseTitle in related docs for fast reads
- **Flat Collections**: Avoid sub-collections for easier querying
- **Pre-aggregated Analytics**: Computed metrics in separate collection
- **Soft Deletes**: isDeleted flag instead of actual deletion

---

#### 2.3 Authentication: Firebase Authentication

**Decision**: Use Firebase Auth with custom claims.

**Rationale**:
- ✅ **Built-in Security**: Industry-standard OAuth, JWT tokens
- ✅ **Multi-provider**: Email, Google, future: Microsoft, SAML
- ✅ **Custom Claims**: Role-based access (admin/teacher/student)
- ✅ **Token Refresh**: Automatic token renewal
- ✅ **No Password Storage**: Security handled by Firebase

**Custom Claims Implementation**:
```javascript
// Set during registration
claims = {
  role: "student" | "teacher" | "admin"
}
```

**Middleware Flow**:
```
Client Request → Vercel Edge → Go Handler → Verify Token → Check Role → Process
```

---

#### 2.4 Frontend: Next.js (App Router)

**Decision**: Use Next.js 14 with App Router.

**Rationale**:
- ✅ **React Server Components**: Better performance, SEO
- ✅ **File-based Routing**: Clean dashboard structure per role
- ✅ **Built-in API Routes**: Optional BFF (Backend-for-Frontend)
- ✅ **Vercel Optimization**: Zero-config deployment
- ✅ **TypeScript Support**: Type safety across stack

**Folder Structure Pattern**:
```
/app
  /dashboard
    /student    <- Student-specific pages
    /teacher    <- Teacher-specific pages
    /admin      <- Admin-specific pages
```

Each role has isolated routes with layout-based navigation.

---

#### 2.5 File Storage: Firebase Storage

**Decision**: Use Firebase Storage for course materials, submissions.

**Rationale**:
- ✅ **Direct Upload**: Client uploads directly (no backend bottleneck)
- ✅ **Security Rules**: Declarative access control
- ✅ **CDN**: Automatic global distribution
- ✅ **Signed URLs**: Temporary access for authenticated users
- ✅ **Integration**: Works seamlessly with Firestore metadata

**Upload Flow**:
```
Client → Firebase Storage (with auth) → Get URL → Send URL to Go API → Store in Firestore
```

---

### 3. Architectural Patterns

#### 3.1 API Design: Function-per-Endpoint

**Pattern**: Each API endpoint is a separate Go file/function.

```
/api/courses/create.go   <- POST /api/courses/create
/api/courses/list.go     <- GET /api/courses/list
```

**Benefits**:
- 🎯 **Modularity**: Easy to test, modify, deploy independently
- 🎯 **Cold Start Optimization**: Only load required code
- 🎯 **Team Scalability**: Different devs can work on different endpoints

---

#### 3.2 Middleware: Composable Auth

**Pattern**: Wrap handlers with middleware for auth and role checks.

```go
utils.AuthMiddleware(handler, "teacher", "admin")(w, r)
```

**Benefits**:
- 🔒 **Centralized Security**: Auth logic in one place
- 🔒 **Reusable**: Apply to any endpoint
- 🔒 **Flexible**: Support multiple roles per endpoint

---

#### 3.3 Response Standardization

**Pattern**: Consistent JSON response structure.

```json
{
  "success": true,
  "data": {},
  "message": "",
  "error": ""
}
```

**Benefits**:
- 📡 **Predictable**: Frontend knows exact response shape
- 📡 **Error Handling**: Consistent error structure
- 📡 **Type Safety**: Easy to type in TypeScript

---

#### 3.4 Pagination: Cursor-based

**Pattern**: Use Firestore's native pagination with cursor/offset.

```go
query.OrderBy("createdAt", Desc).Limit(20).Offset(page * pageSize)
```

**Benefits**:
- ⚡ **Performance**: Indexed queries are fast
- ⚡ **Consistency**: No skipped/duplicate results

---

### 4. Security Architecture

#### 4.1 Defense in Depth

**Layers**:
1. **Client-side**: Firebase Auth token required
2. **Edge**: Vercel firewall, rate limiting
3. **API**: Token verification, role validation
4. **Database**: Firestore security rules
5. **Storage**: Firebase Storage rules

#### 4.2 Role-Based Access Control (RBAC)

**Hierarchy**:
```
Admin (full access)
  ├── Teacher (courses, quizzes, exams they own)
  └── Student (enrolled courses, own submissions)
```

**Enforcement Points**:
- API middleware checks JWT claims
- Firestore rules validate ownership
- Frontend conditionally renders UI

---

### 5. Scalability Considerations

#### 5.1 Serverless Benefits

**Auto-scaling**:
- Vercel automatically scales functions based on traffic
- No server management required
- Pay only for execution time

**Limits to Consider**:
- **Function Timeout**: 10s max (Vercel hobby), 60s (pro)
- **Memory**: 1024MB default
- **Concurrency**: Auto-scaled by Vercel

**Optimization Strategies**:
- Keep functions focused (single responsibility)
- Avoid heavy computation in sync APIs
- Use Firestore batch operations for bulk updates

---

#### 5.2 Database Indexing

**Critical Indexes**:
```
- users: (role, isActive)
- courses: (teacherId, isPublished, createdAt)
- enrollments: (studentId, courseId) [composite]
- quiz_submissions: (quizId, studentId)
- exam_submissions: (examId, studentId)
```

**Impact**:
- Queries under 100ms even with 100k+ documents
- Prevents full collection scans

---

#### 5.3 Caching Strategy

**Future Enhancement** (not in MVP):
- Add Redis/Upstash for session caching
- Cache expensive aggregations (analytics)
- Use SWR on frontend for optimistic UI

---

### 6. Cost Optimization

#### 6.1 Serverless Cost Model

**Vercel**:
- Free tier: 100GB-hrs/month (sufficient for MVP)
- Pro: $20/month for production use

**Firebase**:
- Free tier:
  - 1GB storage
  - 50K reads/day
  - 20K writes/day
- Pay-as-you-go after limits

**Estimated Cost** (1000 active users):
- Vercel: $20/month
- Firebase: ~$50/month
- **Total**: ~$70/month

vs Traditional Hosting:
- VPS: $40/month minimum
- Database: $20/month
- CDN: $10/month
- Managed Backups: $15/month
- **Total**: ~$85/month (less scalable)

---

### 7. Development Workflow

#### 7.1 Local Development

```bash
# Backend
vercel dev  # Runs Go functions locally

# Frontend
cd frontend && npm run dev
```

**Hot Reload**: Both support file watching.

---

#### 7.2 Deployment Pipeline

```
Git Push → GitHub Actions (optional CI) → Vercel Auto-Deploy → Production
```

**Environments**:
- `main` branch → Production
- `develop` branch → Preview deployment

---

### 8. Testing Strategy

#### 8.1 Backend Testing

**Unit Tests**:
```go
// Test auth middleware
// Test response helpers
// Test model validation
```

**Integration Tests**:
```bash
# Test actual API calls with test Firebase project
```

---

#### 8.2 Frontend Testing

**Component Tests**: Jest + React Testing Library
**E2E Tests**: Playwright/Cypress (future)

---

### 9. Monitoring & Observability

#### 9.1 Vercel Analytics

- Function execution time
- Error rates
- Cold start frequency

#### 9.2 Firebase Console

- Firestore usage metrics
- Auth activity
- Storage bandwidth

#### 9.3 Custom Logging

```go
log.Printf("[ERROR] %s: %v", context, err)
```

Logs visible in Vercel dashboard.

---

### 10. Future Enhancements

#### 10.1 Real-time Features

**Potential**: Use Firestore real-time listeners for:
- Live quiz leaderboard
- Instant grade notifications
- Chat/discussion forums

**Implementation**: Already supported by Firestore, add frontend listeners.

---

#### 10.2 AI Integration

**Use Cases**:
- Auto-generate quiz questions from course materials
- Plagiarism detection in assignments
- Personalized learning recommendations

**Tech**: OpenAI API, Vertex AI (future modules)

---

#### 10.3 Multi-tenancy

**Approach**: Add `organizationId` to all collections for SaaS model.

---

### 11. Compliance & Data Privacy

#### 11.1 GDPR/CCPA

**Features**:
- User data export (API endpoint)
- Account deletion (soft + hard delete options)
- Consent management (future)

#### 11.2 Education Standards

**FERPA Compliance** (US Education):
- Role-based access ensures only authorized users see student data
- Audit logs for data access (future)

---

## Conclusion

This architecture balances:
- ✅ **Performance**: Serverless Go functions
- ✅ **Scalability**: Auto-scaling infrastructure
- ✅ **Security**: Multi-layered auth and authorization
- ✅ **Cost**: Pay-per-use model
- ✅ **Developer Experience**: Clean separation of concerns

**Total Implementation Time Estimate**: 4-6 weeks for full production-ready system with one senior full-stack developer.

---

**Document Version**: 1.0  
**Last Updated**: 2025-01-27  
**Status**: Architecture Finalized, Implementation In Progress

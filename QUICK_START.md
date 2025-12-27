# Quick Start Guide - LMS Platform

Get your LMS platform running locally in under 10 minutes.

---

## Prerequisites

- [x] Node.js 18+ installed
- [x] Go 1.21+ installed
- [x] Firebase account created
- [x] Vercel account created (for deployment)

---

## Step 1: Firebase Setup (5 minutes)

### 1.1 Create Firebase Project

1. Go to https://console.firebase.google.com
2. Click "Add Project"
3. Name: `lms-platform-dev`
4. Disable Google Analytics (optional)
5. Click "Create Project"

### 1.2 Enable Services

**Authentication:**
1. Left sidebar → Authentication → Get Started
2. Enable "Email/Password"
3. (Optional) Enable Google sign-in

**Firestore:**
1. Left sidebar → Firestore Database → Create database
2. Select "Test mode" (for development)
3. Choose region: `us-central1`
4. Click "Enable"

**Storage:**
1. Left sidebar → Storage → Get Started
2. Use default rules
3. Click "Done"

### 1.3 Get Credentials

**For Backend (Service Account):**
1. Project Settings (gear icon) → Service Accounts
2. Click "Generate new private key"
3. Save the JSON file

**For Frontend (Web Config):**
1. Project Settings → General
2. Under "Your apps", click Web icon `</>`
3. Register app: `lms-frontend`
4. Copy the `firebaseConfig` object

---

## Step 2: Backend Setup (2 minutes)

### 2.1 Install Dependencies

```bash
cd d:\LMS
go mod download
```

### 2.2 Configure Environment

Create `.env` in project root:

```env
FIREBASE_PROJECT_ID=lms-platform-dev
FIREBASE_CLIENT_EMAIL=firebase-adminsdk-xxxxx@lms-platform-dev.iam.gserviceaccount.com
FIREBASE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nYour key from JSON file\n-----END PRIVATE KEY-----\n"
FIREBASE_STORAGE_BUCKET=lms-platform-dev.appspot.com
```

**Get values from the service account JSON:**
- `project_id` → FIREBASE_PROJECT_ID
- `client_email` → FIREBASE_CLIENT_EMAIL
- `private_key` → FIREBASE_PRIVATE_KEY (keep the \n characters)

### 2.3 Start Backend

```bash
# Install Vercel CLI
npm install -g vercel

# Run locally
vercel dev
```

Backend will start on http://localhost:3000

---

## Step 3: Frontend Setup (3 minutes)

### 3.1 Create Frontend

```bash
cd d:\LMS\frontend
npm install
```

### 3.2 Configure Environment

Create `frontend/.env.local`:

```env
NEXT_PUBLIC_FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXX
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=lms-platform-dev.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=lms-platform-dev
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=lms-platform-dev.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abcdefg

NEXT_PUBLIC_API_URL=http://localhost:3000
```

Use values from Firebase web config.

### 3.3 Start Frontend

```bash
npm run dev
```

Frontend will start on http://localhost:3001

---

## Step 4: Create First Admin User

### 4.1 Register via API

```bash
# Using PowerShell
$body = @{
    email = "admin@example.com"
    password = "Admin123!"
    displayName = "Admin User"
    role = "admin"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/api/auth/register" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body
```

### 4.2 Set Custom Claims

The registration already sets the role, but to verify:

1. Go to Firebase Console → Authentication → Users
2. Find the user you just created
3. Note the UID

Create `set-admin.js` in project root:

```javascript
const admin = require('firebase-admin');
const serviceAccount = require('./path-to-your-service-account.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

const uid = 'PASTE_UID_HERE';

admin.auth().setCustomUserClaims(uid, { role: 'admin' })
  .then(() => {
    console.log('✅ Admin role set');
    process.exit(0);
  })
  .catch(err => {
    console.error('❌ Error:', err);
    process.exit(1);
  });
```

Run:
```bash
node set-admin.js
```

---

## Step 5: Test the Application

### 5.1 Login

1. Open http://localhost:3001/auth/login
2. Email: `admin@example.com`
3. Password: `Admin123!`
4. Click "Login"

### 5.2 Test Course Creation

```bash
# Get auth token from browser (F12 → Application → IndexedDB → firebaseLocalStorageDb)
# Or login and check Network tab for the token

$token = "PASTE_FIREBASE_TOKEN_HERE"

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

$body = @{
    title = "Introduction to Programming"
    description = "Learn the basics of programming"
    category = "Computer Science"
    difficulty = "beginner"
    syllabus = "Week 1: Variables\nWeek 2: Functions"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/api/courses/create" `
    -Method POST `
    -Headers $headers `
    -Body $body
```

### 5.3 Verify in Firestore

1. Firebase Console → Firestore Database
2. Check `courses` collection
3. You should see the newly created course

---

## Common Issues & Solutions

### Issue: "Firebase Admin SDK not initialized"

**Solution:** Check your `.env` file:
- Ensure FIREBASE_PRIVATE_KEY includes `\n` characters
- Verify project ID matches Firebase Console

### Issue: "CORS error"

**Solution:** Make sure:
- Backend is running on port 3000
- Frontend is running on port 3001
- `NEXT_PUBLIC_API_URL` in frontend matches backend URL

### Issue: "Token expired"

**Solution:** 
- Firebase tokens expire after 1 hour
- Logout and login again
- Or implement auto-refresh (see auth-context.tsx)

### Issue: Port already in use

**Solution:**
```bash
# Kill process on port 3000 (backend)
Get-Process -Id (Get-NetTCPConnection -LocalPort 3000).OwningProcess | Stop-Process

# Kill process on port 3001 (frontend)
Get-Process -Id (Get-NetTCPConnection -LocalPort 3001).OwningProcess | Stop-Process
```

---

## Project Structure Overview

```
d:\LMS\
├── api/              # Golang serverless functions
│   ├── auth/         # ✅ Authentication APIs (complete)
│   └── courses/      # ✅ Course APIs (complete)
├── utils/            # ✅ Firebase, Auth, Response utilities
├── models/           # ✅ Data models
├── frontend/         # Next.js application
│   ├── app/          # ✅ App router pages
│   ├── lib/          # ✅ Firebase config, API client, auth
│   └── components/   # ✅ Reusable components
├── vercel.json       # ✅ Vercel configuration
├── go.mod            # ✅ Go dependencies
└── .env              # Your Firebase credentials
```

---

## Next Steps After Setup

1. **Read Documentation:**
   - `README.md` - Project overview
   - `ARCHITECTURE.md` - Design decisions
   - `FIRESTORE_SCHEMA.md` - Database structure
   - `IMPLEMENTATION_GUIDE.md` - Next features to build

2. **Implement Remaining APIs:**
   - Quizzes (6 files)
   - Exams (6 files)
   - Assignments (5 files)
   - Analytics (3 files)
   - Admin (3 files)

3. **Build Frontend:**
   - Student dashboard
   - Teacher dashboard
   - Admin dashboard

4. **Deploy to Production:**
   - Follow `DEPLOYMENT.md`
   - Set up custom domain
   - Configure production Firebase project

---

## Testing Checklist

- [ ] Backend starts without errors
- [ ] Frontend loads at localhost:3001
- [ ] Can register a new user
- [ ] Can login with credentials
- [ ] Can create a course (admin/teacher)
- [ ] Course appears in Firestore
- [ ] Can list courses via API
- [ ] Authentication token is validated

---

## Development Workflow

```bash
# Terminal 1 - Backend
cd d:\LMS
vercel dev

# Terminal 2 - Frontend
cd d:\LMS\frontend
npm run dev

# Terminal 3 - Testing
# Use for curl/Invoke-RestMethod commands
```

---

## Useful Commands

### Backend

```bash
# Test API endpoint
Invoke-RestMethod -Uri "http://localhost:3000/api/auth/profile" `
    -Headers @{"Authorization"="Bearer YOUR_TOKEN"}

# Check Go syntax
go fmt ./...

# Build (for testing)
go build ./...
```

### Frontend

```bash
# Install new package
npm install package-name

# Type check
npm run type-check

# Lint
npm run lint
```

### Firebase

```bash
# Deploy Firestore rules
firebase deploy --only firestore:rules

# Deploy Storage rules
firebase deploy --only storage:rules

# View Firestore data
firebase firestore:get users/USER_ID
```

---

## Support Resources

- **Firebase Docs:** https://firebase.google.com/docs
- **Next.js Docs:** https://nextjs.org/docs
- **Vercel Docs:** https://vercel.com/docs
- **Go Firebase SDK:** https://pkg.go.dev/firebase.google.com/go/v4

---

**You're all set! Start building your LMS platform. 🚀**

If you encounter issues, check:
1. Environment variables are correct
2. Firebase services are enabled
3. Ports 3000 and 3001 are available
4. Node.js and Go are properly installed

Happy coding!

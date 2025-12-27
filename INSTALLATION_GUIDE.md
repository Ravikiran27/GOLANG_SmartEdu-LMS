# LMS Platform - Complete Installation Guide

## System Requirements

- **Node.js** 18.x or higher
- **Go** 1.21 or higher
- **Firebase** project (free tier works)
- **Vercel** account (free tier works)
- **Git** installed

---

## 📋 Pre-Installation Checklist

- [ ] Node.js installed (`node --version`)
- [ ] Go installed (`go version`)
- [ ] Firebase CLI installed (`npm install -g firebase-tools`)
- [ ] Vercel CLI installed (`npm install -g vercel`)
- [ ] Git configured
- [ ] Firebase project created

---

## 🔥 Part 1: Firebase Setup

### 1. Create Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Click "Add project"
3. Enter project name: `lms-platform`
4. Disable Google Analytics (optional)
5. Click "Create project"

### 2. Enable Authentication

1. In Firebase Console, click **Authentication**
2. Click "Get started"
3. Enable **Email/Password** sign-in method
4. Save

### 3. Create Firestore Database

1. Click **Firestore Database**
2. Click "Create database"
3. Select **Production mode** (we'll deploy rules later)
4. Choose location (e.g., `us-central`)
5. Click "Enable"

### 4. Enable Storage

1. Click **Storage**
2. Click "Get started"
3. Use default security rules (we'll deploy custom rules)
4. Choose same location as Firestore
5. Click "Done"

### 5. Get Firebase Credentials

1. Click ⚙️ **Settings** → **Project settings**
2. Scroll to "Your apps"
3. Click **Web app** icon (</>) 
4. Register app: `LMS Frontend`
5. Copy the config object:

```javascript
const firebaseConfig = {
  apiKey: "AIza...",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project-id",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "123456789",
  appId: "1:123456789:web:abc123"
};
```

### 6. Generate Service Account

1. Go to **Project settings** → **Service accounts**
2. Click "Generate new private key"
3. Save the JSON file as `serviceAccountKey.json`
4. **Keep this file secure!**

---

## 🖥️ Part 2: Backend Setup (Golang)

### 1. Clone Repository

```bash
git clone <your-repo-url>
cd LMS
```

### 2. Install Go Dependencies

```bash
cd backend
go mod download
```

### 3. Configure Environment

Create `.env` in root directory:

```bash
# Firebase Admin SDK
FIREBASE_CREDENTIALS_JSON='{"type":"service_account","project_id":"your-project",...}'

# Or use file path (choose one)
FIREBASE_CREDENTIALS_PATH=./serviceAccountKey.json

# API Configuration
PORT=3000
CORS_ORIGIN=http://localhost:3000
```

**Option A: Use JSON directly**
```bash
# Copy content from serviceAccountKey.json
FIREBASE_CREDENTIALS_JSON='<paste entire JSON here>'
```

**Option B: Use file path**
```bash
# Place serviceAccountKey.json in root directory
FIREBASE_CREDENTIALS_PATH=./serviceAccountKey.json
```

### 4. Test Backend Locally

```bash
# Install Vercel CLI
npm install -g vercel

# Run locally
vercel dev
```

Visit: `http://localhost:3000/api/auth/register`

---

## 🎨 Part 3: Frontend Setup (Next.js)

### 1. Navigate to Frontend

```bash
cd frontend
```

### 2. Run Setup Script

**Windows:**
```bash
setup.bat
```

**Mac/Linux:**
```bash
chmod +x setup.sh
./setup.sh
```

**Or manually:**
```bash
npm install
cp .env.example .env.local
```

### 3. Configure Environment

Edit `frontend/.env.local`:

```env
NEXT_PUBLIC_FIREBASE_API_KEY=AIza...
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abc123
NEXT_PUBLIC_API_URL=http://localhost:3000/api
```

### 4. Start Development Server

```bash
npm run dev
```

Visit: `http://localhost:3000`

---

## 🔐 Part 4: Deploy Security Rules

### 1. Login to Firebase

```bash
firebase login
```

### 2. Initialize Firebase in Project

```bash
# From root directory
firebase init
```

Select:
- ✅ Firestore
- ✅ Storage

Choose existing project: `your-project-id`

### 3. Deploy Firestore Rules

```bash
firebase deploy --only firestore:rules
```

This deploys `firestore.rules` with RBAC security.

### 4. Deploy Storage Rules

```bash
firebase deploy --only storage:rules
```

This deploys `storage.rules` with file validation.

---

## 🚀 Part 5: Deploy to Production

### 1. Deploy Backend to Vercel

```bash
# From root directory
vercel
```

Follow prompts:
- Link to existing project or create new
- Set root directory: `.`
- Build command: (leave empty)
- Output directory: (leave empty)

**Add Environment Variables in Vercel:**

1. Go to Vercel Dashboard → Project → Settings → Environment Variables
2. Add:
   - `FIREBASE_CREDENTIALS_JSON` = (paste JSON from serviceAccountKey.json)

### 2. Deploy Frontend to Vercel

```bash
cd frontend
vercel
```

**Add Environment Variables in Vercel:**

1. Go to Vercel Dashboard → Project → Settings → Environment Variables
2. Add all `NEXT_PUBLIC_*` variables from `.env.local`
3. Update `NEXT_PUBLIC_API_URL` to your backend URL:
   - `https://your-backend.vercel.app/api`

### 3. Update CORS Settings

Update backend `.env` or Vercel environment:

```bash
CORS_ORIGIN=https://your-frontend.vercel.app
```

Redeploy backend:
```bash
vercel --prod
```

---

## ✅ Part 6: Verification

### 1. Test Authentication

1. Visit your frontend URL
2. Click "Register"
3. Create account with email/password
4. Check Firebase Console → Authentication → Users

### 2. Test API

```bash
# Health check
curl https://your-backend.vercel.app/api/health

# Register user
curl -X POST https://your-backend.vercel.app/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "name": "Test User",
    "role": "student"
  }'
```

### 3. Create Admin User

Use Firebase Console to set custom claims:

```bash
# Install Firebase Admin
npm install -g firebase-admin

# Create script to set admin
node -e "
const admin = require('firebase-admin');
const serviceAccount = require('./serviceAccountKey.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

admin.auth().getUserByEmail('admin@example.com')
  .then(user => {
    return admin.auth().setCustomUserClaims(user.uid, { role: 'admin' });
  })
  .then(() => console.log('Admin role set!'));
"
```

---

## 📊 Part 7: Initialize Database

### 1. Create First Course (via API or Firestore Console)

**Via Firestore Console:**

1. Go to Firestore Database
2. Create collection: `courses`
3. Add document with auto-ID:

```json
{
  "title": "Introduction to Computer Science",
  "description": "Learn programming basics",
  "teacherId": "teacher-uid-here",
  "teacherName": "John Doe",
  "category": "Computer Science",
  "level": "Beginner",
  "enrollmentCount": 0,
  "isPublished": true,
  "createdAt": "2025-12-27T00:00:00Z",
  "updatedAt": "2025-12-27T00:00:00Z"
}
```

### 2. Create Sample Quiz

Use the API:

```bash
curl -X POST https://your-backend.vercel.app/api/quizzes/create \
  -H "Authorization: Bearer <teacher-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Python Basics Quiz",
    "courseId": "course-id-here",
    "duration": 30,
    "totalMarks": 100,
    "passingMarks": 60,
    "preventTabSwitch": true,
    "maxTabSwitches": 3,
    "requireFullscreen": true,
    "isPublished": true
  }'
```

---

## 🛠️ Troubleshooting

### Issue: "Firebase not initialized"

**Solution:** Check `FIREBASE_CREDENTIALS_JSON` is set correctly

### Issue: "CORS error"

**Solution:** 
1. Check `CORS_ORIGIN` in backend environment
2. Ensure frontend URL matches exactly (no trailing slash)

### Issue: "Module not found" in frontend

**Solution:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Issue: API returns 500 error

**Solution:**
1. Check Vercel function logs
2. Verify Firebase credentials
3. Check Firestore security rules

### Issue: Can't login

**Solution:**
1. Check Firebase Authentication is enabled
2. Verify email/password provider is enabled
3. Check browser console for errors

---

## 📚 Next Steps

1. ✅ **Create Admin Account**
   - Register via frontend
   - Set admin role via Firebase Console

2. ✅ **Create Teacher Account**
   - Use admin to set teacher role
   - API: `POST /api/auth/set-role`

3. ✅ **Create Sample Content**
   - Add courses
   - Create quizzes
   - Upload materials

4. ✅ **Test Features**
   - Student enrollment
   - Quiz taking with cheating prevention
   - Results viewing

5. ✅ **Monitor Performance**
   - Vercel Analytics
   - Firebase Usage
   - Error tracking

---

## 🔗 Useful Links

- [Firebase Console](https://console.firebase.google.com/)
- [Vercel Dashboard](https://vercel.com/dashboard)
- [Next.js Docs](https://nextjs.org/docs)
- [Go Firebase Admin SDK](https://firebase.google.com/docs/admin/setup)

---

## 📞 Support

For issues:
1. Check logs in Vercel Dashboard
2. Review Firebase Console for errors
3. See documentation files:
   - `README.md` - Project overview
   - `ARCHITECTURE.md` - System design
   - `QUIZ_SYSTEM_GUIDE.md` - Quiz features
   - `frontend/README.md` - Frontend guide

# Deployment Guide - LMS Platform

## Prerequisites Checklist

- [ ] Firebase Project created
- [ ] Vercel Account created
- [ ] GitHub repository (optional but recommended)
- [ ] Domain name (optional)

---

## Part 1: Firebase Setup

### Step 1: Create Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com)
2. Click "Add Project"
3. Name it (e.g., "lms-platform-prod")
4. Disable Google Analytics (optional)
5. Click "Create Project"

### Step 2: Enable Authentication

1. In Firebase Console, go to **Authentication** → **Sign-in method**
2. Enable **Email/Password**
3. Enable **Google** (add OAuth client ID)
4. Save changes

### Step 3: Create Firestore Database

1. Go to **Firestore Database** → **Create database**
2. Select **Production mode**
3. Choose region closest to your users (e.g., `us-central1`)
4. Click **Enable**

### Step 4: Enable Firebase Storage

1. Go to **Storage** → **Get Started**
2. Use default security rules
3. Choose same region as Firestore
4. Click **Done**

### Step 5: Get Service Account Credentials

1. Go to **Project Settings** → **Service Accounts**
2. Click **Generate new private key**
3. Save the JSON file securely (you'll need it for Vercel)

### Step 6: Configure Security Rules

#### Firestore Rules (`firestore.rules`)

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    
    // Helper functions
    function isAuthenticated() {
      return request.auth != null;
    }
    
    function hasRole(role) {
      return isAuthenticated() && request.auth.token.role == role;
    }
    
    function isOwner(uid) {
      return isAuthenticated() && request.auth.uid == uid;
    }
    
    // Users collection
    match /users/{userId} {
      allow read: if isOwner(userId) || hasRole('admin');
      allow update: if isOwner(userId) && 
                     request.resource.data.role == resource.data.role; // Can't change own role
      allow delete: if hasRole('admin');
    }
    
    // Courses collection
    match /courses/{courseId} {
      allow read: if resource.data.isPublished == true || 
                     hasRole('teacher') || 
                     hasRole('admin');
      allow create: if hasRole('teacher') || hasRole('admin');
      allow update, delete: if (hasRole('teacher') && resource.data.teacherId == request.auth.uid) || 
                              hasRole('admin');
    }
    
    // Enrollments collection
    match /enrollments/{enrollmentId} {
      allow read: if isOwner(resource.data.studentId) || 
                     hasRole('teacher') || 
                     hasRole('admin');
      allow create: if isOwner(request.resource.data.studentId) && 
                       hasRole('student');
      allow update: if isOwner(resource.data.studentId); // Student can update progress
      allow delete: if hasRole('admin');
    }
    
    // Quizzes collection
    match /quizzes/{quizId} {
      allow read: if resource.data.isPublished == true || 
                     resource.data.teacherId == request.auth.uid || 
                     hasRole('admin');
      allow write: if (hasRole('teacher') && request.resource.data.teacherId == request.auth.uid) || 
                      hasRole('admin');
    }
    
    // Questions collection
    match /questions/{questionId} {
      allow read: if hasRole('teacher') || hasRole('admin');
      allow write: if hasRole('teacher') || hasRole('admin');
    }
    
    // Quiz Submissions
    match /quiz_submissions/{submissionId} {
      allow read: if isOwner(resource.data.studentId) || 
                     hasRole('teacher') || 
                     hasRole('admin');
      allow create: if isOwner(request.resource.data.studentId);
      allow update: if hasRole('teacher') || hasRole('admin'); // For evaluation
    }
    
    // Exam Submissions
    match /exam_submissions/{submissionId} {
      allow read: if isOwner(resource.data.studentId) || 
                     hasRole('teacher') || 
                     hasRole('admin');
      allow create: if isOwner(request.resource.data.studentId);
      allow update: if hasRole('teacher') || hasRole('admin');
    }
    
    // Assignment Submissions
    match /assignment_submissions/{submissionId} {
      allow read: if isOwner(resource.data.studentId) || 
                     hasRole('teacher') || 
                     hasRole('admin');
      allow create: if isOwner(request.resource.data.studentId);
      allow update: if hasRole('teacher') || hasRole('admin');
    }
    
    // Analytics (read-only for students, write for teachers/admin)
    match /analytics/{analyticsId} {
      allow read: if isAuthenticated();
      allow write: if hasRole('teacher') || hasRole('admin');
    }
    
    // Notifications
    match /notifications/{notificationId} {
      allow read: if isOwner(resource.data.userId);
      allow write: if hasRole('teacher') || hasRole('admin');
    }
  }
}
```

**Deploy Rules**:
```bash
firebase deploy --only firestore:rules
```

#### Storage Rules (`storage.rules`)

```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    
    // Course materials (teacher upload)
    match /courses/{courseId}/{fileName} {
      allow read: if request.auth != null;
      allow write: if request.auth.token.role in ['teacher', 'admin'];
    }
    
    // Assignment submissions (student upload)
    match /assignments/{assignmentId}/{studentId}/{fileName} {
      allow read: if request.auth != null && 
                     (request.auth.uid == studentId || 
                      request.auth.token.role in ['teacher', 'admin']);
      allow write: if request.auth.uid == studentId;
    }
    
    // User profile photos
    match /users/{userId}/profile.jpg {
      allow read: if request.auth != null;
      allow write: if request.auth.uid == userId;
    }
  }
}
```

**Deploy Rules**:
```bash
firebase deploy --only storage:rules
```

---

## Part 2: Backend Deployment (Vercel)

### Step 1: Prepare Repository

```bash
# Initialize git if not already
git init
git add .
git commit -m "Initial LMS platform setup"

# Push to GitHub (recommended)
git remote add origin https://github.com/your-username/lms-platform.git
git push -u origin main
```

### Step 2: Install Vercel CLI

```bash
npm install -g vercel
```

### Step 3: Login to Vercel

```bash
vercel login
```

### Step 4: Link Project

```bash
# Run from project root
vercel link
```

Select:
- Scope: Your Vercel team/personal account
- Link to existing project? No
- Project name: `lms-platform-backend`
- Directory: `./` (root)

### Step 5: Configure Environment Variables

#### Option A: Via Vercel CLI

```bash
vercel env add FIREBASE_PROJECT_ID
# Enter value: your-project-id

vercel env add FIREBASE_CLIENT_EMAIL
# Enter value: firebase-adminsdk-xxxxx@your-project.iam.gserviceaccount.com

vercel env add FIREBASE_PRIVATE_KEY
# Enter value: -----BEGIN PRIVATE KEY-----\nYour key here\n-----END PRIVATE KEY-----

vercel env add FIREBASE_STORAGE_BUCKET
# Enter value: your-project.appspot.com
```

#### Option B: Via Vercel Dashboard

1. Go to [Vercel Dashboard](https://vercel.com/dashboard)
2. Select your project
3. Go to **Settings** → **Environment Variables**
4. Add each variable (for Production, Preview, Development)

**Important**: For `FIREBASE_PRIVATE_KEY`, replace newlines with `\n`:
```bash
# Example
-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQ...\n-----END PRIVATE KEY-----
```

### Step 6: Deploy

```bash
# Deploy to production
vercel --prod
```

Expected output:
```
✅ Deployment complete!
   Production: https://lms-platform-backend.vercel.app
```

### Step 7: Test API

```bash
curl https://your-deployment.vercel.app/api/auth/profile
# Should return 401 Unauthorized (expected without token)
```

---

## Part 3: Frontend Deployment

### Step 1: Create Frontend Directory

```bash
# In project root
npx create-next-app@latest frontend --typescript --tailwind --app
cd frontend
```

### Step 2: Install Dependencies

```bash
npm install firebase
npm install swr  # For data fetching
```

### Step 3: Configure Environment Variables

Create `frontend/.env.local`:

```env
# Firebase Client Config
NEXT_PUBLIC_FIREBASE_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXX
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abcdef123456

# Backend API URL
NEXT_PUBLIC_API_URL=https://lms-platform-backend.vercel.app
```

**Get Firebase Config**:
1. Firebase Console → Project Settings → General
2. Scroll to "Your apps"
3. Click "Add app" → Web
4. Copy the config values

### Step 4: Deploy Frontend

```bash
# From frontend directory
vercel --prod
```

Vercel will auto-detect Next.js and configure build settings.

---

## Part 4: Post-Deployment Configuration

### Step 1: Create Admin User

Since you need an admin to create other users, manually create the first admin:

1. **Register via Firebase Console**:
   - Go to Firebase Console → Authentication → Users
   - Click "Add user"
   - Email: `admin@yourdomain.com`
   - Password: (set strong password)
   - Click "Add user"

2. **Set Admin Role**:
   - Note the UID of the created user
   - Use Firebase CLI:

   ```bash
   # Install Firebase CLI
   npm install -g firebase-tools
   firebase login
   
   # Set custom claims
   firebase functions:shell
   
   # In the shell:
   admin.auth().setCustomUserClaims('USER_UID_HERE', {role: 'admin'})
   ```

   **OR** use a one-time script:

   ```javascript
   // set-admin.js
   const admin = require('firebase-admin');
   const serviceAccount = require('./serviceAccountKey.json');

   admin.initializeApp({
     credential: admin.credential.cert(serviceAccount)
   });

   const uid = 'USER_UID_FROM_CONSOLE';
   
   admin.auth().setCustomUserClaims(uid, { role: 'admin' })
     .then(() => {
       console.log('✅ Admin role set successfully');
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

3. **Update Firestore**:
   - Go to Firestore Console
   - Find the user document in `users` collection
   - Set `role: "admin"`

### Step 2: Configure CORS

Update `utils/response.go`:

```go
func EnableCORS(w http.ResponseWriter, r *http.Request) {
	allowedOrigins := []string{
		"http://localhost:3000",
		"https://your-frontend.vercel.app",
	}
	
	origin := r.Header.Get("Origin")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}
```

Redeploy backend:
```bash
vercel --prod
```

### Step 3: Test End-to-End

1. **Frontend**: Go to `https://your-frontend.vercel.app`
2. **Login**: Use admin credentials
3. **Create Teacher**: Admin → Users → Add Teacher
4. **Create Course**: Login as teacher → Courses → Create
5. **Enroll Student**: Create student account → Enroll in course

---

## Part 5: Custom Domain (Optional)

### For Backend:

1. Vercel Dashboard → Project → Settings → Domains
2. Add `api.yourdomain.com`
3. Add DNS records as instructed by Vercel
4. Update frontend `.env.local`:
   ```env
   NEXT_PUBLIC_API_URL=https://api.yourdomain.com
   ```

### For Frontend:

1. Vercel Dashboard → Frontend Project → Settings → Domains
2. Add `yourdomain.com` or `lms.yourdomain.com`
3. Configure DNS

---

## Part 6: Monitoring & Maintenance

### Vercel Logs

```bash
# View real-time logs
vercel logs https://your-deployment.vercel.app --follow
```

### Firebase Console

- **Authentication**: Monitor sign-ups, active users
- **Firestore**: Check document counts, query performance
- **Storage**: Monitor bandwidth usage

### Alerts

Set up in Vercel Dashboard:
- Function errors > 5%
- Execution time > 8s

---

## Part 7: CI/CD Pipeline (Optional)

### GitHub Actions Workflow

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Vercel

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy Backend to Vercel
        uses: amondnet/vercel-action@v20
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          vercel-args: '--prod'
```

**Setup**:
1. Generate Vercel token: Settings → Tokens
2. Add to GitHub Secrets: Settings → Secrets → Actions
3. Get Org ID and Project ID from `.vercel/project.json`

---

## Troubleshooting

### Issue: "Firebase Admin SDK not initialized"

**Solution**: Check environment variables are set correctly in Vercel.

```bash
vercel env ls
```

### Issue: CORS errors

**Solution**: Verify allowed origins in `utils/response.go` match frontend URL.

### Issue: "Token expired"

**Solution**: Frontend should auto-refresh tokens. Check Firebase Auth config.

### Issue: Cold starts taking too long

**Solution**: 
- Reduce function complexity
- Use Vercel Pro for better cold start performance
- Implement warming functions (cron jobs)

---

## Backup Strategy

### Firestore Backup

```bash
# Install gcloud CLI
gcloud auth login

# Export Firestore
gcloud firestore export gs://your-bucket/backups/$(date +%Y%m%d)
```

**Automate**: Set up Cloud Scheduler to run daily exports.

### Code Backup

- GitHub repository (automatic with each push)
- Vercel maintains deployment history

---

## Cost Monitoring

### Vercel

- Dashboard → Usage
- Set billing alerts

### Firebase

- Console → Usage and Billing
- Set budget alerts in Google Cloud Console

---

**Deployment Complete! 🎉**

Your LMS platform is now live and ready for production use.

# Deployment Checklist

## Pre-Deployment

- [ ] Firebase project created
- [ ] Firebase Authentication enabled (Email/Password)
- [ ] Firestore database created (Production mode)
- [ ] Firebase Storage enabled
- [ ] Service account key downloaded
- [ ] All environment variables configured

## Backend Deployment (Vercel)

### 1. Install Vercel CLI
```bash
npm install -g vercel
```

### 2. Login to Vercel
```bash
vercel login
```

### 3. Deploy from Root Directory
```bash
cd D:\LMS
vercel
```

### 4. Add Environment Variables in Vercel Dashboard

Go to: **Vercel Dashboard → Project → Settings → Environment Variables**

Add:
- **Name**: `FIREBASE_CREDENTIALS_JSON`
- **Value**: Paste entire content from `serviceAccountKey.json`
- **Environment**: Production, Preview, Development

### 5. Deploy to Production
```bash
vercel --prod
```

**Save the backend URL**: `https://your-backend.vercel.app`

## Frontend Deployment (Vercel)

### 1. Deploy Frontend
```bash
cd frontend
vercel
```

### 2. Add Frontend Environment Variables

Go to: **Vercel Dashboard → Frontend Project → Settings → Environment Variables**

Add all from `.env.local`:
- `NEXT_PUBLIC_FIREBASE_API_KEY`
- `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN`
- `NEXT_PUBLIC_FIREBASE_PROJECT_ID`
- `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET`
- `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID`
- `NEXT_PUBLIC_FIREBASE_APP_ID`
- `NEXT_PUBLIC_API_URL` = `https://your-backend.vercel.app/api`

### 3. Deploy to Production
```bash
vercel --prod
```

**Save the frontend URL**: `https://your-frontend.vercel.app`

## Firebase Security Rules Deployment

### 1. Login to Firebase
```bash
firebase login
```

### 2. Initialize Firebase (if not done)
```bash
cd D:\LMS
firebase init
```

Select:
- Firestore
- Storage

Choose existing project: `flutterdemo-6d4ac`

### 3. Deploy Firestore Rules
```bash
firebase deploy --only firestore:rules
```

### 4. Deploy Storage Rules
```bash
firebase deploy --only storage:rules
```

## Post-Deployment

### 1. Update CORS Settings

Update backend environment variable:
- `CORS_ORIGIN` = `https://your-frontend.vercel.app`

Redeploy backend:
```bash
vercel --prod
```

### 2. Test the Application

- [ ] Visit frontend URL
- [ ] Test user registration
- [ ] Test login
- [ ] Create a course (as teacher)
- [ ] Create a quiz
- [ ] Enroll as student
- [ ] Take quiz
- [ ] Verify results

### 3. Create Admin User

Option 1: Via Firebase Console
1. Go to Firebase Console → Authentication → Users
2. Find your user
3. Set custom claims manually

Option 2: Via API
```bash
curl -X POST https://your-backend.vercel.app/api/auth/set-role \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "uid": "user-uid-here",
    "role": "admin"
  }'
```

## Monitoring

- [ ] Check Vercel Analytics
- [ ] Monitor Firebase Usage
- [ ] Set up error tracking
- [ ] Configure alerts

## Custom Domain (Optional)

### Backend
1. Vercel Dashboard → Backend Project → Settings → Domains
2. Add domain: `api.yourdomain.com`
3. Update DNS records

### Frontend
1. Vercel Dashboard → Frontend Project → Settings → Domains
2. Add domain: `lms.yourdomain.com`
3. Update DNS records

### Update Environment Variables
- Update `NEXT_PUBLIC_API_URL` to `https://api.yourdomain.com/api`
- Update `CORS_ORIGIN` to `https://lms.yourdomain.com`
- Redeploy both projects

## Troubleshooting

### Backend 500 Error
- Check Vercel function logs
- Verify `FIREBASE_CREDENTIALS_JSON` is set correctly
- Ensure Firebase project ID matches

### Frontend Can't Connect to Backend
- Verify `NEXT_PUBLIC_API_URL` points to backend URL
- Check CORS settings
- Verify backend is deployed

### Authentication Errors
- Check Firebase config in frontend
- Verify Email/Password auth is enabled
- Check Firebase Authentication tab for errors

## Success Criteria

- [ ] Backend responds at `/api/health`
- [ ] Frontend loads successfully
- [ ] Can register new user
- [ ] Can login
- [ ] Firebase shows user in Authentication
- [ ] Can create and view courses
- [ ] Can create and take quizzes
- [ ] Security rules are active

## Rollback Plan

If issues occur:
```bash
# Rollback backend
vercel rollback

# Rollback frontend
cd frontend
vercel rollback
```

# Vercel Deployment Guide

## Prerequisites

✅ Git repository initialized
✅ Code committed to Git
✅ Vercel account created
✅ Firebase project configured

## Step 1: Install Vercel CLI

```bash
npm install -g vercel
```

## Step 2: Login to Vercel

```bash
vercel login
```

## Step 3: Deploy Backend

```bash
cd D:\LMS
vercel
```

Follow prompts:
- Set up and deploy? **Y**
- Which scope? Select your account
- Link to existing project? **N**
- Project name? **lms-backend** (or your choice)
- Which directory? **.** (current directory)
- Override settings? **N**

**Copy the deployment URL** (e.g., `https://lms-backend.vercel.app`)

## Step 4: Add Backend Environment Variables

1. Go to: https://vercel.com/dashboard
2. Select **lms-backend** project
3. Click **Settings** → **Environment Variables**
4. Add:

```
Name: FIREBASE_CREDENTIALS_JSON
Value: (paste entire content from your Firebase service account JSON)
Environment: Production, Preview, Development
```

To get Firebase credentials:
1. Firebase Console → Project Settings → Service Accounts
2. Click "Generate New Private Key"
3. Save as `serviceAccountKey.json` (don't commit!)
4. Copy entire JSON content

## Step 5: Redeploy Backend

```bash
vercel --prod
```

**Backend is now live!** Test: `https://your-backend.vercel.app/api/health`

## Step 6: Deploy Frontend

```bash
cd frontend
vercel
```

Follow prompts:
- Set up and deploy? **Y**
- Link to existing project? **N**
- Project name? **lms-frontend**
- Which directory? **.** (current directory)
- Override settings? **N**

## Step 7: Add Frontend Environment Variables

Go to: **Vercel Dashboard → lms-frontend → Settings → Environment Variables**

Add all from `.env.local`:

```
NEXT_PUBLIC_FIREBASE_API_KEY=AIzaSyDV2vK5mXHKnFjPBF4UwoctJMF_uzTEpfg
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=flutterdemo-6d4ac.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=flutterdemo-6d4ac
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=flutterdemo-6d4ac.firebasestorage.app
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=764379386105
NEXT_PUBLIC_FIREBASE_APP_ID=1:764379386105:web:c62baca4253f4d351cd228
NEXT_PUBLIC_API_URL=https://your-backend.vercel.app/api
```

⚠️ **Important**: Replace `NEXT_PUBLIC_API_URL` with your actual backend URL from Step 5!

## Step 8: Redeploy Frontend

```bash
vercel --prod
```

**Frontend is now live!** Visit: `https://your-frontend.vercel.app`

## Step 9: Deploy Firebase Security Rules

```bash
# Install Firebase CLI
npm install -g firebase-tools

# Login
firebase login

# Initialize (if not done)
firebase init

# Select:
# - Firestore
# - Storage
# - Choose existing project: flutterdemo-6d4ac

# Deploy rules
firebase deploy --only firestore:rules
firebase deploy --only storage:rules
```

## Step 10: Test the Deployment

1. Visit your frontend URL
2. Click "Register" → Create account
3. Login with new account
4. Should see dashboard (as student by default)

## Optional: Link to GitHub (Auto-Deploy)

### 1. Push to GitHub

```bash
cd D:\LMS
git remote add origin https://github.com/yourusername/lms-platform.git
git branch -M main
git push -u origin main
```

### 2. Connect Vercel to GitHub

**Backend:**
1. Vercel Dashboard → lms-backend → Settings → Git
2. Click "Connect Git Repository"
3. Select your GitHub repo
4. Root Directory: `/`

**Frontend:**
1. Vercel Dashboard → lms-frontend → Settings → Git
2. Connect to same repo
3. Root Directory: `/frontend`

Now every push to `main` auto-deploys!

## GitHub Actions CI/CD (Advanced)

Already configured in `.github/workflows/deploy.yml`

Add secrets in GitHub:
- `VERCEL_TOKEN`: Get from https://vercel.com/account/tokens
- `VERCEL_ORG_ID`: From Vercel project settings
- `VERCEL_PROJECT_ID`: From Vercel project settings
- All Firebase env vars

## Vercel Project Settings

### Backend (`vercel.json`)
```json
{
  "functions": {
    "api/**/*.go": {
      "runtime": "go1.x",
      "memory": 1024,
      "maxDuration": 10
    }
  }
}
```

### Frontend (`vercel.json` in frontend/)
Will be auto-detected as Next.js project

## Custom Domain (Optional)

### Add Domain to Backend
1. Vercel → lms-backend → Settings → Domains
2. Add: `api.yourdomain.com`
3. Update DNS:
   - Type: CNAME
   - Name: api
   - Value: cname.vercel-dns.com

### Add Domain to Frontend
1. Vercel → lms-frontend → Settings → Domains
2. Add: `lms.yourdomain.com`
3. Update DNS

### Update Environment Variables
- `NEXT_PUBLIC_API_URL` → `https://api.yourdomain.com/api`
- Redeploy frontend

## Monitoring & Logs

### View Logs
```bash
# Backend logs
vercel logs lms-backend --prod

# Frontend logs
cd frontend
vercel logs lms-frontend --prod
```

### Vercel Analytics
- Enable in: Project → Analytics
- Free tier: 10k events/month

## Troubleshooting

### Backend 500 Error
```bash
# Check logs
vercel logs lms-backend

# Verify env vars
vercel env ls
```

### Frontend Build Fails
```bash
# Test locally first
cd frontend
npm run build

# Check env vars in Vercel dashboard
```

### CORS Errors
Ensure backend allows frontend domain:
```go
// In API handlers, add:
w.Header().Set("Access-Control-Allow-Origin", "https://your-frontend.vercel.app")
```

### Firebase Auth Not Working
- Verify all Firebase env vars in frontend
- Check Firebase Console → Authentication → Settings → Authorized domains
- Add Vercel frontend domain: `your-frontend.vercel.app`

## Security Checklist

- [ ] `.env.local` not committed (in `.gitignore`)
- [ ] `serviceAccountKey.json` not committed
- [ ] Firestore security rules deployed
- [ ] Storage security rules deployed
- [ ] Firebase authorized domains updated
- [ ] Environment variables set in Vercel dashboard
- [ ] CORS properly configured

## Success!

✅ Backend deployed: `https://your-backend.vercel.app`
✅ Frontend deployed: `https://your-frontend.vercel.app`
✅ Firebase rules active
✅ Auto-deploy on Git push (if configured)

Your LMS Platform is now live! 🚀

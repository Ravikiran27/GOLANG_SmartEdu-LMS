# LMS Frontend Setup Guide

## Prerequisites

- Node.js 18+ installed
- npm or yarn package manager
- Firebase project configured

---

## Installation Steps

### 1. Install Dependencies

```bash
cd frontend
npm install
```

This will install:
- **Next.js 14** - React framework with App Router
- **React 18** - UI library
- **Firebase 10** - Authentication, Firestore, Storage
- **TypeScript 5** - Type safety
- **Tailwind CSS 3** - Utility-first CSS
- **SWR 2** - Data fetching and caching

---

### 2. Configure Environment Variables

Create `.env.local` file in the `frontend` directory:

```bash
cp .env.example .env.local
```

Edit `.env.local` with your Firebase credentials:

```env
NEXT_PUBLIC_FIREBASE_API_KEY=AIza...
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789:web:abc123
NEXT_PUBLIC_API_URL=http://localhost:3000/api
```

**Get Firebase credentials**:
1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Click ⚙️ Settings → Project settings
4. Scroll to "Your apps" → Web app
5. Copy the config values

---

### 3. Run Development Server

```bash
npm run dev
```

The app will be available at: **http://localhost:3000**

---

## Project Structure

```
frontend/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # Root layout with AuthProvider
│   ├── page.tsx           # Landing page
│   ├── globals.css        # Global styles
│   └── auth/              # Authentication pages
│       └── login/
│           └── page.tsx   # Login page
│
├── components/            # React components
│   └── ProtectedRoute.tsx # Route guard
│
├── lib/                   # Utilities
│   ├── firebase.ts        # Firebase client initialization
│   ├── auth-context.tsx   # Auth state management
│   └── api.ts             # API client
│
├── public/                # Static assets
├── .env.example           # Environment template
├── .env.local            # Your environment variables (gitignored)
├── next.config.js        # Next.js configuration
├── tailwind.config.js    # Tailwind CSS configuration
├── tsconfig.json         # TypeScript configuration
└── package.json          # Dependencies
```

---

## Available Scripts

```bash
# Development server
npm run dev

# Production build
npm run build

# Start production server
npm start

# Run linter
npm run lint
```

---

## Features Implemented

### ✅ Authentication System
- Firebase Authentication integration
- Email/Password login
- Protected routes with role-based access
- Auth context for global state management

### ✅ API Client
- Automatic JWT token injection
- TypeScript type safety
- Error handling
- Request/response interceptors

### ✅ UI Components
- Responsive layouts
- Tailwind CSS styling
- Loading states
- Protected route wrapper

---

## Adding New Pages

### 1. Student Dashboard

Create `app/student/dashboard/page.tsx`:

```tsx
'use client';

import { useAuth } from '@/lib/auth-context';
import ProtectedRoute from '@/components/ProtectedRoute';

export default function StudentDashboard() {
  const { user } = useAuth();

  return (
    <ProtectedRoute allowedRoles={['student']}>
      <div className="p-8">
        <h1 className="text-3xl font-bold">Welcome, {user?.displayName}</h1>
        {/* Your dashboard content */}
      </div>
    </ProtectedRoute>
  );
}
```

### 2. Teacher Dashboard

Create `app/teacher/dashboard/page.tsx`:

```tsx
'use client';

import ProtectedRoute from '@/components/ProtectedRoute';

export default function TeacherDashboard() {
  return (
    <ProtectedRoute allowedRoles={['teacher']}>
      <div className="p-8">
        <h1 className="text-3xl font-bold">Teacher Dashboard</h1>
        {/* Your dashboard content */}
      </div>
    </ProtectedRoute>
  );
}
```

---

## Using the API Client

### Example: Fetch Courses

```tsx
'use client';

import { useEffect, useState } from 'react';
import { apiRequest } from '@/lib/api';

export default function CourseList() {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchCourses() {
      try {
        const response = await apiRequest('/courses/list');
        if (response.success) {
          setCourses(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch courses:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchCourses();
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      {courses.map((course) => (
        <div key={course.id}>{course.title}</div>
      ))}
    </div>
  );
}
```

### Example: Create Quiz

```tsx
import { apiRequest } from '@/lib/api';

async function createQuiz(data) {
  const response = await apiRequest('/quizzes/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });

  if (response.success) {
    console.log('Quiz created:', response.data);
  } else {
    console.error('Error:', response.error);
  }
}
```

---

## Authentication Flow

### Login Process

1. User enters email/password
2. Firebase authenticates user
3. Custom claims fetched (role: admin/teacher/student)
4. Auth context updates global state
5. JWT token stored for API requests
6. User redirected to role-based dashboard

### Protected Routes

```tsx
<ProtectedRoute allowedRoles={['teacher', 'admin']}>
  <YourComponent />
</ProtectedRoute>
```

**Features**:
- Redirects unauthenticated users to `/auth/login`
- Redirects unauthorized users to `/unauthorized`
- Shows loading state during auth check

---

## Styling with Tailwind

### Example Component

```tsx
export default function Card({ title, children }) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">{title}</h2>
      <div className="text-gray-600">{children}</div>
    </div>
  );
}
```

### Common Tailwind Classes

```css
/* Layout */
flex, grid, container, mx-auto

/* Spacing */
p-4, m-4, px-8, py-4, space-x-4

/* Typography */
text-xl, font-bold, text-gray-700

/* Colors */
bg-blue-600, text-white, hover:bg-blue-700

/* Borders */
border, border-gray-300, rounded-lg

/* Shadows */
shadow-md, shadow-lg, hover:shadow-xl
```

---

## Troubleshooting

### Issue: "Module not found"

**Solution**: Install dependencies
```bash
npm install
```

### Issue: Firebase errors

**Solution**: Check `.env.local` has correct Firebase config

### Issue: API calls failing

**Solution**: 
1. Verify backend is running
2. Check `NEXT_PUBLIC_API_URL` in `.env.local`
3. Ensure Firebase Auth token is valid

### Issue: TypeScript errors

**Solution**: 
```bash
# Delete cache and rebuild
rm -rf .next
npm run build
```

---

## Deployment to Vercel

### 1. Install Vercel CLI

```bash
npm i -g vercel
```

### 2. Deploy

```bash
cd frontend
vercel
```

### 3. Configure Environment Variables

In Vercel dashboard:
1. Go to Project Settings → Environment Variables
2. Add all `NEXT_PUBLIC_*` variables
3. Redeploy

---

## Next Steps

1. **Build Quiz Interface**
   - Quiz taking component with timer
   - Tab switching detection
   - Fullscreen enforcement
   - Answer submission

2. **Build Dashboards**
   - Student: View courses, quizzes, grades
   - Teacher: Create courses, manage quizzes
   - Admin: User management, reports

3. **Add Real-time Features**
   - Firestore listeners for live updates
   - Notifications system
   - Live quiz results

4. **Optimize Performance**
   - Code splitting
   - Image optimization
   - Caching strategies with SWR

---

## Support

For issues or questions:
- Check documentation in `/docs`
- Review `ARCHITECTURE.md` for system design
- See `QUIZ_SYSTEM_GUIDE.md` for quiz features

# LMS Platform

Full-stack Learning Management System with Quiz & Exam Platform built with Golang, Next.js, and Firebase.

## Features

- 🔐 Role-based authentication (Admin, Teacher, Student)
- 📚 Course management with enrollment
- 📝 Quiz system with cheating prevention
- 📊 Exam system with auto/manual grading
- 📋 Assignment management
- 📈 Analytics and performance tracking
- 👨‍💼 Admin dashboard

## Tech Stack

- **Backend**: Golang serverless functions (Vercel)
- **Frontend**: Next.js 14 with TypeScript
- **Database**: Firebase Firestore
- **Authentication**: Firebase Auth
- **Storage**: Firebase Storage
- **Deployment**: Vercel

## Quick Start

### Prerequisites

- Node.js 18+
- Go 1.21+
- Firebase project
- Vercel account

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd LMS
```

2. Setup backend:
```bash
# Add Firebase credentials to .env
echo 'FIREBASE_CREDENTIALS_JSON={"type":"service_account",...}' > .env
```

3. Setup frontend:
```bash
cd frontend
npm install
cp .env.example .env.local
# Edit .env.local with your Firebase config
```

4. Run locally:
```bash
# Terminal 1 - Backend
vercel dev

# Terminal 2 - Frontend
cd frontend
npm run dev
```

Visit http://localhost:3000

## Deployment

See [INSTALLATION_GUIDE.md](./INSTALLATION_GUIDE.md) for complete deployment instructions.

### Quick Deploy

```bash
# Deploy backend
vercel --prod

# Deploy frontend
cd frontend
vercel --prod
```

## Documentation

- [Installation Guide](./INSTALLATION_GUIDE.md) - Complete setup instructions
- [Architecture](./ARCHITECTURE.md) - System design and patterns
- [Quiz System Guide](./QUIZ_SYSTEM_GUIDE.md) - Quiz features and cheating prevention
- [API Documentation](./IMPLEMENTATION_GUIDE.md) - API endpoints reference

## Project Structure

```
LMS/
├── api/                    # Golang serverless functions
│   ├── auth/              # Authentication APIs
│   ├── courses/           # Course management APIs
│   ├── quizzes/           # Quiz system APIs
│   └── ...
├── frontend/              # Next.js frontend
│   ├── app/              # App router pages
│   ├── components/       # React components
│   └── lib/              # Utilities
├── models/               # Go data models
├── utils/                # Go utilities
├── firestore.rules       # Firestore security rules
└── storage.rules         # Storage security rules
```

## License

MIT

## Support

For issues and questions, please open a GitHub issue.

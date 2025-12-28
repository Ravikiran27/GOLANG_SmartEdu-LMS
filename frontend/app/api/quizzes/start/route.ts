// Quizzes API: Start quiz attempt
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  createdResponse, 
  verifyAuthToken,
  generateId 
} from '@/lib/api-utils';
import { Quiz, QuizSubmission, Question, User } from '@/lib/types';

// POST: Start a quiz attempt
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const body = await request.json();
    const { quizId } = body;

    if (!quizId) {
      return errorResponse('Quiz ID is required', 400);
    }

    // Get quiz
    const quizDoc = await db.collection('quizzes').doc(quizId).get();
    if (!quizDoc.exists) {
      return errorResponse('Quiz not found', 404);
    }

    const quiz = { id: quizDoc.id, ...quizDoc.data() } as Quiz;

    if (!quiz.isPublished || quiz.isDeleted) {
      return errorResponse('Quiz is not available', 400);
    }

    // Check deadline
    if (quiz.deadline && new Date(quiz.deadline) < new Date()) {
      return errorResponse('Quiz deadline has passed', 400);
    }

    // Check for existing in-progress submission
    const existingSubmission = await db.collection('submissions')
      .where('quizId', '==', quizId)
      .where('studentId', '==', authUser.uid)
      .where('status', '==', 'in_progress')
      .limit(1)
      .get();

    if (!existingSubmission.empty) {
      const existing = { id: existingSubmission.docs[0].id, ...existingSubmission.docs[0].data() } as QuizSubmission;
      return successResponse(existing, 'Resuming existing attempt');
    }

    // Check attempt count
    const attemptCount = await db.collection('submissions')
      .where('quizId', '==', quizId)
      .where('studentId', '==', authUser.uid)
      .get();

    if (quiz.maxAttempts > 0 && attemptCount.size >= quiz.maxAttempts) {
      return errorResponse('Maximum attempts reached', 400);
    }

    // Get user info
    const userDoc = await db.collection('users').doc(authUser.uid).get();
    const user = userDoc.data() as User;

    // Get questions
    const questionsSnapshot = await db.collection('quizzes').doc(quizId)
      .collection('questions')
      .orderBy('order', 'asc')
      .get();

    let questions: Question[] = [];
    questionsSnapshot.forEach((doc) => {
      const q = { id: doc.id, ...doc.data() } as Question;
      // Hide correct answers
      delete q.correctAnswer;
      if (q.options) {
        q.options = q.options.map(opt => ({ ...opt, isCorrect: false }));
      }
      questions.push(q);
    });

    // Shuffle if required
    if (quiz.shuffleQuestions || quiz.randomizeQuestionOrder) {
      questions = questions.sort(() => Math.random() - 0.5);
    }

    // Create submission
    const now = new Date();
    const submissionId = generateId();
    const submission: QuizSubmission = {
      id: submissionId,
      submissionId,
      quizId,
      studentId: authUser.uid,
      studentName: user.displayName,
      studentEmail: user.email,
      courseId: quiz.courseId,
      attemptNumber: attemptCount.size + 1,
      answers: [],
      questions,
      startedAt: now,
      timeLimit: quiz.duration,
      totalMarks: quiz.totalMarks,
      marksObtained: 0,
      percentage: 0,
      passed: false,
      status: 'in_progress',
      tabSwitchCount: 0,
      fullscreenExits: 0,
      suspiciousActivity: [],
      createdAt: now,
      updatedAt: now,
    };

    await db.collection('submissions').doc(submissionId).set(submission);

    return createdResponse(submission, 'Quiz started successfully');

  } catch (error: any) {
    console.error('Start quiz error:', error);
    return errorResponse('Failed to start quiz: ' + error.message, 500);
  }
}

export async function OPTIONS() {
  return new Response(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

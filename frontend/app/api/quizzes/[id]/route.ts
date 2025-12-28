// Quizzes API: Get single quiz
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken 
} from '@/lib/api-utils';
import { Quiz, Question } from '@/lib/types';

interface RouteParams {
  params: Promise<{ id: string }>;
}

// GET: Get single quiz with questions
export async function GET(request: NextRequest, { params }: RouteParams) {
  try {
    const { id } = await params;
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const quizDoc = await db.collection('quizzes').doc(id).get();
    if (!quizDoc.exists) {
      return errorResponse('Quiz not found', 404);
    }

    const quiz = { id: quizDoc.id, ...quizDoc.data() } as Quiz;

    if (quiz.isDeleted) {
      return errorResponse('Quiz not found', 404);
    }

    // Students can only see published quizzes
    if (authUser.role === 'student' && !quiz.isPublished) {
      return errorResponse('Quiz not available', 403);
    }

    // Teachers can only see their own unpublished quizzes
    if (authUser.role === 'teacher' && !quiz.isPublished && quiz.teacherId !== authUser.uid) {
      return errorResponse('Quiz not available', 403);
    }

    // Get questions
    const questionsSnapshot = await db.collection('quizzes').doc(id)
      .collection('questions')
      .orderBy('order', 'asc')
      .get();

    const questions: Question[] = [];
    questionsSnapshot.forEach((doc) => {
      const question = { id: doc.id, ...doc.data() } as Question;
      
      // Hide correct answers for students until submission
      if (authUser.role === 'student') {
        delete question.correctAnswer;
        if (question.options) {
          question.options = question.options.map(opt => ({
            ...opt,
            isCorrect: false, // Hide correct answer
          }));
        }
      }
      
      questions.push(question);
    });

    return successResponse({
      ...quiz,
      questions,
    });

  } catch (error: any) {
    console.error('Get quiz error:', error);
    return errorResponse('Failed to get quiz: ' + error.message, 500);
  }
}

export async function OPTIONS() {
  return new Response(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

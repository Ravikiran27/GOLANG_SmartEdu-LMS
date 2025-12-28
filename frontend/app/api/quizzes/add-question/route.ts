// Quizzes API: Add question to quiz
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  createdResponse, 
  verifyAuthToken, 
  checkRole,
  generateId 
} from '@/lib/api-utils';
import { CreateQuestionRequest, Question, Quiz } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

// POST: Add question to quiz
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Only teachers and admins can add questions
    if (!checkRole(authUser, ['teacher', 'admin'])) {
      return errorResponse('Forbidden: Only teachers and admins can add questions', 403);
    }

    const body: CreateQuestionRequest = await request.json();

    // Validate required fields
    if (!body.quizId) {
      return errorResponse('Quiz ID is required', 400);
    }
    if (!body.type) {
      return errorResponse('Question type is required', 400);
    }
    if (!body.questionText) {
      return errorResponse('Question text is required', 400);
    }
    if (!body.marks || body.marks <= 0) {
      return errorResponse('Marks must be greater than 0', 400);
    }

    // Verify quiz exists and user is the teacher or admin
    const quizDoc = await db.collection('quizzes').doc(body.quizId).get();
    if (!quizDoc.exists) {
      return errorResponse('Quiz not found', 404);
    }

    const quiz = quizDoc.data() as Quiz;
    if (authUser.role !== 'admin' && quiz.teacherId !== authUser.uid) {
      return errorResponse('Forbidden: You can only add questions to your own quizzes', 403);
    }

    // Get current question count for order
    const questionsSnapshot = await db.collection('quizzes').doc(body.quizId)
      .collection('questions')
      .get();
    const order = questionsSnapshot.size + 1;

    // Create question
    const now = new Date();
    const questionId = generateId();
    const question: Question = {
      id: questionId,
      questionId,
      quizId: body.quizId,
      type: body.type,
      text: body.questionText,
      questionText: body.questionText,
      imageUrl: body.imageUrl || '',
      marks: body.marks,
      points: body.marks,
      options: body.options || [],
      correctAnswer: body.correctAnswer || '',
      explanation: body.explanation || '',
      order,
      createdAt: now,
      updatedAt: now,
    };

    // Save question as subcollection
    await db.collection('quizzes').doc(body.quizId)
      .collection('questions').doc(questionId).set(question);

    // Update quiz question count
    await db.collection('quizzes').doc(body.quizId).update({
      questionCount: FieldValue.increment(1),
      updatedAt: FieldValue.serverTimestamp(),
    });

    return createdResponse(question, 'Question added successfully');

  } catch (error: any) {
    console.error('Add question error:', error);
    return errorResponse('Failed to add question: ' + error.message, 500);
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

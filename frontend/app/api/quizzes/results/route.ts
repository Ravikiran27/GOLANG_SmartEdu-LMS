// Quizzes API: Get quiz results
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken,
  getPaginationParams,
  checkRole 
} from '@/lib/api-utils';
import { QuizSubmission, Quiz } from '@/lib/types';

// GET: Get quiz results
export async function GET(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const url = new URL(request.url);
    const quizId = url.searchParams.get('quizId');
    const submissionId = url.searchParams.get('submissionId');
    const { page, pageSize } = getPaginationParams(request);

    // If specific submission requested
    if (submissionId) {
      const submissionDoc = await db.collection('submissions').doc(submissionId).get();
      if (!submissionDoc.exists) {
        return errorResponse('Submission not found', 404);
      }

      const submission = { id: submissionDoc.id, ...submissionDoc.data() } as QuizSubmission;

      // Check access
      if (authUser.role === 'student' && submission.studentId !== authUser.uid) {
        return errorResponse('Forbidden', 403);
      }

      // Get quiz for checking allowReview
      const quizDoc = await db.collection('quizzes').doc(submission.quizId).get();
      const quiz = quizDoc.data() as Quiz;

      return successResponse({
        ...submission,
        allowReview: quiz.allowReview,
      });
    }

    // List submissions
    let query = db.collection('submissions');

    if (authUser.role === 'student') {
      // Students see only their submissions
      query = query.where('studentId', '==', authUser.uid) as any;
    } else if (authUser.role === 'teacher') {
      // Teachers see submissions for their quizzes
      if (quizId) {
        // Verify teacher owns this quiz
        const quizDoc = await db.collection('quizzes').doc(quizId).get();
        if (quizDoc.exists) {
          const quiz = quizDoc.data() as Quiz;
          if (quiz.teacherId !== authUser.uid) {
            return errorResponse('Forbidden', 403);
          }
        }
        query = query.where('quizId', '==', quizId) as any;
      } else {
        // Get all quizzes by this teacher first
        const teacherQuizzes = await db.collection('quizzes')
          .where('teacherId', '==', authUser.uid)
          .get();
        const quizIds = teacherQuizzes.docs.map(d => d.id);
        
        if (quizIds.length === 0) {
          return successResponse({ submissions: [], page, pageSize });
        }
        
        query = query.where('quizId', 'in', quizIds.slice(0, 10)) as any; // Firestore limit
      }
    } else if (quizId) {
      // Admin with filter
      query = query.where('quizId', '==', quizId) as any;
    }

    const snapshot = await query
      .orderBy('submittedAt', 'desc')
      .limit(pageSize)
      .offset((page - 1) * pageSize)
      .get();

    const submissions: QuizSubmission[] = [];
    snapshot.forEach((doc) => {
      submissions.push({ id: doc.id, ...doc.data() } as QuizSubmission);
    });

    return successResponse({
      submissions,
      page,
      pageSize,
    });

  } catch (error: any) {
    console.error('Get results error:', error);
    return errorResponse('Failed to get results: ' + error.message, 500);
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

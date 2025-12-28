// Quizzes API: Resume quiz (Teacher only)
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken,
  checkRole 
} from '@/lib/api-utils';
import { QuizSubmission, Quiz } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

// POST: Resume a student's quiz (Teacher only)
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Only teachers and admins can resume quizzes
    if (!checkRole(authUser, ['teacher', 'admin'])) {
      return errorResponse('Forbidden: Only teachers and admins can resume quizzes', 403);
    }

    const body = await request.json();
    const { submissionId, reason } = body;

    if (!submissionId) {
      return errorResponse('submissionId is required', 400);
    }

    // Get submission
    const submissionDoc = await db.collection('submissions').doc(submissionId).get();
    if (!submissionDoc.exists) {
      return errorResponse('Submission not found', 404);
    }

    const submission = { id: submissionDoc.id, ...submissionDoc.data() } as QuizSubmission;

    // Verify quiz allows teacher resume
    const quizDoc = await db.collection('quizzes').doc(submission.quizId).get();
    const quiz = quizDoc.data() as Quiz;

    if (!quiz.allowTeacherResume && authUser.role !== 'admin') {
      return errorResponse('This quiz does not allow teacher resume', 400);
    }

    // Verify teacher owns this quiz
    if (authUser.role === 'teacher' && quiz.teacherId !== authUser.uid) {
      return errorResponse('Forbidden: You can only manage your own quizzes', 403);
    }

    // Update submission to allow resume
    const now = new Date();
    await db.collection('submissions').doc(submissionId).update({
      status: 'in_progress',
      resumedBy: authUser.uid,
      resumedAt: now,
      resumeReason: reason || 'Resumed by teacher',
      suspiciousActivity: FieldValue.arrayUnion(`Quiz resumed by teacher: ${reason || 'No reason provided'}`),
      updatedAt: FieldValue.serverTimestamp(),
    });

    // Get updated submission
    const updatedDoc = await db.collection('submissions').doc(submissionId).get();
    const updatedSubmission = { id: updatedDoc.id, ...updatedDoc.data() } as QuizSubmission;

    return successResponse(updatedSubmission, 'Quiz resumed successfully');

  } catch (error: any) {
    console.error('Resume quiz error:', error);
    return errorResponse('Failed to resume quiz: ' + error.message, 500);
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

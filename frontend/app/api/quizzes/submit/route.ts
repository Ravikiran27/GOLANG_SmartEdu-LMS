// Quizzes API: Submit quiz
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken 
} from '@/lib/api-utils';
import { SubmitQuizRequest, QuizSubmission, Quiz, Question } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

// POST: Submit quiz
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const body: SubmitQuizRequest = await request.json();

    if (!body.submissionId || !body.quizId) {
      return errorResponse('submissionId and quizId are required', 400);
    }

    // Get submission
    const submissionDoc = await db.collection('submissions').doc(body.submissionId).get();
    if (!submissionDoc.exists) {
      return errorResponse('Submission not found', 404);
    }

    const submission = { id: submissionDoc.id, ...submissionDoc.data() } as QuizSubmission;

    // Verify ownership
    if (submission.studentId !== authUser.uid) {
      return errorResponse('Forbidden: This is not your submission', 403);
    }

    if (submission.status !== 'in_progress') {
      return errorResponse('Quiz already submitted', 400);
    }

    // Get quiz for grading
    const quizDoc = await db.collection('quizzes').doc(body.quizId).get();
    const quiz = quizDoc.data() as Quiz;

    // Get questions with correct answers
    const questionsSnapshot = await db.collection('quizzes').doc(body.quizId)
      .collection('questions')
      .get();

    const questionsMap = new Map<string, Question>();
    questionsSnapshot.forEach((doc) => {
      const q = { id: doc.id, ...doc.data() } as Question;
      questionsMap.set(q.id, q);
    });

    // Grade answers
    let marksObtained = 0;
    const gradedAnswers = body.answers.map((answer) => {
      const question = questionsMap.get(answer.questionId);
      if (!question) return { ...answer, isCorrect: false, pointsAwarded: 0 };

      let isCorrect = false;
      let pointsAwarded = 0;

      if (question.type === 'mcq' || question.type === 'true_false') {
        // Check if selected option is correct
        if (answer.selectedOptions && answer.selectedOptions.length > 0) {
          const selectedId = answer.selectedOptions[0];
          const correctOption = question.options?.find(opt => opt.isCorrect);
          isCorrect = correctOption?.id === selectedId;
        }
      } else if (question.type === 'short_answer') {
        // Simple text matching
        isCorrect = answer.textAnswer?.toLowerCase().trim() === question.correctAnswer?.toLowerCase().trim();
      }

      if (isCorrect) {
        pointsAwarded = question.marks;
        marksObtained += pointsAwarded;
      } else if (quiz.negativeMarking && answer.selectedOptions?.length) {
        // Apply negative marking for wrong answers
        pointsAwarded = -quiz.negativeMarkValue;
        marksObtained += pointsAwarded;
      }

      return {
        ...answer,
        isCorrect,
        pointsAwarded,
      };
    });

    // Ensure marks don't go below 0
    marksObtained = Math.max(0, marksObtained);

    // Calculate results
    const percentage = (marksObtained / quiz.totalMarks) * 100;
    const passed = marksObtained >= quiz.passingMarks;
    const now = new Date();
    const startedAt = submission.startedAt instanceof Date 
      ? submission.startedAt 
      : new Date(submission.startedAt);
    const timeTaken = Math.round((now.getTime() - startedAt.getTime()) / 60000);

    // Update submission
    const updateData: any = {
      answers: gradedAnswers,
      submittedAt: now,
      timeTaken,
      marksObtained,
      score: marksObtained,
      percentage,
      passed,
      status: 'submitted',
      tabSwitchCount: body.tabSwitches || 0,
      fullscreenExits: body.fullscreenExits || 0,
      updatedAt: FieldValue.serverTimestamp(),
    };

    if (body.timedOut) {
      updateData.suspiciousActivity = FieldValue.arrayUnion('Quiz timed out');
    }

    await db.collection('submissions').doc(body.submissionId).update(updateData);

    // Get updated submission
    const updatedDoc = await db.collection('submissions').doc(body.submissionId).get();
    const updatedSubmission = { id: updatedDoc.id, ...updatedDoc.data() } as QuizSubmission;

    return successResponse({
      ...updatedSubmission,
      showResults: quiz.showResultsAfterSubmit,
    }, 'Quiz submitted successfully');

  } catch (error: any) {
    console.error('Submit quiz error:', error);
    return errorResponse('Failed to submit quiz: ' + error.message, 500);
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

// Quizzes API: Create and List quizzes
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  createdResponse, 
  verifyAuthToken, 
  checkRole,
  getPaginationParams,
  generateId 
} from '@/lib/api-utils';
import { CreateQuizRequest, Quiz, Course } from '@/lib/types';

// GET: List quizzes
export async function GET(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const { page, pageSize } = getPaginationParams(request);
    const url = new URL(request.url);
    const courseId = url.searchParams.get('courseId');

    let query = db.collection('quizzes').where('isDeleted', '==', false);

    // Filter by course if provided
    if (courseId) {
      query = query.where('courseId', '==', courseId);
    }

    // Filter based on role
    switch (authUser.role) {
      case 'student':
        query = query.where('isPublished', '==', true);
        break;
      case 'teacher':
        const teacherParam = url.searchParams.get('teacher');
        if (teacherParam === 'me') {
          query = query.where('teacherId', '==', authUser.uid);
        } else {
          query = query.where('isPublished', '==', true);
        }
        break;
      case 'admin':
        // Admins see all
        break;
    }

    const snapshot = await query
      .orderBy('createdAt', 'desc')
      .limit(pageSize)
      .offset((page - 1) * pageSize)
      .get();

    const quizzes: Quiz[] = [];
    snapshot.forEach((doc) => {
      quizzes.push({ id: doc.id, ...doc.data() } as Quiz);
    });

    return successResponse({
      quizzes,
      page,
      pageSize,
    });

  } catch (error: any) {
    console.error('List quizzes error:', error);
    return errorResponse('Failed to fetch quizzes: ' + error.message, 500);
  }
}

// POST: Create a new quiz
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Only teachers and admins can create quizzes
    if (!checkRole(authUser, ['teacher', 'admin'])) {
      return errorResponse('Forbidden: Only teachers and admins can create quizzes', 403);
    }

    const body: CreateQuizRequest = await request.json();

    // Validate required fields
    if (!body.title) {
      return errorResponse('Title is required', 400);
    }
    if (!body.courseId) {
      return errorResponse('Course ID is required', 400);
    }
    if (!body.totalMarks || body.totalMarks <= 0) {
      return errorResponse('Total marks must be greater than 0', 400);
    }
    if (!body.duration || body.duration <= 0) {
      return errorResponse('Duration must be greater than 0', 400);
    }

    // Verify course exists and user is the teacher or admin
    const courseDoc = await db.collection('courses').doc(body.courseId).get();
    if (!courseDoc.exists) {
      return errorResponse('Course not found', 404);
    }

    const course = courseDoc.data() as Course;
    if (authUser.role !== 'admin' && course.teacherId !== authUser.uid) {
      return errorResponse('Forbidden: You can only create quizzes for your own courses', 403);
    }

    // Create quiz
    const now = new Date();
    const quizId = generateId();
    const quiz: Quiz = {
      id: quizId,
      courseId: body.courseId,
      courseTitle: course.title,
      teacherId: authUser.uid,
      title: body.title,
      description: body.description || '',
      duration: body.duration,
      totalMarks: body.totalMarks,
      passingMarks: body.passingMarks || body.totalMarks * 0.4,
      negativeMarking: body.negativeMarking || false,
      negativeMarkValue: body.negativeMarkValue || 0,
      questionCount: 0,
      shuffleQuestions: body.shuffleQuestions || false,
      shuffleOptions: body.shuffleOptions || false,
      showResultsAfterSubmit: body.showResultsAfterSubmit ?? true,
      allowReview: body.allowReview || false,
      maxAttempts: body.maxAttempts || 1,
      instructions: body.instructions || '',
      deadline: body.deadline ? new Date(body.deadline) : undefined,
      preventTabSwitch: body.preventTabSwitch || false,
      maxTabSwitches: body.maxTabSwitches || 3,
      requireFullscreen: body.requireFullscreen || false,
      disableCopyPaste: body.disableCopyPaste || false,
      enableProctoring: body.enableProctoring || false,
      randomizeQuestionOrder: body.randomizeQuestionOrder || false,
      timePerQuestion: body.timePerQuestion || 0,
      lockAfterSubmit: true,
      allowTeacherResume: body.allowTeacherResume || false,
      allowTeacherExtendTime: body.allowTeacherExtendTime || false,
      isPublished: body.isPublished || false,
      createdAt: now,
      updatedAt: now,
      isDeleted: false,
    };

    await db.collection('quizzes').doc(quizId).set(quiz);

    return createdResponse(quiz, 'Quiz created successfully');

  } catch (error: any) {
    console.error('Create quiz error:', error);
    return errorResponse('Failed to create quiz: ' + error.message, 500);
  }
}

export async function OPTIONS() {
  return new Response(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

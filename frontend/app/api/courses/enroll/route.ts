// Courses API: Enroll in a course
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
import { Course, User, Enrollment } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

// POST: Enroll in a course
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const body = await request.json();
    const { courseId } = body;

    if (!courseId) {
      return errorResponse('courseId is required', 400);
    }

    // Get course
    const courseDoc = await db.collection('courses').doc(courseId).get();
    if (!courseDoc.exists) {
      return errorResponse('Course not found', 404);
    }

    const course = courseDoc.data() as Course;

    if (!course.isPublished || course.isDeleted) {
      return errorResponse('Course is not available for enrollment', 400);
    }

    // Check if already enrolled
    const existingEnrollment = await db.collection('enrollments')
      .where('studentId', '==', authUser.uid)
      .where('courseId', '==', courseId)
      .limit(1)
      .get();

    if (!existingEnrollment.empty) {
      return errorResponse('Already enrolled in this course', 400);
    }

    // Get student info
    const userDoc = await db.collection('users').doc(authUser.uid).get();
    const user = userDoc.data() as User;

    // Create enrollment
    const now = new Date();
    const enrollment: Enrollment = {
      enrollmentId: generateId(),
      studentId: authUser.uid,
      studentName: user.displayName,
      courseId,
      courseTitle: course.title,
      enrolledAt: now,
      progress: 0,
      completedMaterials: [],
      status: 'active',
      lastAccessedAt: now,
    };

    await db.collection('enrollments').doc(enrollment.enrollmentId).set(enrollment);

    // Update course enrollment count
    await db.collection('courses').doc(courseId).update({
      enrollmentCount: FieldValue.increment(1),
    });

    return createdResponse(enrollment, 'Enrolled successfully');

  } catch (error: any) {
    console.error('Enroll error:', error);
    return errorResponse('Failed to enroll: ' + error.message, 500);
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

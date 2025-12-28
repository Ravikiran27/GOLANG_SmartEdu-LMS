// Courses API: Create and List courses
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
import { CreateCourseRequest, Course, User } from '@/lib/types';

// GET: List courses
export async function GET(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const { page, pageSize } = getPaginationParams(request);
    const url = new URL(request.url);
    const teacherParam = url.searchParams.get('teacher');

    let query = db.collection('courses').where('isDeleted', '==', false);

    // Filter based on role
    switch (authUser.role) {
      case 'student':
        query = query.where('isPublished', '==', true);
        break;
      case 'teacher':
        if (teacherParam === 'me') {
          query = query.where('teacherId', '==', authUser.uid);
        } else {
          query = query.where('isPublished', '==', true);
        }
        break;
      case 'admin':
        // Admins see all non-deleted courses
        break;
      default:
        query = query.where('isPublished', '==', true);
    }

    // Get courses with pagination
    const snapshot = await query
      .orderBy('createdAt', 'desc')
      .limit(pageSize)
      .offset((page - 1) * pageSize)
      .get();

    const courses: Course[] = [];
    snapshot.forEach((doc) => {
      courses.push(doc.data() as Course);
    });

    return successResponse({
      courses,
      page,
      pageSize,
    });

  } catch (error: any) {
    console.error('List courses error:', error);
    return errorResponse('Failed to fetch courses: ' + error.message, 500);
  }
}

// POST: Create a new course
export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Only teachers and admins can create courses
    if (!checkRole(authUser, ['teacher', 'admin'])) {
      return errorResponse('Forbidden: Only teachers and admins can create courses', 403);
    }

    const body: CreateCourseRequest = await request.json();

    // Validate required fields
    if (!body.title || !body.description || !body.category || !body.difficulty) {
      return errorResponse('Title, description, category, and difficulty are required', 400);
    }

    // Get user info for teacher name
    const userDoc = await db.collection('users').doc(authUser.uid).get();
    if (!userDoc.exists) {
      return errorResponse('User not found', 404);
    }
    const user = userDoc.data() as User;

    // Create course
    const now = new Date();
    const course: Course = {
      courseId: generateId(),
      title: body.title,
      description: body.description,
      syllabus: body.syllabus || '',
      teacherId: authUser.uid,
      teacherName: user.displayName,
      category: body.category,
      difficulty: body.difficulty as 'beginner' | 'intermediate' | 'advanced',
      thumbnail: body.thumbnail || '',
      materials: [],
      enrollmentCount: 0,
      isPublished: false,
      createdAt: now,
      updatedAt: now,
      isDeleted: false,
    };

    await db.collection('courses').doc(course.courseId).set(course);

    return createdResponse(course, 'Course created successfully');

  } catch (error: any) {
    console.error('Create course error:', error);
    return errorResponse('Failed to create course: ' + error.message, 500);
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

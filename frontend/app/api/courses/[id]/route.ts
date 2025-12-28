// Courses API: Get, Update, Delete single course
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken, 
  checkRole 
} from '@/lib/api-utils';
import { UpdateCourseRequest, Course } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

interface RouteParams {
  params: Promise<{ id: string }>;
}

// GET: Get single course
export async function GET(request: NextRequest, { params }: RouteParams) {
  try {
    const { id } = await params;
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const courseDoc = await db.collection('courses').doc(id).get();
    if (!courseDoc.exists) {
      return errorResponse('Course not found', 404);
    }

    const course = courseDoc.data() as Course;

    // Check access permissions
    if (course.isDeleted) {
      return errorResponse('Course not found', 404);
    }

    // Students can only see published courses
    if (authUser.role === 'student' && !course.isPublished) {
      return errorResponse('Course not available', 403);
    }

    // Teachers can only see their own unpublished courses
    if (authUser.role === 'teacher' && !course.isPublished && course.teacherId !== authUser.uid) {
      return errorResponse('Course not available', 403);
    }

    return successResponse(course);

  } catch (error: any) {
    console.error('Get course error:', error);
    return errorResponse('Failed to get course: ' + error.message, 500);
  }
}

// PUT: Update course
export async function PUT(request: NextRequest, { params }: RouteParams) {
  try {
    const { id } = await params;
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const courseDoc = await db.collection('courses').doc(id).get();
    if (!courseDoc.exists) {
      return errorResponse('Course not found', 404);
    }

    const course = courseDoc.data() as Course;

    // Only course teacher or admin can update
    if (authUser.role !== 'admin' && course.teacherId !== authUser.uid) {
      return errorResponse('Forbidden: You can only update your own courses', 403);
    }

    const body: UpdateCourseRequest = await request.json();

    // Build update object
    const updateData: any = {
      updatedAt: FieldValue.serverTimestamp(),
    };

    if (body.title) updateData.title = body.title;
    if (body.description) updateData.description = body.description;
    if (body.syllabus !== undefined) updateData.syllabus = body.syllabus;
    if (body.category) updateData.category = body.category;
    if (body.difficulty) updateData.difficulty = body.difficulty;
    if (body.thumbnail !== undefined) updateData.thumbnail = body.thumbnail;
    if (body.materials) updateData.materials = body.materials;
    if (body.isPublished !== undefined) updateData.isPublished = body.isPublished;

    await db.collection('courses').doc(id).update(updateData);

    // Get updated course
    const updatedDoc = await db.collection('courses').doc(id).get();
    const updatedCourse = updatedDoc.data() as Course;

    return successResponse(updatedCourse, 'Course updated successfully');

  } catch (error: any) {
    console.error('Update course error:', error);
    return errorResponse('Failed to update course: ' + error.message, 500);
  }
}

// DELETE: Soft delete course
export async function DELETE(request: NextRequest, { params }: RouteParams) {
  try {
    const { id } = await params;
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const courseDoc = await db.collection('courses').doc(id).get();
    if (!courseDoc.exists) {
      return errorResponse('Course not found', 404);
    }

    const course = courseDoc.data() as Course;

    // Only course teacher or admin can delete
    if (authUser.role !== 'admin' && course.teacherId !== authUser.uid) {
      return errorResponse('Forbidden: You can only delete your own courses', 403);
    }

    // Soft delete
    await db.collection('courses').doc(id).update({
      isDeleted: true,
      updatedAt: FieldValue.serverTimestamp(),
    });

    return successResponse({ courseId: id }, 'Course deleted successfully');

  } catch (error: any) {
    console.error('Delete course error:', error);
    return errorResponse('Failed to delete course: ' + error.message, 500);
  }
}

export async function OPTIONS() {
  return new Response(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

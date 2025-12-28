// Courses API: Get my enrollments
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { 
  successResponse, 
  errorResponse, 
  verifyAuthToken,
  getPaginationParams 
} from '@/lib/api-utils';
import { Enrollment } from '@/lib/types';

// GET: Get my enrollments
export async function GET(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const { page, pageSize } = getPaginationParams(request);

    const snapshot = await db.collection('enrollments')
      .where('studentId', '==', authUser.uid)
      .orderBy('enrolledAt', 'desc')
      .limit(pageSize)
      .offset((page - 1) * pageSize)
      .get();

    const enrollments: Enrollment[] = [];
    snapshot.forEach((doc) => {
      enrollments.push(doc.data() as Enrollment);
    });

    return successResponse({
      enrollments,
      page,
      pageSize,
    });

  } catch (error: any) {
    console.error('Get enrollments error:', error);
    return errorResponse('Failed to get enrollments: ' + error.message, 500);
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

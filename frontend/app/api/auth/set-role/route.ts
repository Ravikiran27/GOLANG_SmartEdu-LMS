// Auth API: Set user role (Admin only)
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { auth, db } from '@/lib/firebase-admin';
import { successResponse, errorResponse, verifyAuthToken, checkRole } from '@/lib/api-utils';
import { FieldValue } from 'firebase-admin/firestore';

export async function POST(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Only admins can set roles
    if (!checkRole(authUser, ['admin'])) {
      return errorResponse('Forbidden: Admin access required', 403);
    }

    const body = await request.json();
    const { userId, role } = body;

    if (!userId || !role) {
      return errorResponse('userId and role are required', 400);
    }

    if (!['admin', 'teacher', 'student'].includes(role)) {
      return errorResponse('Invalid role. Must be admin, teacher, or student', 400);
    }

    // Update custom claims
    await auth.setCustomUserClaims(userId, { role });

    // Update Firestore document
    await db.collection('users').doc(userId).update({
      role,
      updatedAt: FieldValue.serverTimestamp(),
    });

    return successResponse({ userId, role }, 'Role updated successfully');

  } catch (error: any) {
    console.error('Set role error:', error);
    return errorResponse('Failed to set role: ' + error.message, 500);
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

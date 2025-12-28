// Auth API: Update user profile
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { successResponse, errorResponse, verifyAuthToken } from '@/lib/api-utils';
import { UpdateUserRequest, User } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

export async function PUT(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    const body: UpdateUserRequest = await request.json();

    // Build update object
    const updateData: any = {
      updatedAt: FieldValue.serverTimestamp(),
    };

    if (body.displayName) updateData.displayName = body.displayName;
    if (body.photoURL) updateData.photoURL = body.photoURL;
    if (body.department) updateData['metadata.department'] = body.department;
    if (body.rollNumber) updateData['metadata.rollNumber'] = body.rollNumber;
    if (body.employeeId) updateData['metadata.employeeId'] = body.employeeId;

    // Update user document
    await db.collection('users').doc(authUser.uid).update(updateData);

    // Get updated user
    const userDoc = await db.collection('users').doc(authUser.uid).get();
    const user = userDoc.data() as User;

    return successResponse(user, 'Profile updated successfully');

  } catch (error: any) {
    console.error('Update profile error:', error);
    return errorResponse('Failed to update profile: ' + error.message, 500);
  }
}

export async function OPTIONS() {
  return new Response(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'PUT, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    },
  });
}

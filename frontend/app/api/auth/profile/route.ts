// Auth API: Get user profile
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { db } from '@/lib/firebase-admin';
import { successResponse, errorResponse, verifyAuthToken } from '@/lib/api-utils';
import { User } from '@/lib/types';
import { FieldValue } from 'firebase-admin/firestore';

export async function GET(request: NextRequest) {
  try {
    const authUser = await verifyAuthToken(request);
    if (!authUser) {
      return errorResponse('Unauthorized', 401);
    }

    // Get user document
    const userDoc = await db.collection('users').doc(authUser.uid).get();
    if (!userDoc.exists) {
      return errorResponse('User not found', 404);
    }

    const user = userDoc.data() as User;

    // Update last login
    await db.collection('users').doc(authUser.uid).update({
      'metadata.lastLogin': FieldValue.serverTimestamp(),
    });

    return successResponse(user);

  } catch (error: any) {
    console.error('Profile error:', error);
    return errorResponse('Failed to get profile: ' + error.message, 500);
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

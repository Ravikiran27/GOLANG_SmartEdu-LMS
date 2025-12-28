// Auth API: Register new user
export const dynamic = 'force-dynamic';

import { NextRequest } from 'next/server';
import { auth, db } from '@/lib/firebase-admin';
import { successResponse, errorResponse, createdResponse } from '@/lib/api-utils';
import { CreateUserRequest, User } from '@/lib/types';

export async function POST(request: NextRequest) {
  try {
    const body: CreateUserRequest = await request.json();

    // Validate required fields
    if (!body.email || !body.password || !body.displayName || !body.role) {
      return errorResponse('Email, password, displayName, and role are required', 400);
    }

    // Validate role
    if (!['admin', 'teacher', 'student'].includes(body.role)) {
      return errorResponse('Invalid role. Must be admin, teacher, or student', 400);
    }

    // Create user in Firebase Auth
    const firebaseUser = await auth.createUser({
      email: body.email,
      password: body.password,
      displayName: body.displayName,
      emailVerified: false,
      disabled: false,
    });

    // Set custom claims for role
    await auth.setCustomUserClaims(firebaseUser.uid, { role: body.role });

    // Create user document in Firestore
    const now = new Date();
    const user: User = {
      uid: firebaseUser.uid,
      email: body.email,
      displayName: body.displayName,
      role: body.role,
      photoURL: '',
      isActive: true,
      createdAt: now,
      updatedAt: now,
      metadata: {
        lastLogin: now,
        department: body.department,
        rollNumber: body.rollNumber,
        employeeId: body.employeeId,
      },
    };

    await db.collection('users').doc(firebaseUser.uid).set(user);

    return createdResponse({
      uid: firebaseUser.uid,
      email: body.email,
      role: body.role,
    }, 'User registered successfully');

  } catch (error: any) {
    console.error('Register error:', error);
    
    // Handle specific Firebase errors
    if (error.code === 'auth/email-already-exists') {
      return errorResponse('Email already in use', 400);
    }
    if (error.code === 'auth/invalid-email') {
      return errorResponse('Invalid email address', 400);
    }
    if (error.code === 'auth/weak-password') {
      return errorResponse('Password is too weak', 400);
    }
    
    return errorResponse('Failed to create user: ' + error.message, 500);
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

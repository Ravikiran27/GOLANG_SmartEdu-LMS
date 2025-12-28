// API Utilities for Next.js API Routes
import { NextRequest, NextResponse } from 'next/server';
import { auth } from './firebase-admin';

// Standard API Response
export function jsonResponse(data: any, status: number = 200) {
  return NextResponse.json(data, { status });
}

export function successResponse(data: any, message?: string) {
  return jsonResponse({
    success: true,
    data,
    message: message || 'Success',
  });
}

export function errorResponse(message: string, status: number = 400) {
  return jsonResponse({
    success: false,
    error: message,
  }, status);
}

export function createdResponse(data: any, message?: string) {
  return jsonResponse({
    success: true,
    data,
    message: message || 'Created successfully',
  }, 201);
}

// CORS Headers
export function corsHeaders() {
  return {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
  };
}

// Auth Token Verification
export interface AuthUser {
  uid: string;
  email: string;
  role: string;
}

export async function verifyAuthToken(request: NextRequest): Promise<AuthUser | null> {
  try {
    const authHeader = request.headers.get('Authorization');
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return null;
    }

    const token = authHeader.split('Bearer ')[1];
    const decodedToken = await auth.verifyIdToken(token);
    
    return {
      uid: decodedToken.uid,
      email: decodedToken.email || '',
      role: (decodedToken.role as string) || 'student',
    };
  } catch (error) {
    console.error('Auth verification failed:', error);
    return null;
  }
}

// Role-based authorization
export function checkRole(user: AuthUser, allowedRoles: string[]): boolean {
  return allowedRoles.includes(user.role);
}

// Pagination helpers
export function getPaginationParams(request: NextRequest) {
  const url = new URL(request.url);
  const page = parseInt(url.searchParams.get('page') || '1', 10);
  const pageSize = parseInt(url.searchParams.get('pageSize') || '10', 10);
  return { page: Math.max(1, page), pageSize: Math.min(50, Math.max(1, pageSize)) };
}

// UUID generator
export function generateId(): string {
  return crypto.randomUUID();
}

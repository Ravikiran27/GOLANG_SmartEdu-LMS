// API Client for Backend Communication
import { auth } from './firebase';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

async function getAuthToken(): Promise<string | null> {
  if (!auth) return null;
  const user = auth.currentUser;
  if (!user) return null;
  return await user.getIdToken();
}

export async function apiRequest<T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const token = await getAuthToken();

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string> || {}),
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'API request failed');
    }

    return data;
  } catch (error) {
    console.error('API Error:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

// Specific API methods
export const api = {
  // Auth
  auth: {
    getProfile: () => apiRequest('/api/auth/profile'),
    updateProfile: (data: any) =>
      apiRequest('/api/auth/update', {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
  },

  // Courses
  courses: {
    list: () => apiRequest('/api/courses/list'),
    get: (id: string) => apiRequest(`/api/courses/get?id=${id}`),
    create: (data: any) =>
      apiRequest('/api/courses/create', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    update: (id: string, data: any) =>
      apiRequest(`/api/courses/update?id=${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: string) =>
      apiRequest(`/api/courses/delete?id=${id}`, { method: 'DELETE' }),
    enroll: (courseId: string) =>
      apiRequest('/api/courses/enroll', {
        method: 'POST',
        body: JSON.stringify({ courseId }),
      }),
    myEnrollments: () => apiRequest('/api/courses/my-enrollments'),
  },

  // Quizzes
  quizzes: {
    list: (courseId?: string) =>
      apiRequest(`/api/quizzes/list${courseId ? `?courseId=${courseId}` : ''}`),
    get: (id: string) => apiRequest(`/api/quizzes/get?id=${id}`),
    create: (data: any) =>
      apiRequest('/api/quizzes/create', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    start: (quizId: string) =>
      apiRequest('/api/quizzes/start', {
        method: 'POST',
        body: JSON.stringify({ quizId }),
      }),
    submit: (quizId: string, answers: any[]) =>
      apiRequest('/api/quizzes/submit', {
        method: 'POST',
        body: JSON.stringify({ quizId, answers }),
      }),
  },

  // Add more API endpoints as needed
};

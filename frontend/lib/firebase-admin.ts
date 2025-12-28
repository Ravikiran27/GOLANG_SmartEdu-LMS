// Firebase Admin SDK for Server-side operations
import { initializeApp, getApps, cert, App } from 'firebase-admin/app';
import { getAuth, Auth } from 'firebase-admin/auth';
import { getFirestore, Firestore } from 'firebase-admin/firestore';

// Firebase Admin credentials
const firebaseAdminConfig = {
  projectId: process.env.FIREBASE_PROJECT_ID || "flutterdemo-6d4ac",
  clientEmail: process.env.FIREBASE_CLIENT_EMAIL || "",
  privateKey: (process.env.FIREBASE_PRIVATE_KEY || "").replace(/\\n/g, '\n'),
};

let app: App | null = null;
let adminAuth: Auth | null = null;
let adminDb: Firestore | null = null;
let initialized = false;

function getFirebaseAdmin() {
  if (initialized) {
    return { app, adminAuth, adminDb };
  }

  // Check if we have valid credentials
  if (!firebaseAdminConfig.clientEmail || !firebaseAdminConfig.privateKey) {
    console.warn('Firebase Admin credentials not available - API routes will not work');
    initialized = true;
    return { app: null, adminAuth: null, adminDb: null };
  }

  try {
    if (getApps().length === 0) {
      app = initializeApp({
        credential: cert(firebaseAdminConfig),
        projectId: firebaseAdminConfig.projectId,
      });
    } else {
      app = getApps()[0];
    }
    
    adminAuth = getAuth(app);
    adminDb = getFirestore(app);
    initialized = true;
    
    return { app, adminAuth, adminDb };
  } catch (error) {
    console.error('Failed to initialize Firebase Admin:', error);
    initialized = true;
    return { app: null, adminAuth: null, adminDb: null };
  }
}

// Lazy getters
export function getAdminAuth(): Auth {
  const { adminAuth } = getFirebaseAdmin();
  if (!adminAuth) {
    throw new Error('Firebase Admin Auth not initialized - check FIREBASE_CLIENT_EMAIL and FIREBASE_PRIVATE_KEY env vars');
  }
  return adminAuth;
}

export function getAdminDb(): Firestore {
  const { adminDb } = getFirebaseAdmin();
  if (!adminDb) {
    throw new Error('Firebase Admin Firestore not initialized - check FIREBASE_CLIENT_EMAIL and FIREBASE_PRIVATE_KEY env vars');
  }
  return adminDb;
}

// Legacy exports for backward compatibility
export const auth = {
  get instance() { return getAdminAuth(); },
  createUser: (...args: Parameters<Auth['createUser']>) => getAdminAuth().createUser(...args),
  setCustomUserClaims: (...args: Parameters<Auth['setCustomUserClaims']>) => getAdminAuth().setCustomUserClaims(...args),
  verifyIdToken: (...args: Parameters<Auth['verifyIdToken']>) => getAdminAuth().verifyIdToken(...args),
};

export const db = {
  get instance() { return getAdminDb(); },
  collection: (...args: Parameters<Firestore['collection']>) => getAdminDb().collection(...args),
};

export default app;

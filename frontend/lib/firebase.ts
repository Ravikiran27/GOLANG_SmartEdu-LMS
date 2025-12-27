// Firebase Client Configuration
import { initializeApp, getApps, FirebaseApp } from 'firebase/app';
import { getAuth, Auth } from 'firebase/auth';
import { getStorage, FirebaseStorage } from 'firebase/storage';

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
};

// Check if we're in a browser environment and have valid config
const isValidConfig = typeof window !== 'undefined' && firebaseConfig.apiKey;

// Initialize Firebase (singleton pattern) only in browser with valid config
const app = isValidConfig && getApps().length === 0 ? initializeApp(firebaseConfig) : getApps()[0];

export const auth = isValidConfig && app ? getAuth(app) : null;
export const storage = isValidConfig && app ? getStorage(app) : null;

export default app;

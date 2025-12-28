// Firebase Client Configuration
import { initializeApp, getApps, FirebaseApp } from 'firebase/app';
import { getAuth, Auth } from 'firebase/auth';
import { getStorage, FirebaseStorage } from 'firebase/storage';

const firebaseConfig = {
  apiKey: "AIzaSyDV2vK5mXHKnFjPBF4UwoctJMF_uzTEpfg",
  authDomain: "flutterdemo-6d4ac.firebaseapp.com",
  projectId: "flutterdemo-6d4ac",
  storageBucket: "flutterdemo-6d4ac.firebasestorage.app",
  messagingSenderId: "764379386105",
  appId: "1:764379386105:web:c62baca4253f4d351cd228",
};

// Check if we're in a browser environment
const isBrowser = typeof window !== 'undefined';

// Initialize Firebase (singleton pattern) only in browser
let app: FirebaseApp | undefined;
if (isBrowser) {
  app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApps()[0];
}

export const auth = isBrowser && app ? getAuth(app) : null;
export const storage = isBrowser && app ? getStorage(app) : null;

export default app;

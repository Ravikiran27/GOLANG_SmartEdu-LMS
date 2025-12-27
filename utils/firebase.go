package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firebaseApp  *firebase.App
	authClient   *auth.Client
	firestoreClient *firestore.Client
	once         sync.Once
	initError    error
)

// InitFirebase initializes Firebase Admin SDK (singleton pattern)
func InitFirebase(ctx context.Context) error {
	once.Do(func() {
		// Read Firebase credentials from environment
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")
		clientEmail := os.Getenv("FIREBASE_CLIENT_EMAIL")

		if projectID == "" || privateKey == "" || clientEmail == "" {
			log.Println("ERROR: Missing Firebase credentials in environment variables")
			initError = &FirebaseError{Message: "Missing Firebase credentials"}
			return
		}

		// Create service account credentials JSON
		credentials := map[string]string{
			"type":                        "service_account",
			"project_id":                  projectID,
			"private_key":                 privateKey,
			"client_email":                clientEmail,
			"token_uri":                   "https://oauth2.googleapis.com/token",
			"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		}

		credentialsJSON, err := json.Marshal(credentials)
		if err != nil {
			log.Printf("ERROR: Failed to marshal credentials: %v", err)
			initError = err
			return
		}

		// Initialize Firebase App
		conf := &firebase.Config{
			ProjectID:     projectID,
			StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
		}

		opt := option.WithCredentialsJSON(credentialsJSON)
		firebaseApp, initError = firebase.NewApp(ctx, conf, opt)
		if initError != nil {
			log.Printf("ERROR: Failed to initialize Firebase App: %v", initError)
			return
		}

		// Initialize Auth Client
		authClient, initError = firebaseApp.Auth(ctx)
		if initError != nil {
			log.Printf("ERROR: Failed to initialize Auth Client: %v", initError)
			return
		}

		// Initialize Firestore Client
		firestoreClient, initError = firebaseApp.Firestore(ctx)
		if initError != nil {
			log.Printf("ERROR: Failed to initialize Firestore Client: %v", initError)
			return
		}

		log.Println("✅ Firebase initialized successfully")
	})

	return initError
}

// GetAuthClient returns the Firebase Auth client
func GetAuthClient(ctx context.Context) (*auth.Client, error) {
	if err := InitFirebase(ctx); err != nil {
		return nil, err
	}
	if authClient == nil {
		return nil, &FirebaseError{Message: "Auth client not initialized"}
	}
	return authClient, nil
}

// GetFirestoreClient returns the Firestore client
func GetFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	if err := InitFirebase(ctx); err != nil {
		return nil, err
	}
	if firestoreClient == nil {
		return nil, &FirebaseError{Message: "Firestore client not initialized"}
	}
	return firestoreClient, nil
}

// CloseFirestore closes Firestore connection (call on shutdown)
func CloseFirestore() error {
	if firestoreClient != nil {
		return firestoreClient.Close()
	}
	return nil
}

// FirebaseError represents a Firebase-related error
type FirebaseError struct {
	Message string
}

func (e *FirebaseError) Error() string {
	return e.Message
}

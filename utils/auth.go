package utils

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware validates Firebase token and extracts user claims
func AuthMiddleware(next http.HandlerFunc, requiredRole ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Verify token with Firebase
		authClient, err := GetAuthClient(ctx)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to initialize auth client")
			return
		}

		decodedToken, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Check if role is required
		if len(requiredRole) > 0 {
			userRole, ok := decodedToken.Claims["role"].(string)
			if !ok {
				RespondError(w, http.StatusForbidden, "User role not found")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range requiredRole {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				RespondError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "uid", decodedToken.UID)
		ctx = context.WithValue(ctx, "email", decodedToken.Claims["email"])
		
		if role, ok := decodedToken.Claims["role"]; ok {
			ctx = context.WithValue(ctx, "role", role)
		}

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserFromContext extracts user information from context
func GetUserFromContext(ctx context.Context) (uid, email, role string) {
	if uidVal := ctx.Value("uid"); uidVal != nil {
		uid = uidVal.(string)
	}
	if emailVal := ctx.Value("email"); emailVal != nil {
		email = emailVal.(string)
	}
	if roleVal := ctx.Value("role"); roleVal != nil {
		role = roleVal.(string)
	}
	return
}

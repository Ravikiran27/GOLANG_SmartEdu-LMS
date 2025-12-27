package models

import "time"

// User represents a user in the system
type User struct {
	UID         string    `firestore:"uid" json:"uid"`
	Email       string    `firestore:"email" json:"email"`
	DisplayName string    `firestore:"displayName" json:"displayName"`
	Role        string    `firestore:"role" json:"role"` // admin | teacher | student
	PhotoURL    string    `firestore:"photoURL,omitempty" json:"photoURL,omitempty"`
	IsActive    bool      `firestore:"isActive" json:"isActive"`
	CreatedAt   time.Time `firestore:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt" json:"updatedAt"`
	Metadata    UserMetadata `firestore:"metadata" json:"metadata"`
}

// UserMetadata contains additional user information
type UserMetadata struct {
	LastLogin   time.Time `firestore:"lastLogin,omitempty" json:"lastLogin,omitempty"`
	Department  string    `firestore:"department,omitempty" json:"department,omitempty"`
	RollNumber  string    `firestore:"rollNumber,omitempty" json:"rollNumber,omitempty"`
	EmployeeID  string    `firestore:"employeeId,omitempty" json:"employeeId,omitempty"`
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	DisplayName string `json:"displayName" validate:"required"`
	Role        string `json:"role" validate:"required,oneof=admin teacher student"`
	Department  string `json:"department,omitempty"`
	RollNumber  string `json:"rollNumber,omitempty"`
	EmployeeID  string `json:"employeeId,omitempty"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	DisplayName string `json:"displayName,omitempty"`
	PhotoURL    string `json:"photoURL,omitempty"`
	Department  string `json:"department,omitempty"`
	RollNumber  string `json:"rollNumber,omitempty"`
	EmployeeID  string `json:"employeeId,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

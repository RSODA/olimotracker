package auth

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

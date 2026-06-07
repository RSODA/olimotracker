package user

import "github.com/google/uuid"

type User struct {
	ID       *uuid.UUID
	Email    string
	Username string
	APIKey   *uuid.UUID
}

type UserProfileResponse struct {
	ID       *uuid.UUID
	Email    string
	Username string
	APIKey   *uuid.UUID
}

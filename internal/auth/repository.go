package auth

import "context"

type Repository interface {
	CreateUser(ctx context.Context, email, username, password string)
}

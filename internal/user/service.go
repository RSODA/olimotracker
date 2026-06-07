package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Service interface {
	GetUserByID(ctx context.Context, id *uuid.UUID) (*User, error)
}

type service struct {
	repo Repository
	l    *slog.Logger
}

func NewService(repo Repository, l *slog.Logger) Service {
	return &service{
		repo: repo,
		l:    l,
	}
}

func (s *service) GetUserByID(ctx context.Context, id *uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

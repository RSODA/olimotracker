package auth

import (
	"context"
	"log/slog"
	"olimotracker/pkg/jwttoken"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (string, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type service struct {
	repo     Repository
	l        *slog.Logger
	tokenGen *jwttoken.TokenGenerator
}

func NewService(repo Repository, l *slog.Logger, tokenGen *jwttoken.TokenGenerator) Service {
	return &service{
		repo:     repo,
		l:        l,
		tokenGen: tokenGen,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.l.Error("failed to hash password", "err", err)
		return "", err
	}

	id, err := s.repo.CreateUser(ctx, User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	})
	if err != nil {
		return "", err
	}

	token, err := s.tokenGen.GenerateToken(id)
	if err != nil {
		s.l.Error("failed to generate token", "err", err)
		return "", err
	}

	return token, nil

}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.l.Error("invalid password", "err", err)
		return "", ErrInvalidPassword
	}

	token, err := s.tokenGen.GenerateToken(&user.ID)
	if err != nil {
		s.l.Error("failed to generate token", "err", err)
		return "", err
	}

	return token, nil
}

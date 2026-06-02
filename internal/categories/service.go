package categories

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type Service interface {
	CreateCategory(ctx context.Context, userID *uuid.UUID, category *Category) (*Category, error)
	GetCategoriesByUserID(ctx context.Context, userID *uuid.UUID) ([]*Category, error)
	GetCategoryByID(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) (*Category, error)
	UpdateCategory(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) error
}

type service struct {
	repo Repository
	l    *slog.Logger
}

func NewService(repo Repository, l *slog.Logger) Service {
	return &service{repo: repo, l: l}
}

func (s *service) CreateCategory(ctx context.Context, userID *uuid.UUID, category *Category) (*Category, error) {
	if len(category.Color) == 0 || category.Color[0] != '#' {
		s.l.Error("error invalid color: ", "color", category.Color, "err", ErrInvalidColor)
		return nil, ErrInvalidColor
	}

	category.UserID = *userID

	res, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetCategoriesByUserID(ctx context.Context, userID *uuid.UUID) ([]*Category, error) {
	res, err := s.repo.GetCategoriesByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetCategoryByID(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) (*Category, error) {
	res, err := s.repo.GetCategoryByID(ctx, categoryID.String(), userID.String())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateCategory(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID, category *Category) (*Category, error) {
	if len(category.Color) != 0 {
		if category.Color[0] != '#' {
			s.l.Error("error invalid color: ", "color", category.Color, "err", ErrInvalidColor)
			return nil, ErrInvalidColor
		}
	}

	category.UserID = *userID
	category.ID = *categoryID

	res, err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) DeleteCategory(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) error {
	err := s.repo.DeleteCategory(ctx, categoryID.String(), userID.String())
	if err != nil {
		return err
	}

	return nil
}

package api

import (
	"context"
	"log/slog"
	"olimotracker/pkg/db"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Repository interface {
	GetUserByAPIKey(ctx context.Context, apiKey *uuid.UUID) (*User, error)
}

type repo struct {
	db db.DBClient
	l  *slog.Logger
}

func NewRepository(db db.DBClient, l *slog.Logger) Repository {
	return &repo{
		db: db,
		l:  l,
	}
}

func (r *repo) GetUserByAPIKey(ctx context.Context, apiKey *uuid.UUID) (*User, error) {
	builder := squirrel.Select(db.UsersIDColumn, db.UsersEmailColumn, db.UsersUsernameColumn, db.UsersAPIKeyColumn).
		From(db.UsersTable).
		Where(squirrel.Eq{db.UsersAPIKeyColumn: apiKey}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("failed to build query", "error", err)
		return nil, err
	}

	var user User
	if err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Username, &user.APIKey); err != nil {
		r.l.Error("failed to scan user", "error", err)
		return nil, err
	}

	return &user, nil
}

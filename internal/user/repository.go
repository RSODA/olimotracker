package user

import (
	"context"
	"log/slog"
	"olimotracker/pkg/db"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Repository interface {
	GetUserByID(ctx context.Context, id *uuid.UUID) (*User, error)
	UpdateAPIKeyByUserID(ctx context.Context, id *uuid.UUID, newAPI *uuid.UUID) error
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

func (r *repo) GetUserByID(ctx context.Context, id *uuid.UUID) (*User, error) {
	var user User

	builder := squirrel.Select(db.UsersIDColumn, db.UsersEmailColumn, db.UsersUsernameColumn, db.UsersAPIKeyColumn).
		From(db.UsersTable).
		Where(squirrel.Eq{db.UsersIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("failed to build query", "error", err)
		return nil, err
	}

	if err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Username, &user.APIKey); err != nil {
		r.l.Error("failed to scan user", "error", err)
		return nil, err
	}

	return &user, nil
}

func (r *repo) UpdateAPIKeyByUserID(ctx context.Context, id *uuid.UUID, newAPI *uuid.UUID) error {
	builder := squirrel.Update(db.UsersTable).
		Set(db.UsersAPIKeyColumn, newAPI).
		Where(squirrel.Eq{db.UsersIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		r.l.Error("error executing query: ", "err", err)
		return err
	}

	return nil
}

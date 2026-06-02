package auth

import (
	"context"
	"errors"
	"log/slog"
	"olimotracker/pkg/db"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateUser(ctx context.Context, u User) (*uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db db.DBClient
	l  *slog.Logger
}

func NewRepository(db db.DBClient, l *slog.Logger) Repository {
	return &repository{db: db, l: l}
}

func (r *repository) CreateUser(ctx context.Context, u User) (*uuid.UUID, error) {
	builder := squirrel.Insert(db.UsersTable).
		Columns(db.UsersEmailColumn, db.UsersPasswordColumn, db.UsersUsernameColumn).
		Values(u.Email, u.Password, u.Username).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING " + db.UsersIDColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("fail build query", "err", err)
		return nil, err
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserEmailAlreadyExists
		}
		r.l.Error("fail create user", "err", err)
		return nil, err
	}

	return &id, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	builder := squirrel.Select(db.UsersIDColumn, db.UsersEmailColumn, db.UsersPasswordColumn, db.UsersUsernameColumn).
		From(db.UsersTable).
		Where(squirrel.Eq{db.UsersEmailColumn: email}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("fail build query", "err", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		r.l.Error("fail get user", "err", err)
		return nil, err
	}

	return &user, nil
}

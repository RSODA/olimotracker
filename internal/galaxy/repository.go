package galaxy

import (
	"context"
	"log/slog"
	"olimotracker/pkg/db"

	"github.com/Masterminds/squirrel"
)

type Repository interface {
	UpdateAllSeeds(ctx context.Context) error
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

func (r *repo) UpdateAllSeeds(ctx context.Context) error {
	builder := squirrel.Update(db.UserStatsTable).
		Set(db.UserStatsGalaxySeedColumn, squirrel.Expr("floor(random() * 9999999)::bigint")).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("failed to build update all seeds query", "error", err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		r.l.Error("failed to update all seeds", "error", err)
		return err
	}

	return nil
}

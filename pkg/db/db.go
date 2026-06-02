package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// users table
	UsersTable = "users"

	UsersIDColumn        = "id"
	UsersEmailColumn     = "email"
	UsersPasswordColumn  = "password"
	UsersUsernameColumn  = "username"
	UsersAPIKeyColumn    = "api_key"
	UsersTgIDColumn      = "telegram_id"
	UsersCreatedAtColumn = "created_at"

	// categories table
	CategoriesTable = "categories"

	CategoriesIDColumn        = "id"
	CategoriesTitleColumn     = "title"
	CategoriesUserIDColumn    = "user_id"
	CategoriesColorColumn     = "color"
	CategoriesCreatedAtColumn = "created_at"

	// sessions table
	SessionsTable = "sessions"

	SessionsIDColumn         = "id"
	SessionsUserIDColumn     = "user_id"
	SessionsCategoryIDColumn = "category_id"
	SessionsDurationColumn   = "duration"
	SessionsNotesColumn      = "note"
	SessionsCreatedAtColumn  = "created_at"

	// user_stats table
	UserStatsTable = "user_stats"

	UserStatsIDColumn            = "id"
	UserStatsUserIDColumn        = "user_id"
	UserStatsTotalHoursColumn    = "total_hour"
	UserStatsXPColumn            = "xp"
	UserStatsLevelColumn         = "level"
	UserStatsCurrentStreakColumn = "current_streak"
	UserStatsMaxStreakColumn     = "max_streak"
	UserStatsGalaxySeedColumn    = "galaxy_seed"
	UserStatsCreatedAtColumn     = "created_at"
	UserStatsUpdatedAtColumn     = "updated_at"
)

type DBClient interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
	Ping(ctx context.Context) error
}

type DB struct {
	l  *slog.Logger
	db *pgxpool.Pool
}

func NewDB(pgx *pgxpool.Pool, l *slog.Logger) DBClient {
	return &DB{db: pgx, l: l}
}

func (d *DB) Ping(ctx context.Context) error {
	d.l.Info("pinging db")
	return d.db.Ping(ctx)
}

func (d *DB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	d.l.Debug("exec", "sql", sql, "args", args)
	tag, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		d.l.Error("exec error", "err", err)
	}
	return tag, err
}

func (d *DB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	d.l.Debug("query", "sql", sql, "args", args)
	rows, err := d.db.Query(ctx, sql, args...)
	if err != nil {
		d.l.Error("query error", "err", err)
	}
	return rows, err
}

func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	d.l.Debug("querying row", "sql", sql, "args", args)
	row := d.db.QueryRow(ctx, sql, args...)
	if row == nil {
		d.l.Error("query row error", "err", row)
	}
	return row
}

func (d *DB) Close() {
	d.l.Info("closing db")
	d.db.Close()
}

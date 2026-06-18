package stats

import (
	"context"
	"fmt"
	"log/slog"
	"olimotracker/pkg/db"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Repository interface {
	GetByUserID(ctx context.Context, id *uuid.UUID) (*UserStats, error)
	Create(ctx context.Context, stats *UserStats) error
	Update(ctx context.Context, stats *UserStats) error
	UpdateStreaks(ctx context.Context) error
	UpdateGoal(ctx context.Context, userID *uuid.UUID, goal int64) error
}

type repo struct {
	db db.DBClient
	l  *slog.Logger
}

func NewRepo(db db.DBClient, l *slog.Logger) Repository {
	return &repo{
		db: db,
		l:  l,
	}
}

func (r *repo) GetByUserID(ctx context.Context, id *uuid.UUID) (*UserStats, error) {
	var res UserStats

	builder := squirrel.Select(fmt.Sprintf("us.%v, us.%v, us.%v, us.%v, us.%v, us.%v, us.%v, us.%v, us.%v, us.%v, u.%v", db.UserStatsUserIDColumn, db.UserStatsTotalHoursColumn, db.UserStatsCurrentStreakColumn, db.UserStatsMaxStreakColumn, db.UserStatsLevelColumn, db.UserStatsXPColumn, db.UserStatsCreatedAtColumn, db.UserStatsUpdatedAtColumn, db.UserStatsLastSessionsAtColumn, db.UserStatsGoalColumn, db.UsersUsernameColumn)).
		From(db.UserStatsTable + " us").
		LeftJoin(fmt.Sprintf("%v u ON u.%v = us.%v", db.UsersTable, db.UsersIDColumn, db.UserStatsUserIDColumn)).
		Where(squirrel.Eq{"us." + db.UserStatsUserIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.l.Error("error building query: ", "err", err)
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&res.UserID, &res.TotalMinutes, &res.CurrentStreak, &res.MaxStreak, &res.Level, &res.XP, &res.CreatedAt, &res.UpdatedAt, &res.LastSessionAt, &res.Goal, &res.Username)
	if err != nil {
		r.l.Error("error scanning row: ", "err", err)
		return nil, err
	}

	return &res, nil
}

func (r *repo) Create(ctx context.Context, stats *UserStats) error {
	builder := squirrel.Insert(db.UserStatsTable).
		Columns(db.UserStatsUserIDColumn, db.UserStatsTotalHoursColumn, db.UserStatsCurrentStreakColumn, db.UserStatsMaxStreakColumn, db.UserStatsLevelColumn, db.UserStatsXPColumn, db.UserStatsUpdatedAtColumn, db.UserStatsLastSessionsAtColumn).
		Values(stats.UserID, stats.TotalMinutes, stats.CurrentStreak, stats.MaxStreak, stats.Level, stats.XP, stats.UpdatedAt, stats.LastSessionAt).
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

func (r *repo) UpdateGoal(ctx context.Context, userID *uuid.UUID, goal int64) error {
	builder := squirrel.Update(db.UserStatsTable).
		Set(db.UserStatsGoalColumn, goal).
		Where(squirrel.Eq{db.UserStatsUserIDColumn: userID}).
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

func (r *repo) Update(ctx context.Context, stats *UserStats) error {
	builder := squirrel.Update(db.UserStatsTable).
		Set(db.UserStatsTotalHoursColumn, stats.TotalMinutes).
		Set(db.UserStatsCurrentStreakColumn, stats.CurrentStreak).
		Set(db.UserStatsMaxStreakColumn, stats.MaxStreak).
		Set(db.UserStatsLevelColumn, stats.Level).
		Set(db.UserStatsXPColumn, stats.XP).
		Set(db.UserStatsIsStudyTodayColumn, stats.IsStudyToday).
		Set(db.UserStatsLastSessionsAtColumn, stats.LastSessionAt).
		Set(db.UserStatsUpdatedAtColumn, time.Now()).
		Where(squirrel.Eq{db.UserStatsUserIDColumn: stats.UserID}).
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

func (r *repo) UpdateStreaks(ctx context.Context) error {
	builder := squirrel.Update(db.UserStatsTable).
		Set(db.UserStatsCurrentStreakColumn, 0).
		Where(squirrel.Eq{db.UserStatsIsStudyTodayColumn: false}).
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

	builder = squirrel.Update(db.UserStatsTable).
		Set(db.UserStatsIsStudyTodayColumn, false).
		Where(squirrel.Eq{db.UserStatsIsStudyTodayColumn: true}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err = builder.ToSql()
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

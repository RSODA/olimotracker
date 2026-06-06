package stats

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	RecalculateStats(ctx context.Context, userID *uuid.UUID, sessionDuration int) error
	GetByUserID(ctx context.Context, userID *uuid.UUID) (*UserStats, error)
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

func (s *service) RecalculateStats(ctx context.Context, userID *uuid.UUID, sessionDuration int) error {
	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			xp := sessionDuration
			level := levelLogic(xp)

			err := s.repo.Create(ctx, &UserStats{UserID: *userID, XP: xp, Level: level, CurrentStreak: 1, MaxStreak: 1})
			if err != nil {
				s.l.Error("error creating stats", "userID", userID, "error", err)
				return err
			}
			return nil
		}

		s.l.Error("error getting stats", "userID", userID, "error", err)
		return err
	}

	xp := stats.XP + sessionDuration
	level := levelLogic(xp)
	newStreak := calcStreak(stats)
	maxStreak := stats.MaxStreak

	if maxStreak < newStreak {
		maxStreak = newStreak
	}

	now := time.Now()

	err = s.repo.Update(ctx, &UserStats{
		UserID:        *userID,
		TotalMinutes:  stats.TotalMinutes + sessionDuration,
		XP:            xp,
		Level:         level,
		CurrentStreak: newStreak,
		MaxStreak:     maxStreak,
		LastSessionAt: &now,
	})
	if err != nil {
		s.l.Error("error updating stats", "userID", userID, "error", err)
		return err
	}

	return nil
}

func (s *service) GetByUserID(ctx context.Context, userID *uuid.UUID) (*UserStats, error) {
	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		s.l.Error("error getting stats", "userID", userID, "error", err)
		return nil, err
	}

	return stats, nil
}

func calcStreak(stats *UserStats) int {
	moscow := time.FixedZone("Moscow", 3*60*60)
	now := time.Now().In(moscow)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, moscow)
	yestardayStart := todayStart.Add(-24 * time.Hour)

	if stats.LastSessionAt == nil {
		return 1
	}

	last := stats.LastSessionAt.In(moscow)

	if !last.Before(todayStart) {
		return stats.CurrentStreak
	}
	if !last.Before(yestardayStart) {
		return stats.CurrentStreak + 1
	}

	return 1
}

func levelLogic(xp int) int {
	return (xp / 100) + 1
}

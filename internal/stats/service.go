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
	AdjustStats(ctx context.Context, userID *uuid.UUID, oldDuration int, newDuration int) error
	RecalculateStats(ctx context.Context, userID *uuid.UUID, sessionDuration int) error
	GetByUserID(ctx context.Context, userID *uuid.UUID) (*UserStatsResponse, error)
	UpdateStreaks(ctx context.Context) error
	UpdateGoal(ctx context.Context, userID *uuid.UUID, goal int64) error
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
		IsStudyToday:  true,
		LastSessionAt: &now,
	})
	if err != nil {
		s.l.Error("error updating stats", "userID", userID, "error", err)
		return err
	}

	return nil
}

func (s *service) UpdateGoal(ctx context.Context, userID *uuid.UUID, goal int64) error {
	err := s.repo.UpdateGoal(ctx, userID, goal)
	if err != nil {
		s.l.Error("error updating goal", "userID", userID, "error", err)
		return err
	}

	return nil
}

func (s *service) GetByUserID(ctx context.Context, userID *uuid.UUID) (*UserStatsResponse, error) {
	var res UserStatsResponse

	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		s.l.Error("error getting stats", "userID", userID, "error", err)
		return nil, err
	}

	res = UserStatsResponse{
		UserID:        stats.UserID,
		Username:      stats.Username,
		TotalMinutes:  stats.TotalMinutes,
		XP:            stats.XP,
		Level:         stats.Level,
		IsStudyToday:  stats.IsStudyToday,
		CurrentStreak: stats.CurrentStreak,
		MaxStreak:     stats.MaxStreak,
		Goal:          stats.Goal,
		GalaxySeed:    stats.GalaxySeed,
		LastSessionAt: stats.LastSessionAt,
	}

	return &res, nil
}

func (s *service) AdjustStats(ctx context.Context, userID *uuid.UUID, oldDuration int, newDuration int) error {
	stats, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	diff := oldDuration - newDuration
	newXp := stats.XP + diff

	if newXp < 0 {
		newXp = 0
	}

	newTotalMinutes := stats.TotalMinutes + newDuration
	if newTotalMinutes < 0 {
		newTotalMinutes = 0
	}

	err = s.repo.Update(ctx, &UserStats{
		UserID:        *userID,
		TotalMinutes:  newTotalMinutes,
		XP:            newXp,
		Level:         levelLogic(newXp),
		CurrentStreak: stats.CurrentStreak,
		MaxStreak:     stats.MaxStreak,
		LastSessionAt: stats.LastSessionAt,
	})

	if err != nil {
		s.l.Error("error updating stats", "userID", userID, "err", err)
		return err
	}

	return nil
}

func (s *service) UpdateStreaks(ctx context.Context) error {
	err := s.repo.UpdateStreaks(ctx)
	if err != nil {
		s.l.Error("error updating streaks", "err", err)
		return err
	}

	return nil
}

func calcStreak(stats *UserStats) int {
	if !stats.IsStudyToday {
		if stats.LastSessionAt == nil {
			return 1
		}

		return stats.CurrentStreak + 1
	}

	return stats.CurrentStreak
}

func levelLogic(xp int) int {
	return (xp / 100) + 1
}

package galaxy

import (
	"context"
	"log/slog"
	"olimotracker/internal/sessions"
	"olimotracker/internal/stats"

	"github.com/google/uuid"
)

type Service interface {
	GetGalaxy(ctx context.Context, userID *uuid.UUID) (*GalaxyResponse, error)
}

type service struct {
	statsService stats.Service
	sessionsRepo sessions.Repository
	l            *slog.Logger
}

func NewService(statsService stats.Service, sessionsRepo sessions.Repository, l *slog.Logger) Service {
	return &service{
		statsService: statsService,
		sessionsRepo: sessionsRepo,
		l:            l,
	}
}

func (s *service) GetGalaxy(ctx context.Context, userID *uuid.UUID) (*GalaxyResponse, error) {
	categories, err := s.sessionsRepo.GetMinutesByCategoryForUser(ctx, userID)
	if err != nil {
		s.l.Error("failed to get minutes by category for user", "error", err)
		return nil, err
	}

	stats, err := s.statsService.GetByUserID(ctx, userID)
	if err != nil {
		s.l.Error("failed to get stats", "error", err)
		return nil, err
	}

	return &GalaxyResponse{
		Seed:          stats.GalaxySeed,
		TotalMinutes:  stats.XP,
		CurrentStreak: stats.CurrentStreak,
		MaxStreak:     stats.MaxStreak,
		Level:         stats.Level,
		Categories:    categories,
	}, nil
}

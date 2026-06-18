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
	RegenerateAllSeeds(ctx context.Context) error
}

type service struct {
	statsService stats.Service
	sessionsRepo sessions.Repository
	l            *slog.Logger
	se           Repository
}

func NewService(statsService stats.Service, sessionsRepo sessions.Repository, l *slog.Logger, se Repository) Service {
	return &service{
		statsService: statsService,
		sessionsRepo: sessionsRepo,
		l:            l,
		se:           se,
	}
}

func (s *service) GetGalaxy(ctx context.Context, userID *uuid.UUID) (*GalaxyResponse, error) {
	grouped, err := s.sessionsRepo.GetSessionsForCategory(ctx, userID)
	if err != nil {
		s.l.Error("failed to get minutes by category for user", "error", err)
		return nil, err
	}

	stats, err := s.statsService.GetByUserID(ctx, userID)
	if err != nil {
		s.l.Error("failed to get stats", "error", err)
		return nil, err
	}

	result := make([]CategorySessions, 0, len(grouped))

	for category, sessions := range grouped {
		result = append(result, CategorySessions{
			Category: category,
			Sessions: sessions,
		})
	}

	return &GalaxyResponse{
		Seed:             stats.GalaxySeed,
		TotalMinutes:     stats.XP,
		CurrentStreak:    stats.CurrentStreak,
		MaxStreak:        stats.MaxStreak,
		Level:            stats.Level,
		CategorySessions: result,
	}, nil
}

func (s *service) RegenerateAllSeeds(ctx context.Context) error {
	err := s.se.UpdateAllSeeds(ctx)
	if err != nil {
		s.l.Error("failed to update all seeds", "error", err)
		return err
	}
	return nil
}

package sessions

import (
	"context"
	"log/slog"
	"olimotracker/internal/stats"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID *uuid.UUID, req *CreateSessionRequest) (*uuid.UUID, error)
	GetByID(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID) (*SessionResponse, error)
	GetByUserID(ctx context.Context, userID *uuid.UUID) ([]*SessionResponse, error)
	GetByCategoryID(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) ([]*SessionResponse, error)
	Update(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID, req *UpdateSessionRequest) (*uuid.UUID, error)
	Delete(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID) error
	GetMinutesBySessionForUser(ctx context.Context, userID *uuid.UUID) ([]SessionsMinutes, error)
}

type service struct {
	repo         Repository
	l            *slog.Logger
	statsService stats.Service
}

func NewService(repo Repository, l *slog.Logger, statsService stats.Service) Service {
	return &service{repo: repo, l: l, statsService: statsService}
}

func (s *service) Create(ctx context.Context, userID *uuid.UUID, req *CreateSessionRequest) (*uuid.UUID, error) {
	session := &Session{
		UserID:     *userID,
		CategoryID: req.CategoryID,
		Duration:   req.Duration,
		Note:       req.Note,
	}

	res, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	err = s.statsService.RecalculateStats(ctx, userID, int(session.Duration))
	if err != nil {
		s.l.Error("failed to recalculate stats", slog.Any("error", err))
	}

	return res, nil
}

func (s *service) GetByID(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID) (*SessionResponse, error) {
	res, err := s.repo.GetByID(ctx, *sessionID, *userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetByUserID(ctx context.Context, userID *uuid.UUID) ([]*SessionResponse, error) {
	res, err := s.repo.GetByUserID(ctx, *userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GetByCategoryID(ctx context.Context, categoryID *uuid.UUID, userID *uuid.UUID) ([]*SessionResponse, error) {
	res, err := s.repo.GetByCategoryID(ctx, *categoryID, *userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID, req *UpdateSessionRequest) (*uuid.UUID, error) {
	session := &Session{
		CategoryID: req.CategoryID,
		Note:       req.Note,
	}

	if req.Duration != nil {
		session.Duration = *req.Duration
	}

	oldStats, err := s.repo.GetByID(ctx, *sessionID, *userID)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Update(ctx, *sessionID, *userID, session)
	if err != nil {
		return nil, err
	}

	if req.Duration != nil {
		err = s.statsService.AdjustStats(ctx, userID, oldStats.Duration, *req.Duration)
		if err != nil {
			s.l.Error("failed to adjust stats", "sessionID", sessionID, "userID", userID, "err", err)
		}
	}

	return res, nil
}

func (s *service) Delete(ctx context.Context, sessionID *uuid.UUID, userID *uuid.UUID) error {
	return s.repo.Delete(ctx, *sessionID, *userID)
}

func (s *service) GetMinutesBySessionForUser(ctx context.Context, userID *uuid.UUID) ([]SessionsMinutes, error) {
	res, err := s.repo.GetMinutesBySessionForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

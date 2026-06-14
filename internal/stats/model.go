package stats

import (
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	UserID        uuid.UUID
	Username      string
	TotalMinutes  int
	XP            int
	Level         int
	CurrentStreak int
	MaxStreak     int
	GalaxySeed    int64
	LastSessionAt *time.Time
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

type UserStatsResponse struct {
	UserID        uuid.UUID  `json:"user_id"`
	Username      string     `json:"username"`
	TotalMinutes  int        `json:"total_minutes"`
	TotalHours    int        `json:"total_hours"`
	XP            int        `json:"xp"`
	Level         int        `json:"level"`
	CurrentStreak int        `json:"current_streak"`
	MaxStreak     int        `json:"max_streak"`
	GalaxySeed    int64      `json:"galaxy_seed"`
	LastSessionAt *time.Time `json:"last_session_at"`
}

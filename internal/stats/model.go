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
	Level         int   `json:"level"`
	XP            int   `json:"xp"`
	TotalHours    int   `json:"total_hours"`
	TotalMinutes  int   `json:"total_minutes"`
	CurrentStreak int   `json:"current_streak"`
	MaxStreak     int   `json:"max_streak"`
	GalaxySeed    int64 `json:"galaxy_seed"`
}

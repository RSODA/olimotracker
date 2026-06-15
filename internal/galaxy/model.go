package galaxy

import "olimotracker/internal/sessions"

type GalaxyResponse struct {
	Seed             int64              `json:"seed"`
	TotalMinutes     int                `json:"total_minutes"`
	CurrentStreak    int                `json:"current_streak"`
	MaxStreak        int                `json:"max_streak"`
	Level            int                `json:"level"`
	CategorySessions []CategorySessions `json:"categories"`
}

type CategorySessions struct {
	Category sessions.Category
	Sessions []sessions.SessionsMinutes
}

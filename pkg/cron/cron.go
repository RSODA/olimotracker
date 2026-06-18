package cron

import (
	"context"
	"log/slog"
	"olimotracker/internal/galaxy"
	"olimotracker/internal/stats"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	c            *cron.Cron
	l            *slog.Logger
	gs           galaxy.Service
	statsService stats.Service
}

func NewCron(c *cron.Cron, l *slog.Logger, gs galaxy.Service, statsService stats.Service) *Cron {
	return &Cron{
		c:            c,
		l:            l,
		gs:           gs,
		statsService: statsService,
	}
}

func (c *Cron) AddsCronJobs() {
	_, err := c.c.AddFunc("0 0 * * 0", func() {
		if err := c.gs.RegenerateAllSeeds(context.Background()); err != nil {
			c.l.Error("failed to regenerate seeds", "err", err)
		}
	})
	if err != nil {
		c.l.Error("failed to add cron job", "err", err)
	}

	_, err = c.c.AddFunc("0 0 * * *", func() {
		if err := c.statsService.UpdateStreaks(context.Background()); err != nil {
			c.l.Error("failed to update streaks", "err", err)
		}
	})
}

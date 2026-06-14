package cron

import (
	"context"
	"log/slog"
	"olimotracker/internal/galaxy"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	c  *cron.Cron
	l  *slog.Logger
	gs galaxy.Service
}

func NewCron(c *cron.Cron, l *slog.Logger, gs galaxy.Service) *Cron {
	return &Cron{
		c:  c,
		l:  l,
		gs: gs,
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
}

package logger

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	colorReset = "\033[0m"

	bgRed    = "\033[41m"
	bgYellow = "\033[43m"
	bgCyan   = "\033[46m"
	bgWhite  = "\033[47m"
	bgBlue   = "\033[44m"
)

type PrettyHandler struct {
	opts slog.HandlerOptions
}

func NewPrettyHandler(opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{opts: *opts}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	var color string
	var emoji string

	switch r.Level {
	case slog.LevelDebug:
		color = bgCyan
		emoji = "🔍"
	case slog.LevelInfo:
		color = bgBlue
		emoji = "ℹ️"
	case slog.LevelWarn:
		color = bgYellow
		emoji = "⚠️"
	case slog.LevelError:
		color = bgRed
		emoji = "❌"
	default:
		color = bgWhite
		emoji = "ℹ️"
	}

	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	logLine := fmt.Sprintf("[%s] %s  %s %s %s: %s%s%s\n",
		r.Time.Format("2006-01-02 15:04:05"),
		emoji,
		color,
		r.Level.String(),
		colorReset,
		r.Message,
		attrs,
		colorReset,
	)
	_, err := fmt.Print(logLine)
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h
}

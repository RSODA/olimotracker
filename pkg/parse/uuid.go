package parse

import (
	"log/slog"

	"github.com/google/uuid"
)

func ParseUUID(s string, l *slog.Logger) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		l.Error("error pars string to uuid", "string", s, "err", err)
		return uuid.Nil, err
	}

	return id, nil
}

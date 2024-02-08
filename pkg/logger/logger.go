package logger

import (
	"log/slog"
	"os"
)

func New(isLocal bool) *slog.Logger {
	if isLocal {
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

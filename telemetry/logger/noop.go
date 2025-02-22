package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/io"
)

func newNoopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&io.NoopWriter{}, nil))
}

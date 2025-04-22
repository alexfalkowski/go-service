package logger

import "log/slog"

// Error just wraps a slog.Any with key of error.
func Error(err error) slog.Attr {
	if err != nil {
		return slog.Any("error", err)
	}

	return slog.Attr{}
}

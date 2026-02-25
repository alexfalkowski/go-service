package logger

import "log/slog"

// Error returns a standardized error attribute for log records.
//
// When err is non-nil it returns `slog.Any("error", err)`. When err is nil it returns
// an empty `slog.Attr{}` so callers can append it unconditionally.
//
// The key is always "error" to keep error fields consistent across handlers/exporters.
func Error(err error) slog.Attr {
	if err != nil {
		return slog.Any("error", err)
	}
	return slog.Attr{}
}

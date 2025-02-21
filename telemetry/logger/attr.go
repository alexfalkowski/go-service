package logger

import "log/slog"

// Error returns an Attr for a string value if err is set, otherwise returns an empty Attr.
func Error(err error) slog.Attr {
	if err != nil {
		return slog.Attr{Key: "error", Value: slog.StringValue(err.Error())}
	}

	return slog.Attr{}
}

package logger

import (
	"context"
)

// Logger is an empty logger.
type Logger struct{}

// NewLogger is empty.
func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Printf(ctx context.Context, format string, v ...any) {}

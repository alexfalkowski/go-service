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

func (l *Logger) Printf(_ context.Context, _ string, _ ...any) {}

package zap

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Logger for redis.
type Logger struct {
	logger *zap.Logger
}

// NewLogger for redis.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Printf(_ context.Context, format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

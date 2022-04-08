package zap

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Logger for zap.
type Logger struct {
	logger *zap.Logger
}

// NewLogger for zap.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Printf(ctx context.Context, format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

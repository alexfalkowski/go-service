package zap

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger for proxy.
type Logger struct {
	logger *zap.Logger
}

// NewLogger for proxy.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

// Printf for proxy.
func (l *Logger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

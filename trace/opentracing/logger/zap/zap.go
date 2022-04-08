package zap

import (
	"fmt"
	"strings"

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

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Infof(msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Log(msg string) {
	if strings.Contains(msg, "INFO") {
		l.logger.Info(msg)

		return
	}

	if strings.Contains(msg, "WARN") {
		l.logger.Warn(msg)

		return
	}

	if strings.Contains(msg, "ERROR") {
		l.logger.Error(msg)

		return
	}
}

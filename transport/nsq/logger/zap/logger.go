package zap

import (
	"strings"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger for zap.
type Logger struct {
	logger *zap.Logger
}

// NewLogger for zap.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

// Output for the logger.
func (l *Logger) Output(calldepth int, s string) error {
	fields := []zapcore.Field{
		zap.String(component, nsqComponent),
	}

	if strings.HasPrefix(s, nsq.LogLevelInfo.String()) {
		l.logger.Info(s, fields...)

		return nil
	}

	if strings.HasPrefix(s, nsq.LogLevelWarning.String()) {
		l.logger.Warn(s, fields...)

		return nil
	}

	if strings.HasPrefix(s, nsq.LogLevelError.String()) {
		l.logger.Error(s, fields...)

		return nil
	}

	l.logger.Debug(s, fields...)

	return nil
}

package zap

import (
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

// NewLogger for zap.
func NewLogger(logger *zap.Logger) ddtrace.Logger {
	return &dataDoglogger{logger: logger}
}

type dataDoglogger struct {
	logger *zap.Logger
}

func (l *dataDoglogger) Log(msg string) {
	l.logger.Info(msg)
}

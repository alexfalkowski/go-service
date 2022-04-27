package logger

// Logger is an empty logger.
type Logger struct{}

// NewLogger is empty.
func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(msg string) {}

package logger

// Logger is an empty logger.
type Logger struct{}

// NewLogger is empty.
func NewLogger() *Logger {
	return &Logger{}
}

// Output for the logger.
func (l *Logger) Output(_ int, _ string) error {
	return nil
}

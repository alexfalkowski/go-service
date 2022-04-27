package logger

// Logger is an empty logger.
type Logger struct{}

// NewLogger is empty.
func NewLogger() *Logger {
	return &Logger{}
}

// Output for the logger.
func (l *Logger) Output(calldepth int, s string) error {
	return nil
}

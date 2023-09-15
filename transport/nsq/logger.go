package nsq

type logger struct{}

// Output for the logger.
func (l *logger) Output(_ int, _ string) error {
	return nil
}

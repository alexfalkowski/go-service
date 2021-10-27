package retry

// Config for retry.
type Config struct {
	Timeout  uint `yaml:"timeout"`
	Attempts uint `yaml:"attempts"`
}

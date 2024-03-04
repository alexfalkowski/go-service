package retry

import (
	"time"
)

// Config for retry.
type Config struct {
	Timeout  time.Duration `yaml:"timeout" json:"timeout" toml:"timeout"`
	Attempts uint          `yaml:"attempts" json:"attempts" toml:"attempts"`
}

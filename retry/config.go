package retry

import (
	"time"
)

// Config for retry.
type Config struct {
	Timeout  time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Attempts uint          `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}

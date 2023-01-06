package retry

import (
	"time"
)

// Config for retry.
type Config struct {
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
	Attempts uint          `yaml:"attempts" json:"attempts"`
}

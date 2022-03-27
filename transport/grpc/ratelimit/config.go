package ratelimit

import (
	"time"
)

// Config for ratelimit.
type Config struct {
	Every time.Duration `yaml:"every"`
	Burst uint          `yaml:"burst"`
	TTL   time.Duration `yaml:"ttl"`
}

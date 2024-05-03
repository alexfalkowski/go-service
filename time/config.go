package time

import (
	"github.com/alexfalkowski/go-service/time/ntp"
	"github.com/alexfalkowski/go-service/time/nts"
)

// IsEnabled for time.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for time.
type Config struct {
	NTP *ntp.Config `yaml:"ntp,omitempty" json:"ntp,omitempty" toml:"ntp,omitempty"`
	NTS *nts.Config `yaml:"nts,omitempty" json:"nts,omitempty" toml:"nts,omitempty"`
}

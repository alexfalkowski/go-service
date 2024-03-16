package client

import (
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security"
)

// Config for client.
type Config struct {
	Host      string          `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Security  security.Config `yaml:"security,omitempty" json:"security,omitempty" toml:"security,omitempty"`
	Timeout   string          `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Retry     retry.Config    `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`
	UserAgent string          `yaml:"user_agent,omitempty" json:"user_agent,omitempty" toml:"user_agent,omitempty"`
}

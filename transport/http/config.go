package http

import (
	"github.com/alexfalkowski/go-service/transport/http/retry"
)

// Config for HTTP.
type Config struct {
	Retry     retry.Config `yaml:"retry" json:"retry" toml:"retry"`
	UserAgent string       `yaml:"user_agent" json:"user_agent" toml:"user_agent"`
}

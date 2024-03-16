package http

import (
	"github.com/alexfalkowski/go-service/server"
)

// Config for HTTP.
type Config struct {
	server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

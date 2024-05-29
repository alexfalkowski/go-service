package http

import (
	"github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for HTTP.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// Mux for config.
func Mux(cfg *Config) http.MuxKind {
	if !IsEnabled(cfg) {
		return http.StandardMux
	}

	return http.MuxKind(cfg.Mux)
}

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
	Mux            string `yaml:"mux,omitempty" json:"mux,omitempty" toml:"mux,omitempty"`
}

package events

import (
	"github.com/alexfalkowski/go-service/transport/events/http"
)

// Config for events.
type Config struct {
	HTTP http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}

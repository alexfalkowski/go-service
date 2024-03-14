package feature

import (
	"github.com/alexfalkowski/go-service/client"
)

// Config for feature.
type Config struct {
	Kind          string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

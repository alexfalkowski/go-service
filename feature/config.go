package feature

import (
	"github.com/alexfalkowski/go-service/client"
)

// Config for feature.
type Config struct {
	Kind          string `yaml:"kind" json:"kind" toml:"kind"`
	client.Config `yaml:",inline" json:",inline" toml:",inline"`
}

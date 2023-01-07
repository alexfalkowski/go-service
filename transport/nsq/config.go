package nsq

import (
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
)

type Config struct {
	LookupHost string       `yaml:"lookup_host" json:"lookup_host" toml:"lookup_host"`
	Host       string       `yaml:"host" json:"host" toml:"host"`
	Retry      retry.Config `yaml:"retry" json:"retry" toml:"retry"`
	UserAgent  string       `yaml:"user_agent" json:"user_agent" toml:"user_agent"`
}

package nsq

import (
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
)

type Config struct {
	LookupHost string       `yaml:"lookup_host" json:"lookup_host"`
	Host       string       `yaml:"host" json:"host"`
	Retry      retry.Config `yaml:"retry" json:"retry"`
	UserAgent  string       `yaml:"user_agent" json:"user_agent"`
}

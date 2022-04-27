package nsq

import (
	"github.com/alexfalkowski/go-service/transport/nsq/retry"
)

type Config struct {
	LookupHost string       `yaml:"lookup_host"`
	Host       string       `yaml:"host"`
	Retry      retry.Config `yaml:"retry"`
	UserAgent  string       `yaml:"user_agent"`
}

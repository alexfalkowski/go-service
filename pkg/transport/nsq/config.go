package nsq

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/nsqio/go-nsq"
)

type Config struct {
	LookupHost string `envconfig:"NSQ_LOOKUP_HOST" required:"true" default:"localhost:4161"`
	Host       string `envconfig:"NSQ_HOST" required:"true" default:"localhost:4150"`

	Config *nsq.Config
}

// NewConfig for NSQ.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	config.Config = nsq.NewConfig()

	return &config, err
}

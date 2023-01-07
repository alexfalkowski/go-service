package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// Config for transport.
type Config struct {
	Port string      `yaml:"port" json:"port" toml:"port"`
	GRPC grpc.Config `yaml:"grpc" json:"grpc" toml:"grpc"`
	HTTP http.Config `yaml:"http" json:"http" toml:"http"`
	NSQ  nsq.Config  `yaml:"nsq" json:"nsq" toml:"nsq"`
}

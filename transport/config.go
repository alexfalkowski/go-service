package transport

import (
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// Config for transport.
type Config struct {
	GRPC grpc.Config `yaml:"grpc"`
	HTTP http.Config `yaml:"http"`
	NSQ  nsq.Config  `yaml:"nsq"`
}

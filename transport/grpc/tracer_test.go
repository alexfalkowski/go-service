package grpc_test

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

func init() {
	tracer.Register()
}

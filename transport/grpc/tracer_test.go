package grpc_test

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

//nolint:gochecknoinits
func init() {
	tracer.Register()
}

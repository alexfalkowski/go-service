package grpc

import (
	"github.com/alexfalkowski/go-health/pkg/subscriber"
)

// Observer for gRPC.
type Observer struct {
	*subscriber.Observer
}

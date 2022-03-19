package grpc

import (
	"github.com/alexfalkowski/go-health/subscriber"
)

// Observer for gRPC.
type Observer struct {
	*subscriber.Observer
}

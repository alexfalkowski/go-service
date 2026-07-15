package errors

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
)

// ServerError normalizes expected gRPC server shutdown errors.
//
// gRPC returns [grpc.ErrServerStopped] when (*[grpc.Server]).Serve is called
// after Stop or GracefulStop. That order can occur during normal asynchronous
// lifecycle shutdown, so it is not considered a serve failure.
//
// ServerError returns nil when err is [grpc.ErrServerStopped], otherwise it
// returns err unchanged.
func ServerError(err error) error {
	if errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}

	return err
}

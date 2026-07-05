package status

import (
	"github.com/alexfalkowski/go-service/v2/time"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/durationpb"
)

// RetryInfo aliases the standard google.rpc RetryInfo detail.
//
// It is commonly attached to retryable gRPC status errors to communicate how
// long clients should wait before retrying the same request.
type RetryInfo = errdetails.RetryInfo

// Duration aliases the protobuf duration type used by structured status details.
type Duration = durationpb.Duration

// NewDuration returns a protobuf duration for structured gRPC status details.
func NewDuration(d time.Duration) *Duration {
	return durationpb.New(d.Duration())
}

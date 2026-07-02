package retry

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/transport/retry"
)

// Config configures retry behavior for gRPC unary client calls.
//
// It embeds shared retry mechanics and adds gRPC-specific failure classification.
type Config struct {
	// Config carries shared retry mechanics such as attempts, timeout, and backoff.
	*retry.Config `yaml:",inline" json:",inline" toml:",inline"`

	// Codes replaces the default retryable gRPC status codes.
	//
	// When empty, gRPC retry uses its default failure classification. When set, only these
	// configured non-OK status codes are retryable; include the defaults here as well when
	// custom configuration should extend rather than replace the default list.
	Codes []codes.Code `yaml:"codes,omitempty" json:"codes,omitempty" toml:"codes,omitempty" validate:"omitempty,dive,gt=0,lte=16"`
}
